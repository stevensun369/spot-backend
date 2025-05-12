package db

import (
	"backend/env"
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx = context.Background()
var RDB *redis.Client
var Client *mongo.Client

// collections
var Accounts *mongo.Collection
var Spots *mongo.Collection

func InitDB() error {
	var err error

	Client, err = mongo.Connect(
		Ctx,
		options.Client().ApplyURI(env.MongoURI),
	)

	if err != nil {
		return err
	}

	// loading collections
	Accounts = GetCollection("accounts", Client)
	Spots = GetCollection("spots", Client)

	fmt.Println("Connected to mongodb")
	return nil
}

func GetCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	return client.Database("dev").Collection(collectionName)
}

func InitCache() error {
	RDB = redis.NewClient(env.RedisOptions)

	pong, _ := RDB.Ping(context.Background()).Result()
	if pong == "PONG" {
		fmt.Println("Connected to redis")
		return nil
	} else {
		fmt.Println("Not connected to redis")
		return errors.New("not connected to redis")
	}
}

func Set(key string, value string) error {
	err := RDB.Set(Ctx, key, value, 0).Err()

	return err
}

func Get(key string) (string, error) {
	val, err := RDB.Get(Ctx, key).Result()

	return val, err
}

func Del(key string) error {
	_, err := RDB.Del(Ctx, key).Result()

	return err
}
