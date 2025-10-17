package auth

import (
	"donbarrigon/new/internal/utils/config"
	"donbarrigon/new/internal/utils/err"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type Session struct {
	//ID          bson.ObjectID
	Token       string
	User        SessionUser
	IP          string
	Agent       string
	Lang        string
	Fingerprint string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
	writer      http.ResponseWriter
	request     *http.Request
}

var muSession = sync.Map{}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Save() error {
	echan := make(chan error, 2)
	go func() {
		muToken := s.muToken()
		muToken.Lock()
		defer muToken.Unlock()
		echan <- s.saveFileSession()
	}()
	go func() {
		muUser := s.muUser()
		muUser.Lock()
		defer muUser.Unlock()
		echan <- s.addFileUserIndex()
	}()

	if e := <-echan; e != nil {
		return e
	}
	if e := <-echan; e != nil {
		return e
	}
	return nil
}

func (s *Session) Refresh() error {
	muToken := s.muToken()
	muToken.Lock()
	defer muToken.Unlock()
	s.UpdatedAt = time.Now()
	s.ExpiresAt = expiresAt()
	return s.saveFileSession()
}

func (s *Session) Destroy() error {
	muToken := s.muToken()
	muToken.Lock()
	defer muToken.Unlock()
	if e := s.deleteFileSession(); e != nil {
		return e
	}
	s.ClearCookie()

	muUser := s.muUser()
	muUser.Lock()
	defer muUser.Unlock()
	return s.removeFileUserIndex()
}

func (s *Session) SetCookie() {
	http.SetCookie(s.writer, &http.Cookie{
		Name:     "session",
		Value:    s.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   config.ServerHttpsEnabled,
		SameSite: http.SameSiteLaxMode,
		Expires:  s.ExpiresAt,
	})
}

func (s *Session) ClearCookie() {
	http.SetCookie(s.writer, &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   config.ServerHttpsEnabled,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
	})
}

func (s *Session) Can(permission string) bool {
	return s.User.Can(permission)
}

func (s *Session) HasRole(role string) bool {
	return s.User.HasRole(role)
}

func (s *Session) saveFileSession() error {
	path, filename := fileSession(s.Token)
	encoded, e := msgpack.Marshal(s)
	if e != nil {
		return err.New(err.INTERNAL, "No se pudo codificar la sesion", e)
	}
	if e := os.MkdirAll(path, 0755); e != nil {
		return err.New(err.INTERNAL, "No se pudo crear el directorio de sesion", e)
	}
	if e := os.WriteFile(path+filename, encoded, 0644); e != nil {
		return err.New(err.INTERNAL, "No se pudo crear la sesiones", e)
	}
	return nil
}

func (s *Session) saveFileUserIndex(data map[string]time.Time) error {
	path, filename := fileUserIndex(s.User.GetID().Hex())
	encoded, e := msgpack.Marshal(data)
	if e != nil {
		return err.New(err.INTERNAL, "No se pudo codificar el indice de sessiones", e)
	}
	if e := os.MkdirAll(path, 0755); e != nil {
		return err.New(err.INTERNAL, "No se pudo crear el indice de sesiones", e)
	}
	if e := os.WriteFile(path+filename, encoded, 0644); e != nil {
		return err.New(err.INTERNAL, "No se pudo crear el indice de session", e)
	}
	return nil
}

func (s *Session) addFileUserIndex() error {
	data, he := s.readFileUserIndex()
	if he != nil {
		return he
	}
	data[s.Token] = s.CreatedAt
	return s.saveFileUserIndex(data)
}

func (s *Session) deleteFileSession() error {
	path, fileName := fileSession(s.Token)
	if e := os.Remove(path + fileName); e != nil {
		return err.New(err.INTERNAL, "No se cerró la sesion", e)
	}
	return nil
}

func (s *Session) removeFileUserIndex() error {
	data, he := s.readFileUserIndex()
	if he != nil {
		return he
	}
	delete(data, s.Token)
	if len(data) == 0 {
		return s.deleteFileUserIndex()
	}
	return s.saveFileUserIndex(data)
}

func (s *Session) deleteFileUserIndex() error {
	path, fileName := fileUserIndex(s.User.GetID().Hex())
	if e := os.Remove(path + fileName); e != nil {
		return err.New(err.INTERNAL, "No se eliminó el indice de la sesion", e)
	}
	return nil
}

func (s *Session) readFileUserIndex() (map[string]time.Time, error) {
	data := map[string]time.Time{}
	path, filename := fileUserIndex(s.User.GetID().Hex())
	info, e := os.Stat(path + filename)
	if e == nil && !info.IsDir() {
		encoded, e := os.ReadFile(path + filename)
		if e != nil {
			return nil, err.New(err.INTERNAL, "No se pudo leer el indice de sessiones", e)
		}
		if e := msgpack.Unmarshal(encoded, &data); e != nil {
			return nil, err.New(err.INTERNAL, "No se pudo decodificar el indice de sessiones", e)
		}
	}
	return data, nil
}

func (s *Session) muUser() *sync.Mutex {
	mu, _ := muSession.LoadOrStore(s.User.GetID().Hex(), &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func (s *Session) muToken() *sync.Mutex {
	mu, _ := muSession.LoadOrStore(s.Token, &sync.Mutex{})
	return mu.(*sync.Mutex)
}
