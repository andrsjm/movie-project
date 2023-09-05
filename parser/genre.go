package parser

import (
	"encoding/json"
	"io/ioutil"
	"movie-project/model"

	"github.com/gin-gonic/gin"
)

func ParseGenre(c *gin.Context) (genre model.Genre, err error) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)

	err = json.Unmarshal(jsonData, &genre)

	return genre, nil
}
