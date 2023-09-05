package parser

import (
	"encoding/json"
	"io/ioutil"
	"movie-project/model"
	"movie-project/util"

	"github.com/gin-gonic/gin"
)

func ParseUser(c *gin.Context) (user model.User, err error) {
	jsonData, err := ioutil.ReadAll(c.Request.Body)

	err = json.Unmarshal(jsonData, &user)

	user.Password = util.HashPassword(user.Password)

	return user, nil
}
