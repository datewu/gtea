package router

// DefaultRoutesGroup return a new routesgroup with default router config
func DefaultRoutesGroup() *RoutesGroup {
	conf := DefaultConf()
	return &RoutesGroup{r: NewRouter(conf)}
}
