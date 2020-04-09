package lib

import (
	"errors"
	"os/exec"
)

func killPpp() error {
	command := "ps -ef|pgrep ppp"
	cmd := exec.Command("bash", "-c", "kill -9 $("+command+")")
	if _, err := cmd.Output(); err != nil {
		return errors.New("ppp已经呗停止了")
	}
	return nil
}