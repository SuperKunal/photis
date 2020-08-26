package middleware

import (
	"errors"
	"github.com/savsgio/atreugo/v11"
	"log"
	"net/http"
	"os"
)

type Middleware struct {}

func InitMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) PanicHandler(ctx *atreugo.RequestCtx, err interface{}) {
	var e error
	if err != nil {
		switch t := err.(type) {
		case string:
			e = errors.New(t)
		case error:
			e = t
		default:
			e = errors.New("Unknown error")
		}
		f, _ := os.OpenFile("err.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		defer f.Close()
		log.SetOutput(f)
		log.Println(e)
	}
	ctx.SetStatusCode(http.StatusInternalServerError)
}