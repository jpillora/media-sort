// +build !linux,!darwin,!freebsd

package mediasort

import "errors"

func canSysMove() bool {
	return false
}

func sysMove(src, dst string) error {
	return errors.New("system move unsupported")
}
