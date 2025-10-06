package auth

import (
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"time"
)

func SessionStart(c *handler.Context, user UserSession) (*Session, err.Error) {
	s := &Session{
		// ID:          bson.NewObjectID(),
		Token:       GenerateToken(),
		UserID:      user.GetID(),
		Data:        user,
		Roles:       user.GetRoles(),
		Permissions: user.GetPermisions(),
		IP:          c.Request.RemoteAddr,
		Agent:       c.Request.Header.Get("user-agent"),
		Fingerprint: c.Request.Header.Get("x-fingerprint"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ExpiresAt:   expiresAt(),
	}
	return s, s.Save()
}
