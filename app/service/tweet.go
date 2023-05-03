package service

import (
	"context"

	"a21hc3NpZ25tZW50/app/model"
	"a21hc3NpZ25tZW50/app/repository"
)

type TweetService interface {
	CreateTweet(ctx context.Context, id int, userTweet model.UserTweet) error
	GetAllTweet(ctx context.Context) ([]model.TweetContent, error)
	UserName(ctx context.Context, id int) string
	UserFullname(ctx context.Context, id int) string

	UpdateTweet(ctx context.Context, id int, tweet model.UserTweet) error
	DeleteTweet(ctx context.Context, id int) error
}

type tweetService struct {
	tweetRepo repository.TweetRepository
	userRepo  repository.UserRepository
}

func NewTweetService(tweetRepo repository.TweetRepository, userRepo repository.UserRepository) TweetService {
	return &tweetService{tweetRepo, userRepo}
}

func (t *tweetService) CreateTweet(ctx context.Context, id int, userTweet model.UserTweet) error {
	err := t.tweetRepo.CreateTweet(ctx, id, userTweet)
	if err != nil {
		return err
	}

	return nil
}

func (t *tweetService) GetAllTweet(ctx context.Context) ([]model.TweetContent, error) {
	resTweets, err := t.tweetRepo.GetAllTweet(ctx)
	if err != nil {
		return nil, err
	}

	var TweetContents []model.TweetContent
	for _, v := range resTweets {
		TweetContents = append(TweetContents, model.TweetContent{
			ID:        v.ID,
			Fullname:  t.UserFullname(ctx, v.UserID),
			Username:  t.UserName(ctx, v.UserID),
			Tweet:     v.Tweet,
			Image:     v.Image,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return TweetContents, nil
}

func (t *tweetService) UserName(ctx context.Context, id int) string {
	resUser, _ := t.userRepo.GetUserByID(ctx, id)
	return resUser.Username
}

func (t *tweetService) UserFullname(ctx context.Context, id int) string {
	resUser, _ := t.userRepo.GetUserByID(ctx, id)
	return resUser.Fullname
}

func (t *tweetService) UpdateTweet(ctx context.Context, id int, tweet model.UserTweet) error {
	err := t.tweetRepo.UpdateTweet(ctx, id, tweet)
	if err != nil {
		return err
	}

	return nil
}

func (t *tweetService) DeleteTweet(ctx context.Context, id int) error {
	err := t.tweetRepo.DeleteTweet(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
