package policy

import "donbarrigon/new/internal/utils/handler"

func WorldViewAny(c *handler.Context) error {
	return nil
}

func WorldView(c *handler.Context) error {
	return nil
}

func WorldUpdate(c *handler.Context) error {
	return c.Auth.Can("update world")
}

func WorldDelete(c *handler.Context) error {
	return c.Auth.Can("delete world")
}
