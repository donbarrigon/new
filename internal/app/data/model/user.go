package model

import (
	"context"
	"donbarrigon/new/internal/app/handler/validator"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

// ================================================================
//                  FUNCIONES PARA LA INTERFAZ
// ================================================================

func (u *User) GetID() bson.ObjectID { return u.ID }
func (u *User) Coll() string         { return "users" }

// ================================================================
//                      FUNCIONES CRUD
// ================================================================

func UserByHexID(id string) (*User, error) {
	oid, e := bson.ObjectIDFromHex(id)
	if e != nil {
		return nil, err.HexID(e.Error())
	}
	return UserByID(oid)
}

func UserByID(id bson.ObjectID) (*User, error) {
	user := &User{}
	filter := bson.D{{Key: "_id", Value: id}}
	if e := db.Mongo.Collection(user.Coll()).FindOne(context.TODO(), filter).Decode(user); e != nil {
		return nil, err.Mongo(e)
	}
	return user, nil
}

func UserByEmail(email string) (*User, error) {
	user := &User{}
	filter := bson.D{{Key: "email", Value: email}}
	if e := db.Mongo.Collection(user.Coll()).FindOne(context.TODO(), filter).Decode(user); e != nil {
		return nil, err.Mongo(e)
	}
	return user, nil
}

func UserPaginate(c *handler.Context) ([]User, error) {
	users := []User{}
	ctx := context.TODO()

	cursor, e := db.Mongo.Collection((&User{}).Coll()).Find(ctx, bson.D{}, db.PaginateFindOptions(c))
	if e != nil {
		return nil, err.Mongo(e)
	}
	if e = cursor.All(ctx, &users); e != nil {
		return nil, err.Mongo(e)
	}
	return users, nil
}

func UserCreate(dto *validator.UserStore) (*User, error) {
	bytes, e := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if e != nil {
		return nil, err.New(err.INTERNAL, "No fue posible encriptar la contraseña", e)
	}

	user := &User{
		Email:           dto.Email,
		Password:        string(bytes),
		EmailVerifiedAt: nil,
		Profile: &UserProfile{
			Avatar:      "",
			Banner:      "",
			Nickname:    dto.Nickname,
			Name:        dto.Name,
			Phone:       dto.Phone,
			Discord:     dto.Discord,
			CityID:      dto.CityID,
			Preferences: map[string]string{},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
	result, e := db.Mongo.Collection(user.Coll()).InsertOne(context.TODO(), user)
	if e != nil {
		return nil, err.Mongo(e)
	}
	user.ID = result.InsertedID.(bson.ObjectID)
	return user, nil
}

func (user *User) UpdateProfile(dto *validator.UserStore) (OriginalValues, error) {
	old := OriginalValues{}
	if user.Profile.Nickname != dto.Nickname {
		old["nickname"] = user.Profile.Nickname
		user.Profile.Nickname = dto.Nickname
	}
	if user.Profile.Name != dto.Name {
		old["name"] = user.Profile.Name
		user.Profile.Name = dto.Name
	}
	if user.Profile.Phone != dto.Phone {
		old["phone"] = user.Profile.Phone
		user.Profile.Phone = dto.Phone
	}
	if user.Profile.Discord != dto.Discord {
		old["discord"] = user.Profile.Discord
		user.Profile.Discord = dto.Discord
	}
	if user.Profile.CityID != dto.CityID {
		old["cityId"] = user.Profile.CityID
		user.Profile.CityID = dto.CityID
	}

	user.UpdatedAt = time.Now()

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: user}}
	result, e := db.Mongo.Collection(user.Coll()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return nil, err.Mongo(e)
	}
	return old, err.MongoUpdateResult(result)
}

func (user *User) UpdateEmail(dto *validator.UserUpdateEmail) (OriginalValues, error) {
	old := OriginalValues{}
	if user.Email != dto.Email {
		old["email"] = user.Email
		user.Email = dto.Email
	}

	user.UpdatedAt = time.Now()

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: user}}
	result, e := db.Mongo.Collection(user.Coll()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return nil, err.Mongo(e)
	}
	return old, err.MongoUpdateResult(result)
}

func (user *User) UpdatePassword(dto *validator.UserUpdatePassword) error {
	bytes, e := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if e != nil {
		return err.New(err.INTERNAL, "No fue posible encriptar la contraseña", e)
	}

	user.Password = string(bytes)
	user.UpdatedAt = time.Now()

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: user}}
	result, e := db.Mongo.Collection(user.Coll()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return err.Mongo(e)
	}
	return err.MongoUpdateResult(result)
}

// ================================================================
//                         PERMISSIONS
// ================================================================

type Permission struct {
	ID   bson.ObjectID `bson:"_id"  json:"id"`
	Name string        `bson:"name" json:"name"`
}

// ================================================================
//                  FUNCIONES PARA LA INTERFAZ
// ================================================================

func (p *Permission) GetID() bson.ObjectID { return p.ID }
func (p *Permission) Coll() string         { return "permissions" }

