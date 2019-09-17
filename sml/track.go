package sml

import (
	"fmt"
	"os"

	"shanhu.io/base/httputil"
)

func track(server string, args []string) error {
	flags := newFlags()
	sync := flags.Bool("sync", true, "synchronize after tracking a new repo")
	flags.Parse(args)
	repos := flags.Args()
	profileName := defaultProfileName

	if len(repos) == 0 {
		profile, err := loadProfile(profileName)
		if err != nil {
			return err
		}

		for _, repo := range profile.Tracking {
			fmt.Println(repo)
		}
		return nil
	}

	c := httputil.NewClient(server)
	var missing []string
	err := c.JSONCall("/api/repos/not-tracking", repos, &missing)
	if err != nil {
		return fmt.Errorf("cannot confirm if trackable: %s", err)
	}

	if len(missing) > 0 {
		fmt.Println("error: these repos not tracked by smallrepo:")
		for _, repo := range missing {
			fmt.Printf("  %s\n", repo)
		}
		return fmt.Errorf("cannot track all the repos requested")
	}

	profile, err := loadProfile(profileName)
	if err != nil {
		return err
	}

	changed := false
	m := profile.trackingMap()
	for _, repo := range repos {
		if m[repo] {
			fmt.Fprintf(os.Stderr, "%q is already tracked\n", repo)
		} else {
			m[repo] = true
			changed = true
		}
	}
	profile.setTrackingFromMap(m)
	if err := saveProfile(profileName, profile); err != nil {
		return err
	}

	if changed && *sync {
		return doSync(server, profile)
	}
	return nil
}
