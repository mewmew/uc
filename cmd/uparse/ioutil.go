package main

import (
	"io/ioutil"
	"os"

	"github.com/mewkiz/pkg/errutil"
)

// readFile reads from the given file, or standard input if path is "-", and
// returns its file contents.
func readFile(path string) ([]byte, error) {
	if path == "-" {
		// Read from standard input.
		buf, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, errutil.Err(err)
		}
		return buf, nil
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return buf, nil
}
