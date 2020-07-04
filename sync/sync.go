// Package sync provides functions for synchronizing to gopkgs head.
package sync

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/misc/idutil"
	"shanhu.io/misc/osutil"
	"shanhu.io/sml/core"
	"shanhu.io/sml/goenv"
)

// Syncer is a syncrhonizer that syncs the repos in GOPATH.
type Syncer struct {
	Env *goenv.ExecEnv

	KnownHostsFile string
}

func currentCommit(env *goenv.ExecEnv, srcDir string) (string, error) {
	branches, err := env.StrOut(srcDir, "git", "branch")
	if err != nil {
		return "", fmt.Errorf("list branches: %s", err)
	}
	if strings.TrimSpace(branches) == "" {
		return "", nil
	}

	ret, err := env.StrOut(
		srcDir, "git", "show", "HEAD", "-s", "--format=%H",
	)
	if err != nil {
		return "", fmt.Errorf("get HEAD commit: %s", err)
	}
	return strings.TrimSpace(ret), nil
}

func (s *Syncer) execAll(srcDir string, lines [][]string) error {
	for _, args := range lines {
		if err := s.Env.Exec(srcDir, args[0], args[1:]...); err != nil {
			return err
		}
	}

	return nil
}

func (s *Syncer) execGitFetch(srcDir, src string) error {
	cmd := s.Env.PipedCmd(&goenv.ExecJob{
		Dir:  srcDir,
		Name: "git",
		Args: []string{"fetch", "-q", src, "HEAD"},
	})

	if s.KnownHostsFile != "" {
		f := s.KnownHostsFile
		if strings.Contains(f, `'`) {
			return fmt.Errorf(`knownHostsFile has "'", not supported`)
		}
		if strings.Contains(f, `\`) {
			return fmt.Errorf(`KnownHostsFile has "\", not supported`)
		}
		hasKnownHosts, err := osutil.IsRegular(s.KnownHostsFile)
		if err != nil {
			return err
		}

		if hasKnownHosts {
			gitSSH := fmt.Sprintf(
				`ssh -o UserKnownHostsFile='%s'`, s.KnownHostsFile,
			)
			osutil.CmdAddEnv(cmd, "GIT_SSH_COMMAND", gitSSH)
		}
	}

	return cmd.Run()
}

func (s *Syncer) syncRepo(repo, src, commit string) (bool, error) {
	if strings.HasPrefix("commit", "hg/") {
		err := fmt.Errorf(
			"%q: mercurial repo support not implemented yet", repo,
		)
		return false, err
	}
	srcDir := goenv.SrcDir(repo)
	if exist, err := s.Env.IsDir(srcDir); err != nil {
		return false, err
	} else if !exist {
		if err := s.Env.Exec("", "mkdir", "-p", srcDir); err != nil {
			return false, err
		}
	}

	newRepo := false
	srcGitDir := filepath.Join(srcDir, ".git")
	if exist, err := s.Env.IsDir(srcGitDir); err != nil {
		return false, err
	} else if !exist {
		if err := s.Env.Exec(srcDir, "git", "init", "-q"); err != nil {
			return false, err
		}

		if err := s.Env.Exec(
			srcDir, "git", "remote", "add", "origin", src,
		); err != nil {
			return false, err
		}

		fmt.Printf(
			"[new %s] %s\n", idutil.Short(commit), repo,
		)
		newRepo = true
	} else {
		cur, err := currentCommit(s.Env, srcDir)
		if err != nil {
			return false, err
		}
		if cur == commit {
			return false, nil
		}

		if cur != "" {
			isAncestor, err := s.Env.Call(
				srcDir, "git", "merge-base", "--is-ancestor", commit, cur,
			)
			if err != nil {
				return false, err
			}
			if isAncestor {
				// merge will be a noop, just update smlrepo branch.
				return false, s.Env.Exec(
					srcDir, "git", "branch", "-q", "-f", "smlrepo", commit,
				)
			}

			fmt.Printf(
				"[%s..%s] %s\n",
				idutil.Short(cur), idutil.Short(commit), repo,
			)
		} else {
			fmt.Printf(
				"[new %s] %s\n", idutil.Short(commit), repo,
			)
		}
	}

	// fetch to smlrepo branch and then merge

	if err := s.execGitFetch(srcDir, src); err != nil {
		return false, fmt.Errorf("git fetch: %s", err)
	}

	if err := s.execAll(srcDir, [][]string{
		{"git", "branch", "-q", "-f", "smlrepo", commit},
		{"git", "merge", "-q", "smlrepo"},
	}); err != nil {
		return false, fmt.Errorf("git branch and merge: %s", err)
	}

	if newRepo {
		if err := s.execGitFetch(srcDir, "origin"); err != nil {
			return false, err
		}

		if err := s.execAll(srcDir, [][]string{{
			"git", "branch", "-q", "-f", "master",
		}}); err != nil {
			return false, fmt.Errorf("git setup origin: %s", err)
		}
	}

	return true, nil
}

// AutoInstall is the packaget that will be auto installed on a sync
// operation.
type AutoInstall struct {
	Repo string
	Pkg  string
}

func installThis(env *goenv.ExecEnv, auto *AutoInstall) error {
	return env.Exec(goenv.SrcDir(auto.Pkg), "go", "install", auto.Pkg)
}

// Sync syncs repos in GOPATH to the desired state.
func Sync(env *goenv.ExecEnv, state *core.State, auto *AutoInstall) error {
	s := &Syncer{Env: env}
	return s.Sync(state, auto)
}

// Sync syncs repos in GOPATH to the desired state.
func (s *Syncer) Sync(state *core.State, auto *AutoInstall) error {
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
	if auto != nil {
		if !repoMap[auto.Repo] {
			repos = append(repos, auto.Repo)
		}
	}
	sort.Strings(repos)

	if s.Env == nil {
		gopath, err := goenv.GOPATH()
		if err != nil {
			return err
		}
		s.Env = goenv.NewExecEnv(gopath)
	}

	for _, repo := range repos {
		commit := state.Commits[repo]
		src := state.Sources[repo]
		if _, err := s.syncRepo(repo, src, commit); err != nil {
			return fmt.Errorf(
				"sync repo %q from %q to commit %q: %s",
				repo, src, idutil.Short(commit), err,
			)
		}
	}

	if auto != nil {
		if err := installThis(s.Env, auto); err != nil {
			return err
		}
	}
	return nil
}
