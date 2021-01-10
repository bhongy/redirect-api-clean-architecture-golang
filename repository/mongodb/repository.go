package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/bhongy/rediret-api-clean-architecture-golang/shortener"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newMongoClient(mongoURL string, mongoTimeoutSeconds int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeoutSeconds)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepository(mongoURL, mongoDB string, mongoTimeoutSeconds int) (shortener.RedirectRepository, error) {
	repo := mongoRepository{
		database: mongoDB,
		timeout:  time.Duration(mongoTimeoutSeconds) * time.Second,
	}
	client, err := newMongoClient(mongoURL, mongoTimeoutSeconds)
	if err != nil {
		return nil, fmt.Errorf("repository.NewMongoRepo: %v", err)
	}
	repo.client = client
	return &repo, nil
}

func (r *mongoRepository) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var redirect shortener.Redirect
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("repository.Redirect.Find: %w", shortener.ErrRedirectNotFound)
		}
		return nil, fmt.Errorf("repository.Redirect.Find: %v", err)
	}
	return &redirect, nil
}

func (r *mongoRepository) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return fmt.Errorf("repository.Redirect.Store: %v", err)
	}
	return nil
}
