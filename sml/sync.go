package sml

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/misc/httputil"
	"shanhu.io/misc/idutil"
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

func syncRepo(env *goenv.ExecEnv, repo, src, commit string) error {
	srcDir := env.SrcDir(repo)
	if exist, err := env.IsDir(srcDir); err != nil {
		return err
	} else if !exist {
		if err := env.Exec("", "mkdir", "-p", srcDir); err != nil {
			return err
		}
	}

	srcGitDir := filepath.Join(srcDir, ".git")
	if exist, err := env.IsDir(srcGitDir); err != nil {
		return err
	} else if !exist {
		if err := env.Exec(srcDir, "git", "init", "-q"); err != nil {
			return err
		}

		if err := env.Exec(
			srcDir, "git", "remote", "add", "origin", src,
		); err != nil {
			return err
		}

		fmt.Printf(
			"[new %s] %s -%s\n", idutil.Short(commit), repo, src,
		)
	} else {
		cur, err := currentCommit(env, srcDir)
		if err != nil {
			return err
		}
		if cur == commit {
			return nil
		}

		isAncestor, err := env.Call(
			srcDir, "git", "merge-base", "--is-ancestor", commit, cur,
		)
		if err != nil {
			return err
		}
		if isAncestor {
			// merge will be a noop, just update smlrepo branch.
			return env.Exec(
				srcDir, "git", "branch", "-q", "-f", "smlrepo", commit,
			)
		}

		fmt.Printf(
			"[%s -> %s] %s - %s\n",
			idutil.Short(cur), idutil.Short(commit), repo, src,
		)
	}

	// fetch to smlrepo branch and then merge
	return execAll(env, srcDir, [][]string{
		{"git", "fetch", "-q", src},
		{"git", "branch", "-q", "-f", "smlrepo", commit},
		{"git", "merge", "-q", "smlrepo"},
	})
}

func doSync(server string, profile *Profile) error {
	if len(profile.Tracking) == 0 {
		return fmt.Errorf("nothing is tracked")
	}

	state := new(State)
	c := httputil.NewClient(server)
	if err := c.JSONCall("/api/sync", profile, state); err != nil {
		return fmt.Errorf("sync error: %s", err)
	}

	fmt.Printf("#%d  [%s]\n", state.Clock, idutil.Short(state.ID))

	if len(state.Commits) == 0 {
		return fmt.Errorf("got no commits")
	}

	var repos []string
	for repo := range state.Commits {
		repos = append(repos, repo)
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
		if err := syncRepo(env, repo, src, commit); err != nil {
			return err
		}
	}

	return nil
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
