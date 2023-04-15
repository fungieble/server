package router

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

func NewRouter(e *echo.Echo, srv *handler.Server) *echo.Echo {
	port := viper.GetString("PORT")
	host := viper.GetString("HOST")
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "{ok}")
	})

	e.POST("/query", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/subscriptions", func(c echo.Context) error {
		srv.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/graphql", func(c echo.Context) error {
		playground.Handler("GraphQL", "/query").ServeHTTP(c.Response(), c.Request())
		return nil
	})

	log.Printf("GraphQL playground at http://%s:%s/\n\n", host, port)
	return e
}
