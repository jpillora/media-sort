// +build !linux,!darwin,!freebsd

package mediasort

import "errors"

const canSysMove = false

func sysMove(src, dst string) error {
	return errors.New("system move unsupported")
}
