package app

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type App struct {
	gin   *gin.Engine
	db    *gorm.DB
	minio *minio.Client
}

func NewApp(db *gorm.DB, minioClient *minio.Client) *App {
	r := gin.New()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allows all origins
		AllowMethods:     []string{"*"}, // Allows all methods (GET, POST, PUT, DELETE, etc.)
		AllowHeaders:     []string{"*"}, // Allows all headers
		ExposeHeaders:    []string{"*"}, // Exposes all headers to the client
		AllowCredentials: true,
	}))

	return &App{
		gin:  	r,
		db:   	db,
		minio: 	minioClient,
	}
}

func (s *App) Run() error {
	if err := s.MapHandlers(); err != nil {
		return err
	}

	serverURL := fmt.Sprintf(":%s", "8080")
	return s.gin.Run(serverURL)
}
