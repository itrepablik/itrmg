package itrmg

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Initialize the MongoDB's client's pointer.
var client *mongo.Client

// DP type is a collection of parameters.
type DP map[string]interface{}

// InitMG initializes the MongoDB connections.
func InitMG(dbConStr string) (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbConStr))
	if err != nil {
		return client, err
	}
	ctxMG, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctxMG)
	if err != nil {
		return client, err
	}
	return client, nil
}

// InsertOne insert one row in MongoDB collection.
func InsertOne(dbName, collName string, client *mongo.Client, data DP) (bool, error) {
	collection := client.Database(dbName).Collection(collName)

	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateOne update a single row in MongoDB collection.
func UpdateOne(dbName, collName string, client *mongo.Client, data DP, filter bson.M) (bool, error) {
	collection := client.Database(dbName).Collection(collName)

	update := bson.M{
		"$set": data,
	}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

// UpdateOneByID update a single row filtered by MongoDB object ID from a MongoDB collection.
func UpdateOneByID(dbName, collName string, client *mongo.Client, data DP, objID string) (bool, error) {
	collection := client.Database(dbName).Collection(collName)

	id, err := primitive.ObjectIDFromHex(objID)
	if err != nil {
		return false, err
	}

	filter := bson.M{"_id": bson.M{"$eq": id}}
	update := bson.M{
		"$set": data,
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteOneByID delete a single row filetered MongoDB object ID from a MongoDB collection.
func DeleteOneByID(dbName, collName string, client *mongo.Client, objID string) (bool, error) {
	collection := client.Database(dbName).Collection(collName)

	id, err := primitive.ObjectIDFromHex(objID)
	if err != nil {
		return false, err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return false, err
	}
	return true, nil
}

// FindOneByID find a single row filtered by MongoDB object ID from a collection.
func FindOneByID(dbName, collName string, client *mongo.Client, objID string) (DP, error) {
	collection := client.Database(dbName).Collection(collName)

	var result = make(map[string]interface{})
	id, err := primitive.ObjectIDFromHex(objID)
	if err != nil {
		return result, err
	}

	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Find find a multiple rows filtered by MongoDB object ID from a collection.
func Find(dbName, collName string, client *mongo.Client, filter DP) (DP, error) {
	collection := client.Database(dbName).Collection(collName)

	var results = make(map[string]interface{})
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return results, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		cursor.Decode(&results)
	}

	if err := cursor.Err(); err != nil {
		return results, err
	}
	return results, nil
}
