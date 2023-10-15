package controllers

import (
	"myapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAlbums(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": models.FetchAlbums(), "message": ""})
}

func PostAlbums(ctx *gin.Context) {
	var response, err = models.CreateAlbums(ctx)
	if err != nil {
		ctx.JSON(http.StatusLengthRequired, gin.H{"status": http.StatusLengthRequired, "data": response, "message": err})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": response, "message": ""})
}
