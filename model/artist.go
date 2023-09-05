package model

import "time"

type ArtistDTO struct {
	ID   int       `form:"id" json:"id"`
	Name string    `form:"name" json:"name"`
	DOB  time.Time `form:"dob" json:"dob"`
}

type Artist struct {
	ID   int    `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
	DOB  string `form:"dob" json:"dob"`
}
