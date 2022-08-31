package hecho

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type (
	QueryResource interface {
		Query(c echo.Context) error
	}

	GetResource interface {
		Get(c echo.Context) error
	}

	CreateResource interface {
		Create(c echo.Context) error
	}

	UpdateResource interface {
		Update(c echo.Context) error
	}

	PatchResource interface {
		Patch(c echo.Context) error
	}

	DeleteResource interface {
		Delete(c echo.Context) error
	}
)

// ResourceAPI defines every http route which its method is satisfied by the resource.
func ResourceAPI(group *echo.Group, resource any, prefix string, m ...echo.MiddlewareFunc) {
	if r, ok := resource.(QueryResource); ok {
		group.GET("", r.Query, m...).Name = routeName(prefix, "query")
	}

	if r, ok := resource.(GetResource); ok {
		group.GET("/:id", r.Get, m...).Name = routeName(prefix, "get")
	}

	if r, ok := resource.(CreateResource); ok {
		group.POST("", r.Create, m...).Name = routeName(prefix, "create")
	}

	if r, ok := resource.(UpdateResource); ok {
		group.PUT("/:id", r.Update, m...).Name = routeName(prefix, "put")
	}

	if r, ok := resource.(PatchResource); ok {
		group.PATCH("/:id", r.Patch, m...).Name = routeName(prefix, "patch")
	}

	if r, ok := resource.(DeleteResource); ok {
		group.DELETE("/:id", r.Delete, m...).Name = routeName(prefix, "delete")
	}
}

func routeName(prefix, name string) string {
	return fmt.Sprintf("%s::%s", prefix, name)
}
