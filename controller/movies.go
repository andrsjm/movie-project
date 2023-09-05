package controller

import (
	"database/sql"
	"fmt"
	"movie-project/model"
	"movie-project/parser"
	"movie-project/util"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func InserMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	tx, err := db.Begin()
	defer tx.Rollback()

	movie, err := parser.ParseMovie(c)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	res, err := tx.Exec("INSERT INTO movies(title, description, duration, watch_url) values(?,?,?,?)",
		movie.Title,
		movie.Description,
		movie.Duration,
		movie.WatchURL)

	if err != nil {
		util.ReturnError(c, err)
		return
	}

	lid, err := res.LastInsertId()

	err = insertMovieArtist(movie.Artists, lid, tx)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	err = insertMovieGenres(movie.Genres, lid, tx)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	var response model.Response
	if err := tx.Commit(); err != nil {
		util.ReturnError(c, err)
	} else {
		response.Status = 200
		response.Message = "Success Insert"
		c.Status(http.StatusOK)
	}
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, response)

}

func insertMovieArtist(artists []int, lid int64, tx *sql.Tx) error {
	for _, artistID := range artists {
		_, err := tx.Exec("INSERT INTO movie_artist(artist_id, movie_id) values(?,?)",
			artistID,
			lid)
		if err != nil {
			return err
		}
	}

	return nil
}

func insertMovieGenres(genres []int, lid int64, tx *sql.Tx) error {
	for _, genreID := range genres {
		_, err := tx.Exec("INSERT INTO movie_genre(genre_id, movie_id) values(?,?)",
			genreID,
			lid)
		if err != nil {
			return err
		}
	}

	return nil
}

func MostViewedMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	var movie model.MovieDTO
	var movies []model.MovieDTO

	query := "SELECT * FROM movies ORDER BY views DESC LIMIT 5"

	rows, err := db.Query(query)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.WatchURL, &movie.Views); err != nil {
			util.ReturnError(c, err)
			return
		} else {
			movie.Artists, err = getArtistsMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movie.Genres, err = getGenresMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movies = append(movies, movie)
		}
	}

	var response model.Response
	if len(movies) > 0 {
		response.Status = 200
		response.Message = "Get Success"
		response.Data = movies
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = 400
		response.Message = "No data found"
		c.JSON(http.StatusBadRequest, response)
	}
	c.Header("Content-Type", "application/json")
}

func getArtistsMovie(movie_id int) ([]model.ArtistDTO, error) {
	db := util.Connect()
	defer db.Close()

	var artist model.ArtistDTO
	var artists []model.ArtistDTO

	query := "SELECT a.* from artists a join movie_artist b on a.id = b.artist_id WHERE b.movie_id = " + strconv.Itoa(movie_id)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.DOB); err != nil {
			return nil, err
		} else {
			artists = append(artists, artist)
		}
	}

	return artists, nil
}

func getGenresMovie(movie_id int) ([]model.Genre, error) {
	db := util.Connect()
	defer db.Close()

	var genre model.Genre
	var genres []model.Genre

	query := "SELECT g.* from genres g join movie_genre b on g.id = b.genre_id WHERE b.movie_id = " + strconv.Itoa(movie_id)

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&genre.ID, &genre.Genre); err != nil {
			return nil, err
		} else {
			genres = append(genres, genre)
		}
	}

	return genres, nil
}

func GetAllMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	var movie model.MovieDTO
	var movies []model.MovieDTO

	/*
		Disini saya memakai OFFSET dan LIMIT sebagai pengatur pagination. Dalam case ini, front end yang menentukan seberapa banyak limit item per page.
	*/
	filter := parser.ParseMovieFilter(c)

	query := "SELECT * FROM movies LIMIT ? OFFSET ?"

	rows, err := db.Query(query, filter.Limit, filter.Offset)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.WatchURL, &movie.Views); err != nil {
			util.ReturnError(c, err)
			return
		} else {
			movie.Artists, err = getArtistsMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movie.Genres, err = getGenresMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movies = append(movies, movie)
		}
	}

	var response model.Response
	if len(movies) > 0 {
		response.Status = 200
		response.Message = "Get Success"
		response.Data = movies
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = 400
		response.Message = "No data found"
		c.JSON(http.StatusBadRequest, response)
	}
	c.Header("Content-Type", "application/json")
}

