package auth

import (
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"net"
	"os"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SessionUser interface {
	GetID() bson.ObjectID
	Can(permission string) bool
	HasRole(role string) bool
}

func SessionStart(c *handler.Context, user SessionUser) error {
	host, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	tk, e := GenerateHexToken()
	if e != nil {
		return e
	}
	s := &Session{
		//ID:          bson.NewObjectID(),
		Token:       tk,
		User:        user,
		IP:          host,
		Agent:       c.Request.Header.Get("user-agent"),
		Lang:        c.Request.Header.Get("accept-language"),
		Fingerprint: c.Request.Header.Get("x-fingerprint"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   expiresAt(),
		writer:      c.Writer,
		request:     c.Request,
	}
	s.SetCookie()
	c.Auth = s
	return s.Save()
}

func expiresAt() time.Time {
	return time.Now().Add(time.Duration(config.SessionLifetime) * time.Minute)
}

func GetSession(c *handler.Context) (*Session, error) {
	var token string
	cookie, e := c.Request.Cookie("session")
	if e != nil {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader != "" {
			token = strings.TrimSpace(strings.TrimPrefix(authHeader, "bearer "))
		}

		if token == "" {
			return nil, err.Unauthorized(e)
		}
	} else {
		token = cookie.Value
	}

	s, e := GetSessionByToken(token)
	if e != nil {
		return nil, e
	}
	s.writer = c.Writer
	s.request = c.Request

	s.SetCookie()
	if e := s.Refresh(); e != nil {
		return nil, e
	}
	return s, nil
}

func GetSessionByToken(token string) (*Session, error) {
	s := &Session{}
	path, filename := fileSession(token)
	info, e := os.Stat(path + filename)
	if e == nil && !info.IsDir() {
		encoded, e := os.ReadFile(path + filename)
		if e != nil {
			return nil, err.New(err.FORBIDDEN, "No ha iniciado session", e)
		}
		if e := msgpack.Unmarshal(encoded, &s); e != nil {
			return nil, err.New(err.FORBIDDEN, "No ha iniciado session", e)
		}
	}
	return s, nil
}

func GetSessionsByUserID(hex string) ([]*Session, error) {
	sessions := []*Session{}
	tokens := map[string]time.Time{}
	path, filename := fileUserIndex(hex)
	info, e := os.Stat(path + filename)
	if e == nil && !info.IsDir() {
		encoded, e := os.ReadFile(path + filename)
		if e != nil {
			return nil, err.New(err.FORBIDDEN, "No ha iniciado session", e)
		}
		if e := msgpack.Unmarshal(encoded, &tokens); e != nil {
			return nil, err.New(err.FORBIDDEN, "No ha iniciado session", e)
		}
	}

	for token := range tokens {
		s, e := GetSessionByToken(token)
		if e != nil {
			return nil, e
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func fileSession(token string) (string, string) {
	return "tmp/sessions/" + token[:3] + "/" + token[3:6] + "/", token[6:]
}

func fileUserIndex(hex string) (string, string) {
	return "tmp/sessions/index/" + hex[:4] + "/" + hex[4:8] + "/", hex[8:]
}
