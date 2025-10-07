package auth

import (
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/err"
	"net/http"
	"os"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

func SessionStart(w http.ResponseWriter, r *http.Request, user UserSession) (*Session, err.Error) {
	s := &Session{
		// ID:          bson.NewObjectID(),
		Token:       GenerateToken(),
		UserID:      user.GetID(),
		Data:        user,
		Roles:       user.GetRoles(),
		Permissions: user.GetPermisions(),
		IP:          r.RemoteAddr,
		Agent:       r.Header.Get("user-agent"),
		Fingerprint: r.Header.Get("x-fingerprint"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   expiresAt(),
		writer:      w,
		request:     r,
	}
	s.SetCookie()
	return s, s.Save()
}

func expiresAt() time.Time {
	return time.Now().Add(time.Duration(config.SessionLifetime) * time.Minute)
}

func GetSession(w http.ResponseWriter, r *http.Request) (*Session, err.Error) {
	cookie, e := r.Cookie("session")
	if e != nil {
		return nil, err.New(err.FORBIDDEN, "No ha iniciado session", e)
	}
	s, he := GetSessionByToken(cookie.Value)
	if he != nil {
		return nil, he
	}
	s.writer = w
	s.request = r

	s.SetCookie()
	s.Refresh()
	return s, nil
}

func GetSessionByToken(token string) (*Session, err.Error) {
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

func GetSessionsByUserID(hex string) ([]*Session, err.Error) {
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
