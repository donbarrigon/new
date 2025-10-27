package controller

import (
	"donbarrigon/new/internal/app/handler/validator"
	"donbarrigon/new/internal/utils/handler"
)

func UserStore(c *handler.Context) {
	_, e := validator.NewUserStrore(c)
	if e != nil {
		c.ResponseError(e)
		return
	}

}
