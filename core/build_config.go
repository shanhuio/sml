package core

// BuildConfig specifies the global config for the builder.
type BuildConfig struct {
	// TestDataWhiteList specifies testdata directories that we will still
	// treat as Go language packages rather than data.
	TestDataWhiteList []string `json:",omitempty"`

	// PkgBlackList specifies directories that we are not building.
	PkgBlackList []string `json:",omitempty"`

	// RepoFixes specifies repo mirrors required to make things build.
	RepoFixes map[string]string `json:",omitempty"`
}
