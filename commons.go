package mongoconnector

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/ninh0gauch0/mongoconnector/types"
	log "github.com/sirupsen/logrus"
)

// LoggerTrait - a logger trait that let's you configure a log
type LoggerTrait struct {
	logger *log.Entry
}

// MetadataObject interface
type MetadataObject interface {
	GetObjectInfo() string
}

// Manager struct
type Manager struct {
	LoggerTrait
	Ctx         context.Context
	Session     *mgo.Session
	initialized bool
	Address     string
	Conf        *types.MongoConf
}
