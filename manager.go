package mongoconnector

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/leemcloughlin/logfile"
	"github.com/ninh0gauch0/hrstypes"
	log "github.com/sirupsen/logrus"
)

const (
	version = "0.8.3-beta"
)

var (
	baseContext   = context.Background()
	contextLogger *log.Entry
	logFileOn     = true
	logFile       *logfile.LogFile
)

// Init - always called at the begining
func (m *Manager) Init() bool {
	contextLogger.Infof("Initializing HR mongo connector...")
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

	/*

		//Reading configuration file
		dat, err := ioutil.ReadFile("config/mongo.json")
		if err != nil {
			customErrorLogger(s, "Failed to read configuration mongodb file %s: %s", "mongoconf", err.Error())
			return false
		}

		// Taking mongodb conf
		var result mongoCon.MongoConf
		err = json.Unmarshal(dat, &result)
		if err != nil {
			customErrorLogger(s, "Failed to unmarshal configuration json extracted from %s file: %s", "mongoconf", err.Error())
			return false
		}


	*/

	//filename is the path to the json config file
	file, err := os.Open("./config/mongo.json")
	if err != nil {
		return false
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(m.Conf)
	if err != nil {
		return false

	}
	m.Address = m.Conf.GetHost() + ":" + m.Conf.GetPort()
	m.initialized = true
	return true
}

// connect -
func (m *Manager) connect() error {

	if m.initialized != true {
		err := m.Init()
		if err {
			return fmt.Errorf("Error initializating mongo connector")
		}
	}

	m.logger.Infof("Starting mongo connector...")

	session, err := mgo.Dial(m.Address)
	if err != nil {
		return err
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	m.Session = session

	return nil
}

// CloseConnection -
func (m *Manager) CloseConnection() {

	if m.Session != nil {
		m.Session.Close()
	}
}

// ExecuteInsert -
func (m *Manager) ExecuteInsert(collection string, obj MetadataObject) (int, error) {

	err := m.connect()

	if err != nil {
		return -1, err
	}

	c := m.Session.DB(m.Conf.GetDB()).C(collection)

	// Insert Datas
	err = c.Insert(obj)
	if err != nil {
		return -1, err
	}

	return 0, nil
}

// ExecuteSearchByID -
func (m *Manager) ExecuteSearchByID(collection string, id string) (MetadataObject, error) {

	err := m.connect()

	if err != nil {
		return nil, err
	}

	c := m.Session.DB(m.Conf.GetDB()).C(collection)

	result := &hrstypes.Recipe{}
	err = c.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		m.customInfoLogger("No results for ID %s: %s", id, err.Error())
		return nil, err
	}
	return result, nil
}

// ExecuteSearch -
func (m *Manager) ExecuteSearch(collection string, query string) ([]MetadataObject, error) {
	var results []MetadataObject

	err := m.connect()

	if err != nil {
		return nil, err
	}

	c := m.Session.DB(m.Conf.GetDB()).C(collection)

	err = c.Find(nil).Sort("-id").All(&results)
	if err != nil {
		return nil, err
	}
	fmt.Println("Results All: ", results)
	return results, nil
}

// ExecuteUpdate -
func (m *Manager) ExecuteUpdate(collection string, id string, obj MetadataObject) (int, error) {

	err := m.connect()

	if err != nil {
		return -1, err
	}

	c := m.Session.DB(m.Conf.GetDB()).C(collection)
	colQuerier := bson.M{"id": id}
	change := bson.M{"$set": obj}
	err = c.Update(colQuerier, change)
	if err != nil {
		return -1, err
	}

	return 0, nil
}

// ExecuteDelete -
func (m *Manager) ExecuteDelete(collection string, id string) (int, error) {

	err := m.connect()

	if err != nil {
		return -1, err
	}

	c := m.Session.DB(m.Conf.GetDB()).C(collection)
	err = c.RemoveId(id)

	if err != nil {
		return -1, err
	}

	return 0, nil
}

// CustomErrorLogger - Writes error
func (m *Manager) customErrorLogger(msg string, args ...interface{}) {
	MSG := "[ERROR] " + msg

	m.logger.Errorf(MSG, args)
	m.logger.Errorln()

	if logFileOn {
		log.Printf(MSG, args)
	}
}

// customInfoLogger - Writes info
func (m *Manager) customInfoLogger(msg string, args ...interface{}) {
	MSG := "[INFO] " + msg

	m.logger.Infof(MSG, args)
	m.logger.Infoln()

	if logFileOn {
		log.Printf(MSG, args)
	}
}
