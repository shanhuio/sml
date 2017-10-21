package sml

import (
	"fmt"

	"shanhu.io/misc/httputil"
)

func track(server string, args []string) error {
	flags := newFlags()
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
	err := c.JSONCall("/api/tracking:missing", repos, &missing)
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

	m := profile.trackingMap()
	for _, repo := range repos {
		m[repo] = true
	}
	profile.setTrackingFromMap(m)
	return saveProfile(profileName, profile)
}
