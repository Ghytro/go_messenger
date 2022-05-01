package handler

import (
	"database/sql"
	"log"
	"sync"

	"github.com/Ghytro/go_messenger/file_storage_service/interface/config"
	_ "github.com/lib/pq"
)

var fileDataDB *sql.DB
var counter = NewAtomicCounter(len(config.Config.StoragesAddrs))

type AtomicCounter struct {
	value    int
	maxValue int
	mutex    sync.Mutex
}

func NewAtomicCounter(maxValue int) *AtomicCounter {
	return &AtomicCounter{
		maxValue: maxValue,
	}
}

func (ac *AtomicCounter) Inc() int {
	ac.mutex.Lock()
	defer ac.mutex.Unlock()
	oldValue := ac.value
	ac.value++
	if ac.value > ac.maxValue {
		ac.value = 0
	}
	return oldValue
}

func init() {
	var err error
	if fileDataDB, err = sql.Open("postgres", "host=localhost user=postgres password=123123 dbname=file_data sslmode=disable"); err != nil {
		log.Fatal(err)
	}
	if _, err := fileDataDB.Exec(
		`CREATE TABLE IF NOT EXISTS file_storages (
			id CHAR(50) NOT NULL PRIMARY KEY,
			storage_addr VARCHAR(2048) NOT NULL,
			content_type VARCHAR(30) NOT NULL
		)`,
	); err != nil {
		log.Fatal(err)
	}
}
