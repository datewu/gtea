package router

import (
	"errors"

	"github.com/julienschmidt/httprouter"
)

// NewRoutesGroup return a new routesgroup
func NewRoutesGroup(conf *Config) (*RoutesGroup, error) {
	if conf == nil {
		return nil, errors.New("no router config provided")
	}
	r := bag{
		rt:     httprouter.New(),
		config: conf,
	}
	return &RoutesGroup{r: &r}, nil
}

// DefaultRoutesGroup return a new routesgroup with default router config
func DefaultRoutesGroup() *RoutesGroup {
	r := bag{
		rt:     httprouter.New(),
		config: DefaultConf(),
	}
	return &RoutesGroup{r: &r}
}
