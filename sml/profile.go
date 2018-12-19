package sml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"shanhu.io/sml/goenv"
)

// Profile is the config file saved in $GOPATH/src/sml.json
type Profile struct {
	Tracking []string
}

func (p *Profile) trackingMap() map[string]bool {
	m := make(map[string]bool)
	for _, repo := range p.Tracking {
		m[repo] = true
	}
	return m
}

func (p *Profile) setTrackingFromMap(m map[string]bool) {
	var tracking []string
	for repo := range m {
		tracking = append(tracking, repo)
	}
	sort.Strings(tracking)
	p.Tracking = tracking
}

const defaultProfileName = "sml"

func isDir(p string) (bool, error) {
	info, err := os.Stat(p)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}

func checkGOPATH() (string, error) {
	gopath, err := goenv.GOPATH()
	if err != nil {
		return "", err
	}

	ok, err := isDir(gopath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("GOPATH %q not exist", gopath)
	}
	if err != nil {
		return "", fmt.Errorf("reading GOPATH %q: %s", gopath, err)
	}
	if !ok {
		return "", fmt.Errorf("GOPATH %q is not a directory", gopath)
	}

	return gopath, nil
}

func profilePath(name string) (string, error) {
	gopath, err := checkGOPATH()
	if err != nil {
		return "", err
	}

	base := fmt.Sprintf("%s.json", name)
	return filepath.Join(gopath, "sml", base), nil
}

func loadProfile(name string) (*Profile, error) {
	p, err := profilePath(name)
	if err != nil {
		return nil, err
	}

	bs, err := ioutil.ReadFile(p)
	if os.IsNotExist(err) {
		return new(Profile), nil
	}

	c := new(Profile)
	if err := json.Unmarshal(bs, c); err != nil {
		return nil, fmt.Errorf("profile %q: %s", p, err)
	}
	return c, nil
}

func saveProfile(name string, prof *Profile) error {
	p, err := profilePath(name)
	if err != nil {
		return err
	}

	bs, err := json.MarshalIndent(prof, "", "  ")
	if err != nil {
		return err
	}

	// make sure the directory exists.
	if dir := filepath.Dir(p); dir != "" {
		if err := os.MkdirAll(dir, os.ModeDir|0755); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(p, bs, 0644)
}
