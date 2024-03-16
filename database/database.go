package database

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pentag.kr/dimimonster/config"
)

func InitMDB() {
	// Setup the mgm default config
	err := mgm.SetDefaultConfig(nil, "dimimonster", options.Client().ApplyURI(config.DB_URI))
	if err != nil {
		panic(err)
	}
}
