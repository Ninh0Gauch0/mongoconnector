package mongoconnector

import (
	"context"

	"github.com/ninh0gauch0/hrstypes"
	log "github.com/sirupsen/logrus"
)

// LoggerTrait - a logger trait that let's you configure a log
type LoggerTrait struct {
	logger *log.Entry
}

// SetLogger - let's you configure a logger
func (lt *LoggerTrait) SetLogger(l *log.Entry) {
	if l != nil {
		lt.logger = l
	}
}

// GetLogger - returns the logger
func (lt *LoggerTrait) GetLogger() *log.Entry {
	return lt.logger
}

// Manager struct
type Manager struct {
	LoggerTrait
	Ctx         context.Context
	initialized bool
	Address     string
	Conf        *hrstypes.MongoConf
}
