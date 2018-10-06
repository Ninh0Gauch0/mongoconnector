package mongoconnector

// MongoConf DTO
type MongoConf struct {
	host string
	port string
	db   string
}

// GetHost - returns the host
func (mc *MongoConf) GetHost() string {
	return mc.host
}

// GetPort - returns the port
func (mc *MongoConf) GetPort() string {
	return mc.port
}

// GetDB - returns the db
func (mc *MongoConf) GetDB() string {
	return mc.db
}
