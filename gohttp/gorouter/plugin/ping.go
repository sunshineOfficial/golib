package plugin

import (
	"net/http"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

type Ping struct {
}

func NewPing() *Ping {
	return &Ping{}
}

func (p *Ping) BasePath() string {
	return "/debug/ping"
}

func (p *Ping) Register(router *gorouter.Router) {
	router.HandleGet("", func(c gorouter.Context) error {
		return c.WriteText(http.StatusOK, "pong")
	})
}
