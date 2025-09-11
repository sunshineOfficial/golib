package gorouter

type contextOptionHolder struct {
	userInfo bool
	traces   bool
}

type ContextOption interface {
	apply(options contextOptionHolder) contextOptionHolder
}

type contextOptionFunc func(o contextOptionHolder) contextOptionHolder

func (f contextOptionFunc) apply(o contextOptionHolder) contextOptionHolder {
	return f(o)
}

func WithUserInfo() ContextOption {
	return contextOptionFunc(func(o contextOptionHolder) contextOptionHolder {
		o.userInfo = true
		return o
	})
}

func WithTraces() ContextOption {
	return contextOptionFunc(func(o contextOptionHolder) contextOptionHolder {
		o.traces = true
		return o
	})
}
