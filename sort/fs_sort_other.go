// +build !linux,!darwin,!freebsd

package mediasort

import "errors"

func move(src, dst string) error                    { return errors.New("not implemented yet") }
func copy(src, dst string) error                    { return errors.New("not implemented yet") }
func link(src, dst string, linkType linkType) error { return errors.New("not implemented yet") }
