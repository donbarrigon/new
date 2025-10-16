package middleware

import (
	"donbarrigon/new/internal/utils/auth"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"donbarrigon/new/internal/utils/lang"
	"net"
)

func Auth(next handler.ControllerFun) handler.ControllerFun {

	return func(c *handler.Context) {
		score := 0
		s, e := auth.GetSession(c.Writer, c.Request)
		if e != nil {
			c.ResponseError(err.Unauthorized(e))
			return
		}

		if s.IsExpired() {
			if e := s.Destroy(); e != nil {
				c.ResponseError(err.New(err.UNAUTHORIZED, "La sesión ha expirado", e))
				return
			}
			c.ResponseError(err.New(err.UNAUTHORIZED, "La sesión ha expirado", nil))
			return
		}

		ipr, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
		if s.IP == ipr {
			score += 2
		}

		agent := c.Request.Header.Get("user-agent")
		if s.Agent == agent {
			score += 2
		}

		if s.Fingerprint == c.Request.Header.Get("x-fingerprint") {
			score += 2
		}

		if s.Lang == c.Request.Header.Get("accept-language") {
			score += 1
		}

		if score < 5 {
			if e := s.Destroy(); e != nil {
				c.ResponseError(err.Unauthorized(e))
				return
			}
			c.ResponseError(err.Unauthorized(lang.T(c.Lang(), "La sesión es inválida", nil)))
			return
		}

		c.Auth = s

		next(c)
	}
}
