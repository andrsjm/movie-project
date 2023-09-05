package model

type Genre struct {
	ID    int    `form:"id" json:"id"`
	Genre string `form:"genre" json:"genre"`
}

type MostViewGenreDTO struct {
	Genre      string `form:"genre" json:"genre"`
	TotalViews string `form:"total_views" json:"total_views"`
}
