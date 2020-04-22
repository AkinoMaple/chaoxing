package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Model struct {
	Users *mongo.Collection
}

func (m *Model) Init() {
	ctx, cancal := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancal()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDB.ApplyURI))
	if err != nil {
		log.Panic(err)
	}

	ctx, cancal = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancal()
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Panic(err)
	}

	m.Users = client.Database(config.MongoDB.Database).Collection(config.MongoDB.Collection)
}

func (m *Model) UpsertOneUser(data Data) (res *mongo.UpdateResult, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Update().SetUpsert(true)
	res, err = m.Users.UpdateOne(ctx, data, bson.M{"$set": bson.M{"puid": data.Puid}}, opts)

	return
}
