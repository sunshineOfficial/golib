package gorouter

type Plugin interface {
	BasePath() string
	Register(router *Router)
}
