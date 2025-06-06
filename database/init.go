package database

import (
	"github.com/trancecho/open-sdk/config"
	"gorm.io/gorm"
	"log"
	"sync"
)

var (
	dbs = make(map[string]*gorm.DB)
	mux sync.RWMutex
)

func InitDB() {
	sources := config.GetConfig().Databases
	for _, source := range sources {
		setDbByKey(source.Key, mustCreateGorm(source))
		if source.Key == "" {
			source.Key = "*"
		}
		//logx.NameSpace("Dbx").Infoln("create datasource %s => %s:%s", source.Key, source.IP, source.PORT)
		log.Println("create datasource", source.Key, "=>", source.IP, ":", source.PORT)
	}
}

func GetDb(key string) *gorm.DB {
	mux.Lock()
	defer mux.Unlock()
	return dbs[key]
}

func setDbByKey(key string, db *gorm.DB) {
	if key == "" {
		key = "*"
	}
	if GetDb(key) != nil {
		log.Fatalln("duplicate db key: ", key)
	}
	mux.Lock()
	defer mux.Unlock()
	dbs[key] = db
}

func mustCreateGorm(database config.Datasource) *gorm.DB {
	var creator = getCreatorByType(database.Type)
	if creator == nil {
		log.Fatalln("fail to find creator for types:", database.Type)
		return nil
	}
	db, err := creator.Create(database.IP, database.PORT, database.USER, database.PASSWORD, database.DATABASE)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	return db
}