func Search(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	var movie model.MovieDTO
	var movies []model.MovieDTO

	filter := parser.ParseMovieFilter(c)

	query := `SELECT DISTINCT m.* FROM movies m LEFT JOIN movie_artist ma ON m.id = ma.movie_id
		LEFT JOIN artists a ON ma.artist_id = a.id LEFT JOIN movie_genre mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id`

	if filter.Title != "" || filter.Description != "" || filter.Artist != "" || filter.Genre != "" {
		query += " WHERE"
	}

	if filter.Title != "" {
		query += " m.title LIKE '%" + filter.Title + "%' OR"
	}

	if filter.Description != "" {
		query += " m.description LIKE '%" + filter.Description + "%' OR"
	}

	if filter.Artist != "" {
		query += " a.name LIKE '%" + filter.Artist + "%' OR"
	}

	if filter.Genre != "" {
		query += " g.genre LIKE '%" + filter.Genre + "%'"
	}

	query = strings.TrimRight(query, "OR")

	fmt.Println()

	rows, err := db.Query(query)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.WatchURL, &movie.Views); err != nil {
			util.ReturnError(c, err)
			return
		} else {
			movie.Artists, err = getArtistsMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movie.Genres, err = getGenresMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movies = append(movies, movie)
		}
	}

	var response model.Response
	if len(movies) > 0 {
		response.Status = 200
		response.Message = "Get Success"
		response.Data = movies
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = 400
		response.Message = "No data found"
		c.JSON(http.StatusBadRequest, response)
	}
	c.Header("Content-Type", "application/json")
}

func UpdateMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	movie, err := parser.ParseMovie(c)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	_, err = db.Exec("UPDATE movies SET title=?, description=?, duration=?, watch_url=?, views=? WHERE id=?",
		movie.Title,
		movie.Description,
		movie.Duration,
		movie.WatchURL,
		movie.Views,
		movie.ID)

	var response model.Response
	if err != nil {
		util.ReturnError(c, err)
		return
	} else {
		response.Message = "Success Update"
		response.Status = 200
		c.JSON(http.StatusOK, response)
	}
	c.Header("Content-Type", "application/json")
}

func TrackMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	var trackMovie model.TrackMovie

	movieID := c.Param("id")

	err := db.QueryRow("SELECT id, title, views FROM movies WHERE id=?", movieID).Scan(&trackMovie.ID, &trackMovie.Title, &trackMovie.Views)

	var response model.Response
	if err != nil {
		util.ReturnError(c, err)
		return
	} else {
		response.Status = 200
		response.Message = "Success Get Movies Viewership"
		response.Data = trackMovie
		c.JSON(http.StatusOK, response)
	}
	c.Header("Content-Type", "application/json")
}

func WatchMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	session := sessions.Default(c)

	userID := session.Get("userID")
	movieID := c.Param("id")

	if userID == nil {
		util.UnauthorizedError(c)
		return
	}

	_, err := db.Exec("INSERT INTO watch_history(movie_id, user_id) values(?, ?)",
		movieID,
		userID)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	_, err = db.Exec("UPDATE movies SET views=views+1 WHERE id=?", movieID)

	var response model.Response
	if err != nil {
		util.ReturnError(c, err)
		return
	} else {
		response.Message = "Success Watch"
		response.Status = 200
		c.JSON(http.StatusOK, response)
	}
	c.Header("Content-Type", "application/json")
}

func WatchHistory(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	var movie model.MovieDTO
	var movies []model.MovieDTO

	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		util.UnauthorizedError(c)
		return
	}

	/*
		Disini saya mengasumsikan bahwa track movie itu adalah user melihat history film yang telah ditontonnya
	*/

	query := "SELECT m.* FROM movies m INNER JOIN watch_history wh ON m.id = wh.movie_id INNER JOIN users s ON wh.user_id = s.id WHERE s.id = ?"

	rows, err := db.Query(query, userID)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.WatchURL, &movie.Views); err != nil {
			util.ReturnError(c, err)
			return
		} else {
			movie.Artists, err = getArtistsMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movie.Genres, err = getGenresMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movies = append(movies, movie)
		}
	}

	var response model.Response
	if len(movies) > 0 {
		response.Status = 200
		response.Message = "Get Success"
		response.Data = movies
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = 400
		response.Message = "No data found"
		c.JSON(http.StatusBadRequest, response)
	}
	c.Header("Content-Type", "application/json")
}

