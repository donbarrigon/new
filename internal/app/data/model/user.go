package model

import (
	"context"
	"crypto/rand"
	"donbarrigon/new/internal/app/data/validator"
	"donbarrigon/new/internal/utils/db"
	"donbarrigon/new/internal/utils/err"
	"donbarrigon/new/internal/utils/handler"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              bson.ObjectID   `bson:"_id"               json:"id"`
	Email           string          `bson:"email"             json:"email"` // unico
	EmailVerifiedAt *time.Time      `bson:"email_verified_at" json:"emailVerifiedAt"`
	Password        string          `bson:"password"          json:"-"`
	Profile         *UserProfile    `bson:"profile"           json:"profile"`
	Roles           map[string]bool `bson:"roles"             json:"roles"`
	Permissions     map[string]bool `bson:"permissions"       json:"permissions"`
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

func (u *User) Can(permission string) bool {
	return u.Permissions[permission]
}

func (u *User) HasRole(role string) bool {
	return u.Roles[role]
}

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
		Roles:       map[string]bool{"user": true},
		Permissions: map[string]bool{"login": true},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}
	result, e := db.Mongo.Collection(user.Coll()).InsertOne(context.TODO(), user)
	if e != nil {
		return nil, err.Mongo(e)
	}
	user.ID = result.InsertedID.(bson.ObjectID)
	return user, nil
}

func (user *User) UpdateProfile(dto *validator.UserUpdateProfile) (*Changes, error) {
	changes := NewChanges()
	if user.Profile.Nickname != dto.Nickname {
		changes.Old["nickname"] = user.Profile.Nickname
		changes.New["nickname"] = dto.Nickname
		user.Profile.Nickname = dto.Nickname
	}
	if user.Profile.Name != dto.Name {
		changes.Old["name"] = user.Profile.Name
		changes.New["name"] = dto.Name
		user.Profile.Name = dto.Name
	}
	if user.Profile.Phone != dto.Phone {
		changes.Old["phone"] = user.Profile.Phone
		changes.New["phone"] = dto.Phone
		user.Profile.Phone = dto.Phone
	}
	if user.Profile.Discord != dto.Discord {
		changes.Old["discord"] = user.Profile.Discord
		changes.New["discord"] = dto.Discord
		user.Profile.Discord = dto.Discord
	}
	if user.Profile.CityID != dto.CityID {
		changes.Old["city_id"] = user.Profile.CityID
		changes.New["city_id"] = dto.CityID
		user.Profile.CityID = dto.CityID
	}

	user.UpdatedAt = time.Now()

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: user}}
	result, e := db.Mongo.Collection(user.Coll()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return nil, err.Mongo(e)
	}
	if e := err.MongoUpdateResult(result); e != nil {
		return nil, e
	}
	return changes, nil
}

func (user *User) UpdateEmail(dto *validator.UserUpdateEmail) (*Changes, error) {
	changes := NewChanges()
	if user.Email != dto.Email {
		changes.Old["email"] = user.Email
		changes.New["email"] = dto.Email
		user.Email = dto.Email
	}

	user.UpdatedAt = time.Now()

	filter := bson.D{{Key: "_id", Value: user.ID}}
	update := bson.D{{Key: "$set", Value: user}}
	result, e := db.Mongo.Collection(user.Coll()).UpdateOne(context.TODO(), filter, update)
	if e != nil {
		return nil, err.Mongo(e)
	}
	if e := err.MongoUpdateResult(result); e != nil {
		return nil, e
	}
	return changes, nil
}

func (user *User) UpdatePassword(password string) error {
	bytes, e := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
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
	if e := err.MongoUpdateResult(result); e != nil {
		return e
	}
	return nil
}

func (user *User) ResetPassword() (string, error) {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	result := make([]byte, 8)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		result[i] = letters[num.Int64()]
	}
	password := string(result)

	if e := user.UpdatePassword(password); e != nil {
		return "", e
	}
	return password, nil

}
