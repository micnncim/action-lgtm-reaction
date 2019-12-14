package actions

import (
	"github.com/kelseyhightower/envconfig"
)

type Input struct {
	Trigger  string
	Override bool
	Source   string
}

func GetInput() (i Input) {
	envconfig.MustProcess("INPUT", &i)
	return i
}