func InsertVote(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	session := sessions.Default(c)
	userID := session.Get("userID")
	movieID := c.Param("id")

	if userID == nil {
		util.UnauthorizedError(c)
		return
	}

	var response model.Response
	isExist, err := checkUserVoted(userID, movieID)
	if err != nil {
		util.ReturnError(c, err)
		return
	}
	if isExist == 0 {
		_, err = db.Exec("INSERT INTO voted_movie(movie_id, user_id) values(?,?)",
			movieID,
			userID)

		if err != nil {
			util.ReturnError(c, err)
			return
		} else {
			response.Message = "Success Votes"
			response.Status = 200
			c.JSON(http.StatusOK, response)
			c.Header("Content-Type", "application/json")
			return
		}
	} else {
		response.Message = "You have voted this film"
		response.Status = 400
		c.JSON(http.StatusBadRequest, response)
		c.Header("Content-Type", "application/json")
		return
	}

}

func UnvoteMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	session := sessions.Default(c)
	userID := session.Get("userID")
	movieID := c.Param("id")

	if userID == nil {
		util.UnauthorizedError(c)
		return
	}

	_, err := db.Exec("DELETE FROM voted_movie WHERE movie_id=? AND user_id=?", movieID, userID)

	var response model.Response
	if err != nil {
		util.ReturnError(c, err)
		return
	} else {
		response.Status = 200
		response.Message = "Success Unvote"
		c.JSON(http.StatusOK, response)
	}
	c.Header("Content-Type", "application/json")

}

func checkUserVoted(userID interface{}, movieID string) (int, error) {
	db := util.Connect()
	defer db.Close()

	var isExist int

	err := db.QueryRow("SELECT EXISTS(SELECT * FROM voted_movie WHERE user_id = ? AND movie_id = ?)", userID, movieID).Scan(&isExist)

	return isExist, err
}

func GetVotedMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	var movie model.MovieDTO
	var movies []model.MovieDTO

	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		util.UnauthorizedError(c)
		return
	}

	query := "SELECT m.* FROM movies m INNER JOIN voted_movie vm ON m.id = vm.movie_id INNER JOIN users u ON vm.user_id = u.id WHERE u.id = ?"

	rows, err := db.Query(query, userID)
	if err != nil {
		util.ReturnError(c, err)
		return
	}

	for rows.Next() {
		if err := rows.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Duration, &movie.WatchURL, &movie.Views); err != nil {
			util.ReturnError(c, err)
			return
		} else {
			movie.Artists, err = getArtistsMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movie.Genres, err = getGenresMovie(movie.ID)
			if err != nil {
				util.ReturnError(c, err)
				return
			}

			movies = append(movies, movie)
		}
	}

	var response model.Response
	if len(movies) > 0 {
		response.Status = 200
		response.Message = "Get Success"
		response.Data = movies
		c.JSON(http.StatusOK, response)
	} else {
		response.Status = 400
		response.Message = "No data found"
		c.JSON(http.StatusBadRequest, response)
	}
	c.Header("Content-Type", "application/json")
}

func MostVotedMovie(c *gin.Context) {
	db := util.Connect()
	defer db.Close()

	session := sessions.Default(c)
	userID := session.Get("userID")

	if userID == nil {
		util.UnauthorizedError(c)
		return
	}

	var mostVotedMovie model.MostVoteMovie

	err := db.QueryRow(`SELECT m.title AS movie_title, COUNT(vm.movie_id) AS vote_count
			FROM movies m
			LEFT JOIN voted_movie vm ON m.id = vm.movie_id
			GROUP BY m.title
			ORDER BY vote_count DESC
			LIMIT 1`).Scan(&mostVotedMovie.Title, &mostVotedMovie.VoteCount)

	var response model.Response
	if err != nil {
		util.ReturnError(c, err)
		return
	} else {
		response.Status = 200
		response.Message = "Success Get Most Voted Movie"
		response.Data = mostVotedMovie
		c.JSON(http.StatusOK, response)
	}
	c.Header("Content-Type", "application/json")
}
