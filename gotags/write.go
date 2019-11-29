package gotags

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

// Write writes the tag file out.
func Write(files []string, outputFile string) error {
	const relative = true
	baseDir, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("get current working dir: %s", err)
	}

	tags := []Tag{}
	for _, file := range files {
		ts, err := Parse(file, relative, baseDir)
		if err != nil {
			return fmt.Errorf("parse: %s", err)
		}
		tags = append(tags, ts...)
	}

	output := createMetaTags()
	for _, tag := range tags {
		output = append(output, tag.String())
	}

	sort.Strings(output)

	out, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("create output: %s", err)
	}
	defer out.Close()

	for _, s := range output {
		fmt.Fprintln(out, s)
	}

	if err := out.Sync(); err != nil {
		return fmt.Errorf("flush output file: %s", err)
	}
	return nil
}

// createMetaTags returns a list of meta tags.
func createMetaTags() []string {
	// Contants used for the meta tags
	const (
		version     = "0.1"
		name        = "gotags_shanhuio"
		url         = "https://github.com/jstemmer/gotags"
		authorName  = "Joel Stemmer"
		authorEmail = "stemmertech@gmail.com"
	)

	const sorted = 1
	return []string{
		"!_TAG_FILE_FORMAT\t2",
		fmt.Sprintf("!_TAG_FILE_SORTED\t%d\t/0=unsorted, 1=sorted/", sorted),
		fmt.Sprintf("!_TAG_PROGRAM_AUTHOR\t%s\t/%s/", authorName, authorEmail),
		fmt.Sprintf("!_TAG_PROGRAM_NAME\t%s", name),
		fmt.Sprintf("!_TAG_PROGRAM_URL\t%s", url),
		fmt.Sprintf(
			"!_TAG_PROGRAM_VERSION\t%s\t/%s/", version, runtime.Version(),
		),
	}
}
