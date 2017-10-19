package sml

import (
	"fmt"
)

func untrack(server string, args []string) error {
	flags := newFlags()
	flags.Parse(args)
	repos := flags.Args()
	profileName := defaultProfileName

	if len(repos) == 0 {
		return fmt.Errorf("no repo specified")
	}

	profile, err := loadProfile(profileName)
	if err != nil {
		return err
	}

	m := profile.trackingMap()
	for _, repo := range repos {
		if m[repo] {
			delete(m, repo)
		}
	}
	profile.setTrackingFromMap(m)
	return saveProfile(profileName, profile)
}
