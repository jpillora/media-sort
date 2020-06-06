// +build linux darwin freebsd

package mediasort

import (
	"errors"
	"io"
	"os"
	"strings"
)

func move(src, dst string) (err error) {
	err = os.Rename(src, dst)
	// cross device move
	if err != nil && strings.Contains(err.Error(), "cross-device") {
		if err = copy(src, dst); err != nil {
			return err
		}
		err = os.Remove(src)
	}
	return
}

func copy(src, dst string) (err error) {
	srcFile, err := os.Open(src)
	defer srcFile.Close()
	if err != nil {
		return
	}
	dstFile, err := os.Open(dst)
	defer dstFile.Close()
	if err != nil {
		return
	}
	_, err = io.Copy(dstFile, srcFile)
	return
}

func link(src, dst string, linkType linkType) (err error) {
	switch linkType {
	case hardLink:
		err = os.Link(src, dst)
	case symLink:
		err = os.Link(src, dst)
	default:
		err = errors.New("wrong link type, please open an issue")
	}
	return
}
