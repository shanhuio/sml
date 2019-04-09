package sml

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/base/httputil"
	"shanhu.io/base/idutil"
	"shanhu.io/sml/core"
	"shanhu.io/sml/goenv"
)

func currentCommit(env *goenv.ExecEnv, srcDir string) (string, error) {
	ret, err := env.StrOut(
		srcDir, "git", "show", "HEAD", "-s", "--format=%H",
	)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(ret), nil
}

func execAll(env *goenv.ExecEnv, srcDir string, lines [][]string) error {
	for _, args := range lines {
		if err := env.Exec(srcDir, args[0], args[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func syncRepo(env *goenv.ExecEnv, repo, src, commit string) (bool, error) {
	if strings.HasPrefix("commit", "hg/") {
		err := fmt.Errorf(
			"%q: mercurial repo support not implemented yet", repo,
		)
		return false, err
	}
	srcDir := goenv.SrcDir(repo)
	if exist, err := env.IsDir(srcDir); err != nil {
		return false, err
	} else if !exist {
		if err := env.Exec("", "mkdir", "-p", srcDir); err != nil {
			return false, err
		}
	}

	newRepo := false
	srcGitDir := filepath.Join(srcDir, ".git")
	if exist, err := env.IsDir(srcGitDir); err != nil {
		return false, err
	} else if !exist {
		if err := env.Exec(srcDir, "git", "init", "-q"); err != nil {
			return false, err
		}

		if err := env.Exec(
			srcDir, "git", "remote", "add", "origin", src,
		); err != nil {
			return false, err
		}

		fmt.Printf(
			"[new %s] %s\n", idutil.Short(commit), repo,
		)
		newRepo = true
	} else {
		cur, err := currentCommit(env, srcDir)
		if err != nil {
			return false, err
		}
		if cur == commit {
			return false, nil
		}

		isAncestor, err := env.Call(
			srcDir, "git", "merge-base", "--is-ancestor", commit, cur,
		)
		if err != nil {
			return false, err
		}
		if isAncestor {
			// merge will be a noop, just update smlrepo branch.
			return false, env.Exec(
				srcDir, "git", "branch", "-q", "-f", "smlrepo", commit,
			)
		}

		fmt.Printf(
			"[%s..%s] %s\n",
			idutil.Short(cur), idutil.Short(commit), repo,
		)
	}

	// fetch to smlrepo branch and then merge
	if err := execAll(env, srcDir, [][]string{
		{"git", "fetch", "-q", src},
		{"git", "branch", "-q", "-f", "smlrepo", commit},
		{"git", "merge", "-q", "smlrepo"},
	}); err != nil {
		return false, err
	}

	if newRepo {
		if err := execAll(env, srcDir, [][]string{
			{"git", "fetch", "-q", "origin"},
			{
				"git", "branch", "-q",
				"--set-upstream-to=origin/master", "master",
			},
		}); err != nil {
			return false, err
		}
	}

	return true, nil
}

// ThisRepo is the package name of this repo.
const ThisRepo = "shanhu.io/sml"

// the only repo that ThisRepo depends.
const baseRepo = "smallrepo.com/base"

func installThis(env *goenv.ExecEnv, thisRepo string) error {
	return env.Exec(goenv.SrcDir(thisRepo), "go", "install", thisRepo)
}

// SyncTo syncs to the desired state.
func SyncTo(state *core.State, thisRepo string) error {
	fmt.Printf("#%d  [%s]\n", state.Clock, idutil.Short(state.ID))

	if len(state.Commits) == 0 {
		return fmt.Errorf("got no commits")
	}

	var repos []string
	repoMap := make(map[string]bool)
	for repo := range state.Commits {
		repos = append(repos, repo)
		repoMap[repo] = true
	}
	if !repoMap[thisRepo] {
		repos = append(repos, thisRepo)
	}
	sort.Strings(repos)

	gopath, err := goenv.GOPATH()
	if err != nil {
		return err
	}
	env := goenv.NewExecEnv(gopath)

	for _, repo := range repos {
		commit := state.Commits[repo]
		src := state.Sources[repo]
		if _, err := syncRepo(env, repo, src, commit); err != nil {
			return err
		}
	}

	if err := installThis(env, thisRepo); err != nil {
		return err
	}
	return nil
}

func doSync(server string, profile *Profile) error {
	if len(profile.Tracking) == 0 {
		return fmt.Errorf("nothing is tracked")
	}

	state := new(core.State)
	query := &core.StateQuery{
		Tracking: profile.Tracking,
	}
	c := httputil.NewClient(server)
	if err := c.JSONCall("/api/sync", query, state); err != nil {
		return fmt.Errorf("sync error: %s", err)
	}

	return SyncTo(state, ThisRepo)
}

func sync(server string, args []string) error {
	flags := newFlags()
	flags.Parse(args)
	if len(args) > 0 {
		return fmt.Errorf("sync accepts no arguments")
	}
	profileName := defaultProfileName
	profile, err := loadProfile(profileName)
	if err != nil {
		return err
	}
	return doSync(server, profile)
}
