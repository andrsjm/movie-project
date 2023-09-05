package parser

import (
	"encoding/json"
	"io/ioutil"
	"movie-project/model"

	"github.com/gin-gonic/gin"
)

func ParseArtist(c *gin.Context) (artist model.Artist, err error) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)

	err = json.Unmarshal(jsonData, &artist)

	return artist, nil
}
