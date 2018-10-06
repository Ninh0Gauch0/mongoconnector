package mongoconnector

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ninh0gauch0/hrstypes"
)

var (
	session *mgo.Session
	address string
	db      string
	coll    string
	conf    *MongoConf
)

// Init - always called at the begining
func Init() bool {
	//filename is the path to the json config file
	file, err := os.Open("./config/mongo.json")
	if err != nil {
		return false
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		return false
	}
	return true
}

// Connect -
func Connect() error {
	session, err := mgo.Dial(address)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	return nil
}

// CloseConnection -
func CloseConnection() {
	if session != nil {
		session.Close()
	}
}

// ExecuteInsert -
func ExecuteInsert(collection string, obj hrstypes.MetadataObject) (int, error) {

	c := session.DB(db).C(collection)

	// Insert Datas
	err := c.Insert(obj)
	if err != nil {
		return -1, err
	}

	return 0, nil
}

// ExecuteSearchByID -
func ExecuteSearchByID(collection string, id string) (hrstypes.MetadataObject, error) {
	c := session.DB(db).C(collection)

	result := hrstypes.Recipe{}
	err := c.Find(bson.M{"ID": id}).One(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

/*
// ExecuteSearchOne -
func ExecuteSearchOne(queryTimeout int){

}
*/

// ExecuteSearch -
func ExecuteSearch(collection string, queryTimeout int) ([]hrstypes.MetadataObject, error) {
	var results []hrstypes.MetadataObject
	c := session.DB(db).C(collection)

	err := c.Find(nil).Sort("-timestamp").All(&results)
	if err != nil {
		return nil, err
	}
	fmt.Println("Results All: ", results)
	return results, nil
}

/*

// ExecuteSearch -
func ExecuteSearch() int{
	return 0
}
*/

// ExecuteUpdate -
func ExecuteUpdate(collection string) int {
	c := session.DB(db).C(collection)
	colQuerier := bson.M{"name": "Ale"}
	change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": time.Now()}}
	err := c.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}

	return 0
}

/*

// ExecuteUpdate -
func executeUpdate(queryTimeout int) {
	return 0
}
*/

// ExecuteDelete -
func ExecuteDelete() int {
	return 0
}

/**
 * Execute delete.
 *
 * @param delete the delete
 * @return the long
 */
// public abstract long executeDelete(SearchableObject delete);

/**
 * Execute update.
 *
 * @param find the find
 * @param update the update
 * @return the long
 */
// public abstract long executeUpdate(SearchableObject find, MetadataObject update);

/**
 * Execute update.
 *
 * @param find the find
 * @param update the update
 * @param queryTimeout the query timeout
 * @return the long
 */
//public abstract long executeUpdate(SearchableObject find, MetadataObject update, long queryTimeout);

/**
 * Execute bulk insert.
 *
 * @param objects the objects
 */
//public abstract void executeBulkInsert(List<MetadataObject> objects);

//public abstract List<OPLogResume> showOplog(long queryTimeout);

//public abstract List<OPLogResume> showOplog();
