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

func syncRepo(env *goenv.ExecEnv, repo, src, commit string) error {
	// TODO(h8liu): check branch commit first, if already synced, do nothing.
	var err error
	curCommit := ""

	xa := func(dir string, name string, args ...string) {
		if err != nil {
			return
		}
		err = env.Exec(dir, name, args...)
	}
	x := func(name string, args ...string) {
		xa("", name, args...)
	}

	srcDir := env.SrcDir(repo)
	if exist, err := env.IsDir(srcDir); err != nil {
		return err
	} else if !exist {
		x("mkdir", "-p", srcDir)
	}

	srcGitDir := filepath.Join(srcDir, ".git")
	if exist, err := env.IsDir(srcGitDir); err != nil {
		return err
	} else if !exist {
		xa(srcDir, "git", "init", "-q")
	} else {
		curCommit, err = env.StrOut(
			srcDir, "git", "show", "HEAD", "-s", "--format=%H",
		)
		if err != nil {
			return err
		}
		curCommit = strings.TrimSpace(curCommit)
	}

	if curCommit == commit {
		return nil
	}

	fmt.Printf(
		"[%s -> %s] %s - %s\n",
		idutil.Short(curCommit), idutil.Short(commit), repo, src,
	)
	xa(srcDir, "git", "fetch", "-q", src)
	xa(srcDir, "git", "branch", "-q", "-f", "smlrepo", commit)
	xa(srcDir, "git", "merge", "-q", "smlrepo")

	return err
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
