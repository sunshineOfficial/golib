package gorouter

type Middleware func(c Context, h Handler) Handler
