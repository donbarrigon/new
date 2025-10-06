package auth

import (
	"donbarrigon/new/internal/utils/err"
	"os"
	"sync"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserSession interface {
	GetID() bson.ObjectID
	GetPermisions() map[string]bool
	GetRoles() map[string][]string
}

type Session struct {
	// ID          bson.ObjectID
	Token       string
	UserID      bson.ObjectID
	Data        any
	Roles       map[string][]string
	Permissions map[string]bool
	IP          string
	Agent       string
	Fingerprint string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
}

var muSession = sync.Map{}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) Save() err.Error {
	echan := make(chan err.Error, 2)
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

func (s *Session) Refresh() err.Error {
	muToken := s.muToken()
	muToken.Lock()
	defer muToken.Unlock()
	s.UpdatedAt = time.Now()
	s.ExpiresAt = expiresAt()
	return s.saveFileSession()
}

func (s *Session) Destroy() err.Error {
	muToken := s.muToken()
	muToken.Lock()
	defer muToken.Unlock()
	if e := s.deleteFileSession(); e != nil {
		return e
	}
	muUser := s.muUser()
	muUser.Lock()
	defer muUser.Unlock()
	return s.removeFileUserIndex()
}

func (s *Session) saveFileSession() err.Error {
	path, filename := s.fileSession()
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

func (s *Session) saveFileUserIndex(data map[string]time.Time) err.Error {
	path, filename := s.fileUserIndex()
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

func (s *Session) addFileUserIndex() err.Error {
	data, he := s.readFileUserIndex()
	if he != nil {
		return he
	}
	data[s.Token] = s.CreatedAt
	return s.saveFileUserIndex(data)
}

func (s *Session) deleteFileSession() err.Error {
	path, fileName := s.fileSession()
	if e := os.Remove(path + fileName); e != nil {
		return err.New(err.INTERNAL, "No se cerró la sesion", e)
	}
	return nil
}

func (s *Session) removeFileUserIndex() err.Error {
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

func (s *Session) deleteFileUserIndex() err.Error {
	path, fileName := s.fileUserIndex()
	if e := os.Remove(path + fileName); e != nil {
		return err.New(err.INTERNAL, "No se eliminó el indice de la sesion", e)
	}
	return nil
}

func (s *Session) readFileUserIndex() (map[string]time.Time, err.Error) {
	data := map[string]time.Time{}
	path, filename := s.fileUserIndex()
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
	mu, _ := muSession.LoadOrStore(s.UserID.Hex(), &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func (s *Session) muToken() *sync.Mutex {
	mu, _ := muSession.LoadOrStore(s.Token, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func (s *Session) fileSession() (string, string) {
	return "tmp/sessions/" + s.Token[:3] + "/" + s.Token[3:6] + "/", s.Token[6:]
}

func (s *Session) fileUserIndex() (string, string) {
	hex := s.UserID.Hex()
	return "tmp/sessions/index/" + hex[:4] + "/" + hex[4:8] + "/", hex[8:]
}
