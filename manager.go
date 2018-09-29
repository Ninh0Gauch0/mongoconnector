package mongoconnector

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/ninh0gauch0/hrstypes"
)

var (
	session *mgo.Session
)

// Connect -
func Connect() error {
	session, err := mgo.Dial("127.0.0.1")
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
func ExecuteInsert(db string, collection string, obj hrstypes.MetadataObject) (int, error) {

	c := session.DB("hrs").C("recipes")

	// Insert Datas
	err := c.Insert(obj)
	if err != nil {
		return -1, err
	}

	return 0, nil
}

// ExecuteSearchOne -
func ExecuteSearchOne() (hrstypes.MetadataObject, error) {
	c := session.DB("hrs").C("recipes")

	result := hrstypes.Recipe{}
	err := c.Find(bson.M{"name": "Ale"}).Select(bson.M{"phone": 0}).One(&result)
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
func ExecuteSearch(queryTimeout int) ([]hrstypes.MetadataObject, error) {
	var results []hrstypes.MetadataObject
	c := session.DB("hrs").C("recipes")

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
func ExecuteUpdate() int {
	c := session.DB("hrs").C("recipes")
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
