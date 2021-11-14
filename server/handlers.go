package server

import (
	"encoding/hex"
	"github.com/aicam/CryptoNotifier/DB"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body"`
}

func (s *Server) AddUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		newUser := context.Param("username")
		s.DB.Save(&DB.UsersData{
			Username:   newUser,
			LastOnline: time.Now(),
		})
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       "Added",
		})
	}
}

func (s *Server) GetToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		var user DB.UsersData
		username := context.GetHeader("username")
		key := []byte("Ali@Kian")
		if err := s.DB.Where(DB.UsersData{Username: username}).First(&user).Error; err != nil {
			context.JSON(http.StatusUnauthorized, Response{
				StatusCode: -1,
				Body:       "Invalid data",
			})
			return
		}
		user.LastOnline = time.Now()
		s.DB.Save(&user)
		token, err := DesEncrypt([]byte(username), key)
		if err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       hex.EncodeToString(token),
		})
	}
}

func (s *Server) AddInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		var jsData DB.WebData
		var jsExist DB.WebData
		err := context.BindJSON(&jsData)
		if err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		var respBody string
		if s.DB.Where(DB.WebData{UniqueID: jsData.UniqueID}).Find(&jsExist).RecordNotFound() == false {
			if jsExist.UpdatedAt.Before(time.Now().Add(-time.Hour * 10)) {
				go SendNotificationByTelegram(jsData.Body, jsData.Title)
				jsExist.UpdatedAt = time.Now()
				s.DB.Save(&jsExist)
				respBody = "Airdrop reminded again"
			} else {
				respBody = "Airdrop has been added less than 10 hours ago"
			}
		} else {
			go SendNotificationByTelegram(jsData.Body, jsData.Title)
			s.DB.Save(&jsData)
			respBody = "Airdrop added!"
		}
		context.JSON(http.StatusOK, Response{
			StatusCode: 1,
			Body:       respBody,
		})
	}
}

func (s *Server) GetInfo() gin.HandlerFunc {
	return func(context *gin.Context) {
		var DBData []DB.WebData
		offset, err := strconv.Atoi(context.Param("offset"))
		if err != nil {
			context.JSON(http.StatusOK, Response{
				StatusCode: -1,
				Body:       err.Error(),
			})
			return
		}
		s.DB.Find(&DBData)
		context.JSON(http.StatusOK, DBData[len(DBData)-offset:])
	}
}
