package api

import (
	"a21hc3NpZ25tZW50/app/service"

	"github.com/gin-gonic/gin"
)

type TweetAPI interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type tweetAPI struct {
	tweetService service.TweetService
}

func NewTweetAPI(tweetService service.TweetService) *tweetAPI {
	return &tweetAPI{tweetService}
}

func (t *tweetAPI) Create(c *gin.Context) {
	// TODO: answer here
}

func (t *tweetAPI) Get(c *gin.Context) {
	// TODO: answer here
}

func (t *tweetAPI) Update(c *gin.Context) {
	// TODO: answer here
}

func (t *tweetAPI) Delete(c *gin.Context) {
	// TODO: answer here
}
