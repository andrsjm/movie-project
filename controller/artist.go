package controller

import (
	"movie-project/model"
	"movie-project/parser"
	"movie-project/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InsertArtist(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	artist, err := parser.ParseArtist(c)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	_, err = db.Exec("INSERT INTO artists(name, dob) values(?, ?)",
		artist.Name,
		artist.DOB)

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
