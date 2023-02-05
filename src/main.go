package main

import (
	"net/http"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/takeru-a/self-introduction-app-backend/configs"
	"github.com/takeru-a/self-introduction-app-backend/graph"
)

func main() {
	e := echo.New()
	 // ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(
        middleware.CORSConfig{
			AllowCredentials: true,
            // Origin
			AllowOrigins: []string{
				"http://localhost:3000",
			},
        }),
	)
	configs.ConnectDB()
	graphqlHandler  := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	playgroundHandler := playground.Handler("GraphQL", "/query")
	e.POST("/query", func(c echo.Context) error {
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	e.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.GET("/",func(c echo.Context) error {
		return c.JSON(http.StatusOK, "{msg : hello!}")
	})
	e.Logger.Fatal(e.Start(":8080"))
}