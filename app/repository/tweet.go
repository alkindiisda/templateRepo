package repository

import (
	"context"

	"a21hc3NpZ25tZW50/app/model"

	"gorm.io/gorm"
)

type TweetRepository interface {
	GetTweetByUserID(ctx context.Context, id int) ([]model.Tweet, error)
	GetAllTweet(ctx context.Context) ([]model.Tweet, error)
	CreateTweet(ctx context.Context, id int, userTweet model.UserTweet) error
	UpdateTweet(ctx context.Context, id int, userTweet model.UserTweet) error
	DeleteTweet(ctx context.Context, id int) error
}

type tweetRepository struct {
	db *gorm.DB
}

func NewTweetRepository(db *gorm.DB) *tweetRepository {
	return &tweetRepository{db}
}

func (t *tweetRepository) GetAllTweet(ctx context.Context) ([]model.Tweet, error) {
	return []model.Tweet{}, nil // TODO: replace this
}

func (t *tweetRepository) GetTweetByUserID(ctx context.Context, id int) ([]model.Tweet, error) {
	return []model.Tweet{}, nil // TODO: replace this
}

func (t *tweetRepository) CreateTweet(ctx context.Context, id int, userTweet model.UserTweet) error {
	return nil // TODO: replace this
}

func (t *tweetRepository) UpdateTweet(ctx context.Context, id int, userTweet model.UserTweet) error {
	return  nil // TODO: replace this
}

func (t *tweetRepository) DeleteTweet(ctx context.Context, id int) error {
	return nil // TODO: replace this
}
