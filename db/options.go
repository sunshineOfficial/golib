package db

type optionsHandler struct {
	traces bool
}

type Option interface {
	apply(handler optionsHandler) optionsHandler
}

type OptionFunc func(handler optionsHandler) optionsHandler

func (o OptionFunc) apply(handler optionsHandler) optionsHandler {
	return o(handler)
}

func WithTraces() Option {
	return OptionFunc(func(handler optionsHandler) optionsHandler {
		handler.traces = true
		return handler
	})
}

func applyOptions(options ...Option) optionsHandler {
	var handler optionsHandler
	for _, o := range options {
		handler = o.apply(handler)
	}

	return handler
}
