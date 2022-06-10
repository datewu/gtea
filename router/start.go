package router

import (
	"errors"

	"github.com/julienschmidt/httprouter"
)

// NewRouterGroup return a new routergroup
func NewRouterGroup(conf *Config) (*RoutesGroup, error) {
	if conf == nil {
		return nil, errors.New("no router config provided")
	}
	r := bag{
		rt:     httprouter.New(),
		config: conf,
	}
	return &RoutesGroup{r: &r}, nil
}

// DefaultRouterGroup return a new routergroup with default router config
func DefaultRouterGroup() *RoutesGroup {
	r := bag{
		rt:     httprouter.New(),
		config: DefaultConf(),
	}
	return &RoutesGroup{r: &r}
}
