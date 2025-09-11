package plugin

import (
	"net/http"
	"net/http/pprof"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

type PProf struct{}

func NewPProf() *PProf {
	return &PProf{}
}

func (p PProf) BasePath() string {
	return "/debug/pprof"
}

func (p PProf) Register(router *gorouter.Router) {
	router.HandleGet("", gorouter.WrapStdLibFunc(pprof.Index))
	router.HandleGet("/", gorouter.WrapStdLibFunc(pprof.Index))
	router.HandleGet("/cmdline", gorouter.WrapStdLibFunc(pprof.Cmdline))
	router.HandleGet("/profile", gorouter.WrapStdLibFunc(pprof.Profile))

	router.HandleGet("/heap", loadSwHandler("heap"))
	router.HandleGet("/goroutine", loadSwHandler("goroutine"))
	router.HandleGet("/block", loadSwHandler("block"))
	router.HandleGet("/threadcreate", loadSwHandler("threadcreate"))

	router.Handle("/symbol", gorouter.WrapStdLibFunc(pprof.Symbol), http.MethodGet, http.MethodPost)

	router.HandleGet("/trace", loadSwHandler("trace"))
	router.HandleGet("/mutex", loadSwHandler("mutex"))
}

func loadSwHandler(name string) gorouter.Handler {
	return gorouter.WrapStdLib(pprof.Handler(name))
}
