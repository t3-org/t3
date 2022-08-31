package infra

import (
	"github.com/kamva/hexa"
)

const (
	UserIdCommandLine = "_command_line"
	UserIdRoot        = "_root"
)

func NewServiceUser(id string) hexa.User {
	return hexa.NewServiceUser(id, "service_user", true, []string{})
}
