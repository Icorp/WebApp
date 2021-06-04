package pkg

import "go.mongodb.org/mongo-driver/mongo"

func newRepository(db *mongo.Client) *mongo.Client {
	return &mongo.Client{}
}
