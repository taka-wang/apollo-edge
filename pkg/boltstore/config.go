package boltstore

// ConfigType config file structure with default values
type ConfigType struct {
	BoltStore struct {
		BucketName  string `default:"MQTT"`
		OpenTimeout int    `default:"1"` // boltdb open timeout
		Logger      struct {
			// Level: [Panic : 0, Fatal : 1, Error : 2, Warn  : 3, Info  : 4, Debug : 5]
			Level    int    `default:"4"`
			JSON     bool   `default:"false"`
			ToFile   bool   `default:"false"`
			Filename string `default:"boltstore.log"`
		}
	}
}
