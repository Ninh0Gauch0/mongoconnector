package mongoconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/leemcloughlin/logfile"
	"github.com/ninh0gauch0/hrstypes"
	log "github.com/sirupsen/logrus"
)

const (
	version = "0.9.1-beta"
	// INGREDIENTCOLL Constant
	INGREDIENTCOLL = "ingredients"
	// RECIPECOLL Constant
	RECIPECOLL = "recipes"
)

var (
	baseContext   = context.Background()
	contextLogger *log.Entry
	logFileOn     = true
	logFile       *logfile.LogFile
)

// Init - always called at the begining
func (m *Manager) Init() bool {

	// Contex configuration
	mongoContext, cancelFunc := context.WithCancel(baseContext)
	m.Ctx = mongoContext

	// logger configuration
	logger := log.StandardLogger()
	logger.Formatter = &log.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	}
	logger.Out = os.Stdout
	logger.SetLevel(log.InfoLevel)

	contextLogger = logger.WithFields(log.Fields{
		"mongo connector": "Home Recipes DB",
	})
	m.SetLogger(contextLogger)
	m.logger.Infof("Initializing HR mongo connector...")
	defer cancelFunc()

	// Init logfile
	logFile, err := logfile.New(
		&logfile.LogFile{
			FileName: "hrMongoConnector.log",
			MaxSize:  1000 * 1024,
			Flags:    logfile.FileOnly | logfile.RotateOnStart})
	if err != nil {
		m.logger.Errorf("Failed to create log file %s: %s", logFile.FileName, err.Error())
		logFileOn = false
	}
	log.SetOutput(logFile)

	//Reading configuration file
	dat, err := ioutil.ReadFile("config/mongo.json")
	if err != nil {
		m.customErrorLogger("Failed to read configuration mongodb file mongoconf: ", err.Error())
		return false
	}
	conf := hrstypes.MongoConf{}
	// Taking mongodb conf
	err = json.Unmarshal(dat, &conf)
	if err != nil {
		m.customErrorLogger("Failed to unmarshal configuration json extracted from %s file: %s", "mongoconf", err.Error())
		return false
	}
	m.Conf = &conf
	m.Address = m.Conf.GetHost() + ":" + m.Conf.GetPort()
	m.initialized = true
	return true
}

// getCollection - Given a mongo session, returns a collection instance
func (m *Manager) getCollection(session mgo.Session, db string, coll string) *mgo.Collection {

	fmt.Println("Creating collection ...")
	collection := session.Copy().DB(db).C(coll)
	fmt.Println("Collection created!")
	return collection

}

func (m *Manager) getSession(db string) (*mgo.Session, error) {
	fmt.Println("Creating session ...")

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{m.Address},
		Timeout:  5 * time.Second,
		Database: db,
	}

	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)

	if err != nil {
		return nil, err
	}

	mongoSession.SetSocketTimeout(5 * time.Second)
	mongoSession.SetMode(mgo.Monotonic, true)

	session := mongoSession.New()

	fmt.Println("Session created!")
	return session, nil
}

// connect -
func (m *Manager) connect(coll string) (*mgo.Collection, error) {

	if m.initialized != true {
		err := m.Init()
		if err {
			return nil, fmt.Errorf("Error initializating mongo connector")
		}
	}

	m.logger.Infof("Connecting to %s database...", m.Conf.GetDB())

	session, err := m.getSession(m.Conf.GetDB())
	if err != nil {
		return nil, err
	}
	defer session.Close()

	collection := m.getCollection(*session, m.Conf.GetDB(), coll)

	return collection, nil
}

// ExecuteInsert -
func (m *Manager) ExecuteInsert(collection string, obj hrstypes.MetadataObject) (int, error) {

	c, err := m.connect(collection)

	if err != nil {
		return -1, err
	}

	// Insert Datas
	err = c.Insert(obj)
	if err != nil {
		return -1, err
	}
	return 0, nil
}