// ================================================================
//                      FUNCIONES CRUD
// ================================================================

func PermissionByID(id bson.ObjectID) (*Permission, error) {
	permission := &Permission{}
	filter := bson.D{{Key: "_id", Value: id}}
	if e := db.Mongo.Collection(permission.Coll()).FindOne(context.TODO(), filter).Decode(permission); e != nil {
		return nil, err.Mongo(e)
	}
	return permission, nil
}

func PermissionByHexID(id string) (*Permission, error) {
	oid, e := bson.ObjectIDFromHex(id)
	if e != nil {
		return nil, err.HexID(e.Error())
	}
	return PermissionByID(oid)
}

func PermissionByName(name string) (*Permission, error) {
	permission := &Permission{}
	filter := bson.D{{Key: "name", Value: name}}
	if e := db.Mongo.Collection(permission.Coll()).FindOne(context.TODO(), filter).Decode(permission); e != nil {
		return nil, err.Mongo(e)
	}
	return permission, nil
}

func PermissionCreate(name string) (*Permission, error) {
	permission := &Permission{
		Name: name,
	}
	result, e := db.Mongo.Collection(permission.Coll()).InsertOne(context.TODO(), permission)
	if e != nil {
		return nil, err.Mongo(e)
	}
	permission.ID = result.InsertedID.(bson.ObjectID)
	return permission, nil
}

func (permission *Permission) Update(dto *validator.PermissionStore) (OriginalValues, error) {

	old := OriginalValues{}
	if permission.Name != dto.Name {
		old["name"] = permission.Name
		permission.Name = dto.Name
	}

	filter := bson.D{{Key: "_id", Value: permission.ID}}
	update := bson.D{{Key: "$set", Value: permission}}
	result, e := db.Mongo.Collection(permission.Coll()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return nil, err.Mongo(e)
	}
	if e := err.MongoUpdateResult(result); e != nil {
		return nil, e
	}

	// actualizo los roles
	filter = bson.D{}
	update = bson.D{{Key: "$set", Value: bson.D{{Key: "permissions.$[elem]", Value: permission.Name}}}}
	opts := options.UpdateMany().SetArrayFilters([]any{
		bson.D{{Key: "elem", Value: old["name"]}},
	})
	resultr, e := db.Mongo.Collection("roles").UpdateMany(context.TODO(), filter, update, opts)
	if e != nil {
		return nil, err.Mongo(e)
	}
	if e := err.MongoUpdateResult(resultr); e != nil {
		return nil, e
	}

	// actualizo los usuarios

	return old, nil
}

// ================================================================
//                           ROLES
// ================================================================

type Role struct {
	ID          bson.ObjectID `bson:"_id"         json:"id"`
	Name        string        `bson:"name"        json:"name"`
	Permissions []string      `bson:"permissions" json:"permissions"`
}

// ================================================================
//                  FUNCIONES PARA LA INTERFAZ
// ================================================================

func (r *Role) GetID() bson.ObjectID { return r.ID }
func (r *Role) Coll() string         { return "roles" }

// ================================================================
//                  FUNCIONES CRUD
// ================================================================

func RoleByID(id bson.ObjectID) (*Role, error) {
	role := &Role{}
	filter := bson.D{{Key: "_id", Value: id}}
	if e := db.Mongo.Collection(role.Coll()).FindOne(context.TODO(), filter).Decode(role); e != nil {
		return nil, err.Mongo(e)
	}
	return role, nil
}

func RoleByHexID(id string) (*Role, error) {
	oid, e := bson.ObjectIDFromHex(id)
	if e != nil {
		return nil, err.HexID(e.Error())
	}
	return RoleByID(oid)
}

func RoleByName(name string) (*Role, error) {
	role := &Role{}
	filter := bson.D{{Key: "name", Value: name}}
	if e := db.Mongo.Collection(role.Coll()).FindOne(context.TODO(), filter).Decode(role); e != nil {
		return nil, err.Mongo(e)
	}
	return role, nil
}

func RoleCreate(dto *validator.RoleStore) (*Role, error) {
	role := &Role{
		Name:        dto.Name,
		Permissions: []string{},
	}
	result, e := db.Mongo.Collection(role.Coll()).InsertOne(context.TODO(), role)
	if e != nil {
		return nil, err.Mongo(e)
	}
	role.ID = result.InsertedID.(bson.ObjectID)
	return role, nil
}

func (role *Role) Update(dto *validator.RoleStore) (OriginalValues, error) {

	old := OriginalValues{}
	if role.Name != dto.Name {
		old["name"] = role.Name
		role.Name = dto.Name
	}

	filter := bson.D{{Key: "_id", Value: role.ID}}
	result, e := db.Mongo.Collection(role.Coll()).UpdateOne(context.TODO(), filter, role)
	if e != nil {
		return nil, err.Mongo(e)
	}
	return old, err.MongoUpdateResult(result)
}
