package types

// MongoConf DTO
type MongoConf struct {
	Host string `json:"host"`
	Port string `json:"port"`
	DB   string `json:"db"`
}

// GetHost - returns the host
func (mc *MongoConf) GetHost() string {
	return mc.Host
}

// SetHost - sets the host
func (mc *MongoConf) SetHost(host string) {
	mc.Host = host
}

// GetPort - returns the port
func (mc *MongoConf) GetPort() string {
	return mc.Port
}

// SetPort - sets the port
func (mc *MongoConf) SetPort(port string) {
	mc.Port = port
}

// GetDB - returns the db
func (mc *MongoConf) GetDB() string {
	return mc.DB
}

// SetDB - sets the db
func (mc *MongoConf) SetDB(db string) {
	mc.DB = db
}
