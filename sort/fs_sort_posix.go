// +build linux darwin freebsd

package mediasort

import (
	"errors"
	"os/exec"
)

const canSysMove = true

func sysMove(src, dst string) error {
	out, err := exec.Command("mv", src, dst).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	return nil
}
