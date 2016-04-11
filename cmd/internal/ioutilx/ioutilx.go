// Package ioutilx implements extended input/output utility functions.
package ioutilx

import (
	"io/ioutil"
	"os"

	"github.com/mewkiz/pkg/errutil"
)

// ReadFile reads from the given file, or standard input if path is "-", and
// returns its file contents.
func ReadFile(path string) ([]byte, error) {
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
