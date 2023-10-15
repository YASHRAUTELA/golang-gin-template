package models

import (
	"github.com/gin-gonic/gin"
)

type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

// postAlbums adds an album from JSON received in the request body.
func CreateAlbums(ctx *gin.Context) (Album, error) {
	var newAlbum Album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := ctx.BindJSON(&newAlbum); err != nil {
		return newAlbum, err
	}
	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	return newAlbum, nil
}

// getAlbums responds with the list of all albums as JSON.
func FetchAlbums() []Album {
	return albums
}
