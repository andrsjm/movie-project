package model

type Movie struct {
	ID          int    `form:"id" json:"id"`
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	Duration    int    `form:"duration" json:"duration"`
	Views       int    `form:"views" json:"views"`
	WatchURL    string `form:"watch_url" json:"watch_url"`
	Artists     []int  `form:"artists" json:"artists"` //ArtistID
	Genres      []int  `form:"genres" json:"genres"`   //GenresID
}

type MovieDTO struct {
	ID          int         `form:"id" json:"id"`
	Title       string      `form:"title" json:"title"`
	Description string      `form:"description" json:"description"`
	Duration    int         `form:"duration" json:"duration"`
	Views       int         `form:"views" json:"views"`
	WatchURL    string      `form:"watch_url" json:"watch_url"`
	Artists     []ArtistDTO `form:"artists" json:"artists"`
	Genres      []Genre     `form:"genres" json:"genres"`
}

type MovieFilter struct {
	Offset      int    `form:"Offset" json:"Offset"`
	Limit       int    `form:"Limit" json:"Limit"`
	Title       string `form:"title" json:"title"`
	Description string `form:"Description" json:"Description"`
	Artist      string `form:"artist" json:"artist"`
	Genre       string `form:"genre" json:"genre"`
}

type TrackMovie struct {
	ID    int    `form:"id" json:"id"`
	Title string `form:"title" json:"title"`
	Views int    `form:"views" json:"views"`
}
