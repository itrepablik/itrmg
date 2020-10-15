![itrmg](https://user-images.githubusercontent.com/58651329/80572103-4772f180-8a30-11ea-8e4b-19c79ccb753d.png)
The **itrmg** package is a simplified usage of MongoDB official library for Go.

# Installation
```
go get -u github.com/itrepablik/itrmg
```

# Usage
These are the various examples on how you can use the **itrmg** package for your next Go project with MongoDB database.
```
package main

import (
	"fmt"
	"time"

	"github.com/itrepablik/itrlog"
	"github.com/itrepablik/itrmg"
)

const (
	_dbConStr = "mongodb://username:password.@localhost:27017"
	_dbName   = "database_name"
	_collName = "collection_name"
)

// ClientMG set the MongoDB's client variable to be called globally across your Go's project
var ClientMG = itrmg.ClientMG

func main() {
	// Initialize the MongoDB connection
	ClientMG, err := itrmg.InitMG(_dbConStr)
	if err != nil {
		itrlog.Fatal(err)
	}

	// itrmg.FindOneByID usage: find single row by using object id in MongoDB collection
	data, err := itrmg.FindOneByID(_dbName, _collName, ClientMG, "5e2a59e3f1a9a91790a13c37")
	if err != nil {
		itrlog.Fatal(err)
	}

	// Get the specific column name with it's value
	fmt.Println("data: ", data["pc_name"])

	// itrmg.Find usage: use filter, sort order and row limit for your query.
	filter := itrmg.DP{"created_by": "politz", "status": "Used"} // constract your filter query here
	sortOrder := itrmg.DP{"pc_name": -1}                         // descending sort order

	results, err := itrmg.Find(_dbName, _collName, ClientMG, filter, sortOrder, 2)
	if err != nil {
		itrlog.Fatal(err)
	}

	// Iterate the 'results' from the itrmg.Find query
	for _, value := range results {
		fmt.Println(value["pc_name"])
	}

	// InsertOne usage: this will insert one row to your collection
	newRow := itrmg.DP{
		"pc_name":      "pc name 1234",
		"license":      "abc 123",
		"price":        23,
		"ip_address":   "123456",
		"created_by":   "politz",
		"created_date": time.Now(),
		"status":       "Available",
		"is_active":    true,
	}

	isInserted, err := itrmg.InsertOne(_dbName, _collName, ClientMG, newRow)
	if err != nil {
		itrlog.Error(err)
	}

	if isInserted {
		fmt.Println("New row has been inserted!")
	}

	// UpdateOne usage: update only one row to your MongoDB collection
	updateFilter := itrmg.DP{"license": "abc 123", "created_by": "politz"}
	updateRow := itrmg.DP{
		"pc_name":       "pc name 888",
		"license":       "abc 111",
		"price":         30,
		"ip_address":    "123456",
		"modified_by":   "politz",
		"modified_date": time.Now(),
	}

	isUpdated, err := itrmg.UpdateOne(_dbName, _collName, ClientMG, updateRow, updateFilter)
	if err != nil {
		itrlog.Error(err)
	}

	if isUpdated {
		fmt.Println("Row has been modified successfully!")
	}

	// UpdateOneByID usage: single row update by using object id
	objID := "5eab87c0fcde9804abc5fbc9"
	updateRow1 := itrmg.DP{
		"pc_name":       "pc name 000",
		"license":       "abc 000",
		"price":         0,
		"ip_address":    "123456",
		"modified_by":   "politz",
		"modified_date": time.Now(),
	}

	isUpdated1, err := itrmg.UpdateOneByID(_dbName, _collName, ClientMG, updateRow1, objID)
	if err != nil {
		itrlog.Error(err)
	}

	if isUpdated1 {
		fmt.Println("Row has been modified successfully!")
	}

	// DeleteOneByID usage: delete any single row permanently filtered by object id.
	rowObjID := "5eab87c0fcde9804abc5fbc9"
	isDeleted, err := itrmg.DeleteOneByID(_dbName, _collName, ClientMG, rowObjID)
	if err != nil {
		itrlog.Error(err)
	}

	if isDeleted {
		fmt.Println("Row has been deleted successfully!")
	}
	
	// IsExist usage: check if specific information found in certain collection
	filter = itrmg.DP{"created_by": "politz"}
	isFound, err := itrmg.IsExist(_dbName, _collName, ClientMG, filter)
	if err != nil {
		itrlog.Error(err)
	}
	if isFound {
		fmt.Println("Record found!")
	}
	
	// GetFieldValue usage: get the specific field value from a certain collection.
	filter := itrmg.DP{"pc_name": "PPP Name"}
	strVal, err := itrmg.GetFieldValue(_dbName, _collName, ClientMG, filter, "license")
	if err != nil {
		itrlog.Error(err)
	}
	fmt.Println("strVal: ", strVal)
}
```
Optionally, if you've a struct for your data structure, you can specify the object id as "itrmg.ObjID".
```
package models

import "github.com/itrepablik/itrmg"

// YourDataStruct is a collection of your own data structure here.
type YourDataStruct struct {
	ID       itrmg.ObjID `json:"_id" bson:"_id"`
	TypeName string      `json:"type_name" bson:"type_name"`
	IsActive bool        `json:"is_active" bson:"is_active"`
}
```

# Subscribe to Maharlikans Code Youtube Channel:
Please consider subscribing to my Youtube Channel to recognize my work on this package. Thank you for your support!
https://www.youtube.com/channel/UCdAVUmldU9Jn2VntuQChHqQ/

# License
Code is distributed under MIT license, feel free to use it in your proprietary projects as well.
