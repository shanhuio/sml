package core

// BuildConfig specifies the global config for the builder.
type BuildConfig struct {
	// TestDataWhiteList specifies testdata directories that we will still
	// treat as Go language packages rather than data.
	TestDataWhiteList []string

	// PkgBlackList specifies directories that we are not building.
	PkgBlackList []string
}
