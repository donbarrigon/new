package policy

import (
	"donbarrigon/new/internal/app/data/model"
	"donbarrigon/new/internal/utils/handler"
)

func UserViewAny(c *handler.Context) error {
	return c.Auth.Can("view-any user")
}

func UserView(c *handler.Context, user *model.User) error {
	if user.ID == c.Auth.UserID() {
		return nil
	}
	return c.Auth.Can("view user")
}

func UserCreate(c *handler.Context) error {
	return nil
}

func UserUpdate(c *handler.Context, user *model.User) error {
	if user.ID == c.Auth.UserID() {
		return nil
	}
	return c.Auth.Can("update user")
}

func UserDelete(c *handler.Context, user *model.User) error {
	if user.ID == c.Auth.UserID() {
		return nil
	}
	return c.Auth.Can("delete user")
}
