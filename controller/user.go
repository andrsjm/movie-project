package controller

import (
	"movie-project/model"
	"movie-project/parser"
	"movie-project/util"
	"net/http"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	user, err := parser.ParseUser(c)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	_, err = db.Exec("INSERT INTO users(name, email, password) values(?, ?, ?)",
		user.Name,
		user.Email,
		user.Password)

	var response model.Response
	if err != nil {
		util.ReturnError(c, err)
		return
	} else {
		response.Status = 200
		response.Message = "Success Register"
		c.Status(http.StatusOK)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)
}

func Login(c *gin.Context) {
	db := util.Connect()
	defer db.Close()
	session := sessions.Default(c)

	user, err := parser.ParseUser(c)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	var userLogin model.User

	var response model.Response
	err = db.QueryRow("SELECT * FROM users WHERE email=? AND password=?", user.Email, user.Password).Scan(&userLogin.ID, &userLogin.Name, &userLogin.Email, &userLogin.Password)
	if err != nil {
		response.Status = 400
		response.Message = "Email or Password wrong"
		c.JSON(http.StatusBadRequest, response)
		c.Header("Content-Type", "application/json")
		return
	} else {
		response.Status = 200
		response.Message = "Login Success"
		session.Set("userID", userLogin.ID)
		session.Save()
		c.JSON(http.StatusOK, response)
	}

	c.Header("Content-Type", "application/json")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)

	userID := session.Get("userID")

	var response model.Response
	if userID != nil {
		session.Clear()
		session.Save()
		response.Status = 200
		response.Message = "Logout Success"
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = 400
		response.Message = "you have never logged in"
		c.JSON(http.StatusBadRequest, response)
	}
	c.Header("Content-Type", "application/json")
}
