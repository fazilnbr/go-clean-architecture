package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fazilnbr/go-clean-architecture/internal/app/users"
	clients "github.com/fazilnbr/go-clean-architecture/internal/client"
	"github.com/fazilnbr/go-clean-architecture/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	CONFIG_FILE_NAME = "config"
	CONFIG_FILE_EXT  = ".json"
	CONFIG_FILE_PATH = "configs"
)

type serverResource struct {
	postgresDB     *gorm.DB
	tokenSecretKey string
}

type Server struct {
	router    *gin.Engine
	resources serverResource
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	appConfig, err := s.ReadConfig()
	if err != nil {
		return errors.Wrap(err, "error reading config")
	}

	baseRouter := gin.Default()
	s.router = baseRouter

	err = s.InitResources(appConfig)
	if err != nil {
		return errors.Wrap(err, "error initializing resource")
	}
	err = s.InitServices()
	if err != nil {
		return errors.Wrap(err, "error initializing services")
	}

	srv := http.Server{
		Addr:    ":9090",
		Handler: s.router,
	}
	log.Printf("starting server on port %s", srv.Addr)
	err = srv.ListenAndServe()
	return err
}

func (s *Server) ReadConfig() (*utils.AppConfig, error) {
	filename := CONFIG_FILE_NAME + CONFIG_FILE_EXT
	configPath := CONFIG_FILE_PATH
	configFile := filepath.Join(configPath, filename)
	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var appConfig *utils.AppConfig
	err = json.Unmarshal(configBytes, &appConfig)
	return appConfig, err
}

func (s *Server) InitResources(appConfig *utils.AppConfig) error {
	usersDB, err := clients.NewPostgresDB(appConfig.Postgres)
	if err != nil {
		return err
	}
	resources := serverResource{
		postgresDB:     usersDB,
		tokenSecretKey: appConfig.TokenSecretKey,
	}
	s.resources = resources
	return nil
}

func (s *Server) InitServices() error {
	usersDao := users.NewDAO(s.resources.postgresDB)
	middleware := users.NewMiddleware(s.resources.tokenSecretKey)
	usersSvc := users.NewService(usersDao)
	usersHandler := users.NewHttpHandler(usersSvc)

	apiGroup := s.router.Group("/api")
	{
		userGroup := apiGroup.Group("/users")
		{
			usersHandler.InitRoutes(userGroup, middleware)
		}
	}

	return nil
}
