package factory

import (
	"errors"
	"os"
	"os/exec"

	"github.com/Sirupsen/logrus"
)

type Runc struct {
	name string
	ID   string
}

func (this *Runc) init() {
	this.name = ""
	this.ID = ""
}

func (this *Runc) SetRT(runtime string) {
	this.name = "runc"
}

func (this *Runc) GetRT() string {
	return "runc"
}

func (this *Runc) GetRTID() string {
	return this.ID
}

func (this *Runc) StartRT(specDir string) (string, error) {
	logrus.Debugf("Launcing runtime")

	cmd := exec.Command("runc", "start")
	cmd.Dir = specDir
	cmd.Stdin = os.Stdin
	out, err := cmd.CombinedOutput()

	logrus.Debugf("Command done")
	if err != nil {
		return string(out), errors.New(string(out) + err.Error())
	}

	return string(out), nil
}

func (this *Runc) StopRT(id string) error {
	return nil
}
