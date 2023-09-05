package parser

import (
	"encoding/json"
	"io/ioutil"
	"movie-project/model"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseMovie(c *gin.Context) (movie model.Movie, err error) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)

	err = json.Unmarshal(jsonData, &movie)

	return movie, nil
}

func ParseMovieFilter(c *gin.Context) model.MovieFilter {
	filter := model.MovieFilter{}

	filter.Limit, _ = strconv.Atoi(c.Query("limit"))
	filter.Offset, _ = strconv.Atoi(c.Query("offset"))
	filter.Title = c.Query("title")
	filter.Description = c.Query("description")
	filter.Artist = c.Query("artist")
	filter.Genre = c.Query("genre")

	return filter
}
