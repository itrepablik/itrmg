package itrmg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ClientMG initialize the MongoDB's client's pointer.
var ClientMG *mongo.Client

// MGC is a Mongo client type
type MGC *mongo.Client

// ObjID is the MongoDB BSON ObjectID type.
type ObjID primitive.ObjectID

// DP type is a data parameters to be used as common map container for collection results
// or use as the filter parameters, etc.
type DP map[string]interface{}

// DM type is a slice map container for collection results storage.
type DM []map[string]interface{}

// InitMG initializes the MongoDB connections.
func InitMG(dbConStr string) (*mongo.Client, error) {
	ClientMG, err := mongo.NewClient(options.Client().ApplyURI(dbConStr))
	if err != nil {
		return ClientMG, err
	}
	ctxMG, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = ClientMG.Connect(ctxMG)
	if err != nil {
		return ClientMG, err
	}
	return ClientMG, nil
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
func UpdateOne(dbName, collName string, client *mongo.Client, data DP, filter DP) (bool, error) {
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

// DeleteOneByID delete any single row permanently filetered MongoDB object ID from a MongoDB collection.
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

// DeleteOne delete any single row permanently from a collection.
func DeleteOne(dbName, collName string, client *mongo.Client, filter DP) (bool, error) {
	collection := client.Database(dbName).Collection(collName)
	_, err := collection.DeleteOne(context.TODO(), filter)
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

// FindOne find a single row and retrieves all columns from a collection.
func FindOne(dbName, collName string, client *mongo.Client, filter DP) (DP, error) {
	collection := client.Database(dbName).Collection(collName)
	var result = make(map[string]interface{})

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Find find a multiple rows filtered by MongoDB object ID from a collection.
func Find(dbName, collName string, client *mongo.Client, filter DP, sortOrder DP, setLimit int64) (DM, error) {
	opts := options.Find()
	opts.SetSort(sortOrder)
	if setLimit > 0 {
		opts.SetLimit(setLimit)
	}

	collection := client.Database(dbName).Collection(collName)
	results := []map[string]interface{}{}
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return results, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		rowData := make(map[string]interface{})
		cursor.Decode(&rowData)
		results = append(results, rowData)
	}
	if err := cursor.Err(); err != nil {
		return results, err
	}
	return results, nil
}

// IsExist find any single row from a specified collection.
func IsExist(dbName, collName string, client *mongo.Client, filter DP) (bool, error) {
	collection := client.Database(dbName).Collection(collName)
	var result = make(map[string]interface{})
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetFieldValue gets the string value of any specific field name from a collection.
// e.g Struct field "PCName" then the bson field is "pc_name", the bson field name must be
func GetFieldValue(dbName, collName string, client *mongo.Client, filter DP, bsonFieldName string) (string, error) {
	opts := options.Find()
	opts.SetLimit(1)
	opts.SetProjection(bson.M{
		bsonFieldName: 1,
		"score":       bson.M{"$meta": "textScore"},
	})
	opts.SetSort(bson.M{"score": bson.M{"$meta": "textScore"}})

	collection := client.Database(dbName).Collection(collName)
	results := []map[string]interface{}{}
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return "", err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		rowData := make(map[string]interface{})
		cursor.Decode(&rowData)
		results = append(results, rowData)
	}
	if err := cursor.Err(); err != nil {
		return "", err
	}

	strVal := ""
	for _, value := range results {
		strVal = fmt.Sprintf("%v", value[bsonFieldName])
	}
	return strVal, nil
}

// CountRows gets the total number of rows from a collection.
func CountRows(dbName, collName string, client *mongo.Client, filter DP) (int64, error) {
	collection := client.Database(dbName).Collection(collName)
	opts := options.Count().SetMaxTime(2 * time.Second)
	count, err := collection.CountDocuments(context.TODO(), filter, opts)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetFieldValueByID gets the string value of any specific field filtered by object id from a collection.
func GetFieldValueByID(dbName, collName string, client *mongo.Client, objID, bsonFieldName string) (string, error) {
	opts := options.Find()
	opts.SetLimit(1)

	id, err := primitive.ObjectIDFromHex(objID)
	if err != nil {
		return "", err
	}

	collection := client.Database(dbName).Collection(collName)
	results := []map[string]interface{}{}
	cursor, err := collection.Find(context.TODO(), bson.M{"_id": id}, opts)
	if err != nil {
		return "", err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		rowData := make(map[string]interface{})
		cursor.Decode(&rowData)
		results = append(results, rowData)
	}
	if err := cursor.Err(); err != nil {
		return "", err
	}

	strVal := ""
	for _, value := range results {
		strVal = fmt.Sprintf("%v", value[bsonFieldName])
	}
	return strVal, nil
}

// BareObjID removes the "ObjectID()" from the object id and returns the bare string value only.
func BareObjID(objectID string) string {
	res := strings.Replace(objectID, "ObjectID(", "", -1)
	res = strings.Replace(res, ")", "", -1)
	return res
}
