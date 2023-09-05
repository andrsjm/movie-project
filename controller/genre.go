package controller

import (
	"movie-project/model"
	"movie-project/parser"
	"movie-project/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InsertGenre(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	genre, err := parser.ParseGenre(c)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	_, err = db.Exec("INSERT INTO genres(genre) values(?)",
		genre.Genre)

	var response model.Response
	if err != nil {
		util.ReturnError(c, err)
		return
	} else {
		response.Status = 200
		response.Message = "Success Insert"
		c.Status(http.StatusOK)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

func MostViewedGenre(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	var mostViewedGenre model.MostViewGenreDTO
	var mostViewedGenres []model.MostViewGenreDTO

	query := `SELECT g.genre, SUM(m.views) AS total_views
			FROM movies m JOIN movie_genre d ON m.id = d.movie_id
			JOIN genres g ON d.genre_id = g.id GROUP BY g.genre
			ORDER BY total_views DESC LIMIT 5`

	rows, err := db.Query(query)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&mostViewedGenre.Genre, &mostViewedGenre.TotalViews); err != nil {
			util.ReturnError(c, err)
			return
		} else {

			mostViewedGenres = append(mostViewedGenres, mostViewedGenre)
		}
	}

	var response model.Response
	if len(mostViewedGenres) > 0 {
		response.Status = 200
		response.Message = "Get Success"
		response.Data = mostViewedGenres
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = 400
		response.Message = "No data found"
		c.JSON(http.StatusBadRequest, response)
	}
	c.Header("Content-Type", "application/json")
}
