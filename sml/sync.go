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

func syncRepo(env *goenv.ExecEnv, repo, src, commit string) error {
	curCommit := ""

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
			srcDir, "git", "merge-base", "--is-ancestor", commit, curCommit,
		)
		if err != nil {
			return err
		}
		if isAncestor {
			return nil
		}

		fmt.Printf(
			"[%s -> %s] %s - %s\n",
			idutil.Short(curCommit), idutil.Short(commit), repo, src,
		)
	}

	// fetch to smlrepo branch and then merge
	for _, args := range [][]string{
		{"git", "fetch", "-q", src},
		{"git", "branch", "-q", "-f", "smlrepo", commit},
		{"git", "merge", "-q", "smlrepo"},
	} {
		if err := env.Exec(srcDir, args[0], args[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func sync(server string, args []string) error {
	flags := newFlags()
	force := flags.Bool("force", false, "force syncing all repository")
	flags.Parse(args)
	profileName := defaultProfileName

	if *force {
		fmt.Println("using force")
	}

	profile, err := loadProfile(profileName)
	if err != nil {
		return err
	}
	if len(profile.Tracking) == 0 {
		return fmt.Errorf("nothing is tracked")
	}

	state := new(State)
	c := httputil.NewClient(server)
	err = c.JSONCall("/api/sync", profile, state)
	if err != nil {
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
