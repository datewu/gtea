package handler

import (
	"errors"

	"github.com/julienschmidt/httprouter"
)

// NewRouterGroup return a new routergroup
func NewRouterGroup(conf *Config) (*RouterGroup, error) {
	if conf == nil {
		return nil, errors.New("no router config provided")
	}
	r := router{
		router: httprouter.New(),
		config: conf,
	}
	return &RouterGroup{r: &r}, nil
}

// DefaultRouterGroup return a new routergroup with default router config
func DefaultRouterGroup() *RouterGroup {
	r := router{
		router: httprouter.New(),
		config: DefaultConf(),
	}
	return &RouterGroup{r: &r}
}
