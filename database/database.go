package database

import (
	"context"
	"log"
	"time"

	"github.com/akhil/gql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString string = "mongodb://localhost:27017/"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		client: client,
	}
}

func (db *DB) GetProfile(id string) *model.CustomerProfile {
	jobCollec := db.client.Database("graphql-customer-profile").Collection("customers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var customerProfile model.CustomerProfile
	err := jobCollec.FindOne(ctx, filter).Decode(&customerProfile)
	if err != nil {
		log.Fatal(err)
	}
	return &customerProfile
}

func (db *DB) GetProfiles() []*model.CustomerProfile {
	jobCollec := db.client.Database("graphql-customer-profile").Collection("customers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var customerProfiles []*model.CustomerProfile
	cursor, err := jobCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &customerProfiles); err != nil {
		panic(err)
	}

	return customerProfiles
}

func (db *DB) CreateCustomerProfile(profileInfo model.CreateCustomerProfileInput) *model.CustomerProfile {
	jobCollec := db.client.Database("graphql-customer-profile").Collection("customers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := jobCollec.InsertOne(ctx, bson.M{"title": profileInfo.Title, "description": profileInfo.Description, "email": profileInfo.Email, "type": profileInfo.Type})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnJobListing := model.CustomerProfile{ID: insertedID, Title: profileInfo.Title, Type: profileInfo.Type, Description: profileInfo.Description, Email: profileInfo.Email}
	return &returnJobListing
}

func (db *DB) UpdateCustomerProfile(jobId string, jobInfo model.UpdateCustomerProfileInput) *model.CustomerProfile {
	jobCollec := db.client.Database("graphql-customer-profile").Collection("customers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateJobInfo := bson.M{}

	if jobInfo.Title != nil {
		updateJobInfo["title"] = jobInfo.Title
	}
	if jobInfo.Description != nil {
		updateJobInfo["description"] = jobInfo.Description
	}
	if jobInfo.Eamil != nil {
		updateJobInfo["email"] = jobInfo.Eamil
	}

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateJobInfo}

	results := jobCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var customerProfile model.CustomerProfile

	if err := results.Decode(&customerProfile); err != nil {
		log.Fatal(err)
	}

	return &customerProfile
}

func (db *DB) DeleteCustomerProfile(customerId string) *model.DeleteCustomerResponse {
	jobCollec := db.client.Database("graphql-customer-profile").Collection("customers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(customerId)
	filter := bson.M{"_id": _id}
	_, err := jobCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return &model.DeleteCustomerResponse{DeletedCusID: customerId}
}
