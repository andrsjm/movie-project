package main

import (
	"fmt"
	"movie-project/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.POST("/movie", controller.InserMovie)
	router.POST("/artist", controller.InsertArtist)
	router.POST("/genre", controller.InsertGenre)
	router.POST("/user/register", controller.Register)
	router.POST("/user/login", controller.Login)
	router.POST("/user/logout", controller.Logout)
	router.POST("/movie/vote/:id", controller.InsertVote)
	router.POST("/movie/unvote/:id", controller.UnvoteMovie)

	router.GET("/movie/most-viewed", controller.MostViewedMovie)
	router.GET("/genre/most-viewed", controller.MostViewedGenre)
	router.GET("/movie", controller.GetAllMovie)
	router.GET("/movie/search", controller.Search)
	router.GET("/watch/:id", controller.WatchMovie)
	router.GET("/movie/track/:id", controller.TrackMovie)
	router.GET("/user/movie/voted", controller.GetVotedMovie)

	router.PUT("/movie", controller.UpdateMovie)

	router.Run(":9090")
	fmt.Println("Run on port 9090")
}
