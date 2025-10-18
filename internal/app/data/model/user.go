package model

import (
	"donbarrigon/new/internal/utils/db"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID              bson.ObjectID   `bson:"_id"               json:"id"`
	Email           string          `bson:"email"             json:"email"` // unico
	EmailVerifiedAt *time.Time      `bson:"email_verified_at" json:"emailVerifiedAt"`
	Password        string          `bson:"password"          json:"-"`
	Roles           map[string]bool `bson:"roles"             json:"roles"`
	Permissions     map[string]bool `bson:"permissions"       json:"permissions"`
	Profile         *UserProfile    `bson:"profile"           json:"profile"`
	CreatedAt       time.Time       `bson:"created_at"        json:"createdAt"`
	UpdatedAt       time.Time       `bson:"updated_at"        json:"updatedAt"`
	DeletedAt       *time.Time      `bson:"deleted_at"        json:"deletedAt"`
	db.Odm          `bson:"-" json:"-"`
}

type UserProfile struct {
	Avatar      string            `bson:"avatar"            json:"avatar"`
	Banner      string            `bson:"banner"            json:"banner"`
	Nickname    string            `bson:"nickname"          json:"nickname"` // unico
	Name        string            `bson:"name"              json:"name"`
	Phone       string            `bson:"phone,omitempty"   json:"phone,omitempty"`
	Discord     string            `bson:"discord,omitempty" json:"discord,omitempty"`
	CityID      bson.ObjectID     `bson:"city_id"           json:"cityId"`
	Preferences map[string]string `bson:"preferences"       json:"preferences"`
}

func (u *User) CollectionName() string { return "users" }
func (u *User) GetID() bson.ObjectID   { return u.ID }
func (u *User) SetID(id bson.ObjectID) { u.ID = id }

func (u *User) BeforeCreate() error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate() error {
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Can(permission string) bool {
	return u.Permissions[permission]
}

func (u *User) HasRole(role string) bool {
	return u.Roles[role]
}

// ================================================================
// roles y permisos
// ================================================================

type Permissions struct {
	ID   bson.ObjectID `bson:"_id"  json:"id"`
	Name string        `bson:"name" json:"name"`
}

func (p *Permissions) CollectionName() string { return "permisions" }
func (p *Permissions) GetID() bson.ObjectID   { return p.ID }
func (p *Permissions) SetID(id bson.ObjectID) { p.ID = id }

type Roles struct {
	ID          bson.ObjectID `bson:"_id"         json:"id"`
	Name        string        `bson:"name"        json:"name"`
	Permissions []string      `bson:"permissions" json:"permissions"`
}

func (r *Roles) CollectionName() string { return "roles" }
func (r *Roles) GetID() bson.ObjectID   { return r.ID }
func (r *Roles) SetID(id bson.ObjectID) { r.ID = id }
