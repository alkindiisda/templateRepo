package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"a21hc3NpZ25tZW50/app/api"
	"a21hc3NpZ25tZW50/app/middleware"
	"a21hc3NpZ25tZW50/app/repository"
	"a21hc3NpZ25tZW50/app/service"
	"a21hc3NpZ25tZW50/app/utils"

	docs "a21hc3NpZ25tZW50/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type APIHandler struct {
	UserAPIHandler  api.UserAPI
	TweetAPIHandler api.TweetAPI
}

func main() {
	if os.Getenv("DATABASE_URL") == "" {
		os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/kampusmerdeka")
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		router := gin.New()

		router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[%s] \"%s %s %s\"\n",
				param.TimeStamp.Format(time.RFC822),
				param.Method,
				param.Path,
				param.ErrorMessage,
			)
		}))
		router.Use(gin.Recovery())

		err := utils.ConnectDB()
		if err != nil {
			panic(err)
		}

		db := utils.GetDBConnection()

		router = RunServer(db, router)

		fmt.Println("Server is running on port 8080")
		err = router.Run()
		if err != nil {
			panic(err)
		}
	}()

	wg.Wait()
}

// @BasePath /api/v1

func RunServer(db *gorm.DB, gin *gin.Engine) *gin.Engine {
	docs.SwaggerInfo.BasePath = "/api/v1"
	userRepo := repository.NewUserRepository(db)
	tweetRepo := repository.NewTweetRepository(db)

	userService := service.NewUserService(userRepo)
	tweetService := service.NewTweetService(tweetRepo, userRepo)

	userAPIHandler := api.NewUserAPI(userService)
	tweerAPIHandler := api.NewTweetAPI(tweetService)

	apiHandler := APIHandler{
		UserAPIHandler:  userAPIHandler,
		TweetAPIHandler: tweerAPIHandler,
	}

	gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	gin.Use(middleware.CORSMiddleware())
	users := gin.Group("/api/v1/users")
	{
		users.POST("/login", apiHandler.UserAPIHandler.Login)
		users.POST("/register", apiHandler.UserAPIHandler.Register)

		users.Use(middleware.Auth())
		users.GET("/get", apiHandler.UserAPIHandler.GetByID)
		users.PUT("/update", apiHandler.UserAPIHandler.UpdateByID)
		users.DELETE("/delete", apiHandler.UserAPIHandler.Delete)
		users.GET("/logout", apiHandler.UserAPIHandler.Logout)
	}

	tweets := gin.Group("/api/v1/tweets")
	{
		tweets.Use(middleware.Auth())
		tweets.POST("/create", apiHandler.TweetAPIHandler.Create)
		tweets.GET("/get", apiHandler.TweetAPIHandler.Get)
		tweets.PUT("/update", apiHandler.TweetAPIHandler.Update)
		tweets.DELETE("/delete", apiHandler.TweetAPIHandler.Delete)
	}

	return gin
}