// ExecuteSearchByID -
func (m *Manager) ExecuteSearchByID(collection string, id string) (hrstypes.MetadataObject, error) {
	c, err := m.connect(collection)

	if err != nil {
		return nil, err
	}

	switch collection {
	case INGREDIENTCOLL:
		ingredient := &hrstypes.Ingredient{}
		err = c.Find(bson.M{"_id": id}).One(&ingredient)

		if err != nil {
			m.customInfoLogger("No ingredients were found for ID %s: %s", id, err.Error())
			return nil, err
		}
		return ingredient, nil
	case RECIPECOLL:
		fallthrough
	default:
		recipe := &hrstypes.Recipe{}
		err = c.Find(bson.M{"_id": id}).One(&recipe)

		if err != nil {
			m.customInfoLogger("No recipes were found for ID %s: %s", id, err.Error())
			return nil, err
		}
		return recipe, nil
	}
}

// ExecuteSearch -
func (m *Manager) ExecuteSearch(collection string, query string) ([]hrstypes.MetadataObject, error) {
	var results []hrstypes.MetadataObject

	c, err := m.connect(collection)

	if err != nil {
		return nil, err
	}

	err = c.Find(nil).Sort("_id").All(&results)
	if err != nil {
		return nil, err
	}
	fmt.Println("Results All: ", results)
	return results, nil
}

// ExecuteUpdate -
func (m *Manager) ExecuteUpdate(collection string, id string, obj hrstypes.MetadataObject) (hrstypes.MetadataObject, error) {

	c, err := m.connect(collection)
	if err != nil {
		return nil, err
	}
	colQuerier := bson.M{"_id": id}
	change := mgo.Change{
		Update:    nil,
		Upsert:    false,
		Remove:    false,
		ReturnNew: true,
	}

	s := reflect.ValueOf(&obj).Elem()
	typeOfT := s.Type()
	var patchDesc, patchName string

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, field.Type(), field.Interface())
		if typeOfT.Field(i).Name == "Description" {
			patchDesc = field.Interface().(string)
		} else if typeOfT.Field(i).Name == "Name" {
			patchName = field.Interface().(string)
		}
	}
	// Building the update bson
	if patchDesc != "" && patchName != "" {
		change.Update = bson.M{"$set": bson.M{"name": patchName, "description": patchDesc}}
	} else if patchDesc != "" {
		change.Update = bson.M{"$set": bson.M{"description": patchDesc}}
	} else if patchName != "" {
		change.Update = bson.M{"$set": bson.M{"name": patchName}}
	}

	if _, ok := obj.(*hrstypes.Ingredient); ok {
		ingredient := &hrstypes.Ingredient{}

		// err = c.Update(colQuerier, change)
		// if err != nil {
		// 	return -1, err
		// }

		//var doc *mgo.ChangeInfo
		_, err = c.Find(colQuerier).Apply(change, &ingredient)
		/////////
		if err != nil {
			m.customInfoLogger("No ingredients were found for ID %s: %s", id, err.Error())
			return nil, err
		}
		return ingredient, nil
	}
	recipe := &hrstypes.Recipe{}
	_, err = c.Find(bson.M{"_id": id}).Apply(change, &recipe)

	if err != nil {
		m.customInfoLogger("No recipes were found for ID %s: %s", id, err.Error())
		return nil, err
	}
	return recipe, nil

}

// ExecuteDelete -
func (m *Manager) ExecuteDelete(collection string, id string) (int, error) {
	c, err := m.connect(collection)

	if err != nil {
		return -1, err
	}
	err = c.Remove(bson.M{"_id": id})

	if err != nil {
		return -1, err
	}
	return 0, nil
}

// CustomErrorLogger - Writes error
func (m *Manager) customErrorLogger(msg string, args ...interface{}) {
	MSG := "[ERROR] " + msg

	if logFileOn {
		log.Printf(MSG, args)
	}
}

// customInfoLogger - Writes info
func (m *Manager) customInfoLogger(msg string, args ...interface{}) {
	MSG := "[INFO] " + msg

	if logFileOn {
		log.Printf(MSG, args)
	}
}
