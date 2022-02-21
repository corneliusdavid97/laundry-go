// Package postgresql provide mechanism to interact with postgres DB
package postgresql

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"
)

type database struct {
	name        string
	replication databaseReplication
}

type databaseReplication struct {
	master *sqlx.DB
	slaves []*sqlx.DB
}

type databaseConfig struct {
	Name   string   `yaml:"name"`
	Master string   `yaml:"master"`
	Slaves []string `yaml:"slaves"`
}

const postgres = "postgres"

var globalLock = sync.RWMutex{}
var databaseMap map[string]database

// InitPostgresqlConfig initialize postgres databases in config file
func InitPostgresqlConfig(ctx context.Context, basepath string) error {
	data := struct {
		Data []databaseConfig `yaml:"data"`
	}{}

	filepath := basepath + "/etc/config/database/postgresql.yaml"
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("%s, path: %s", err.Error(), filepath)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		return errors.New("failed to parse config")
	}

	globalLock.Lock()
	databaseMap = make(map[string]database, len(data.Data))
	for _, db := range data.Data {
		master, err := sqlx.ConnectContext(ctx, postgres, db.Master)
		if err != nil {
			// add log without return
			return err
		}
		var slaves []*sqlx.DB
		for _, slaveStr := range db.Slaves {
			slave, err := sqlx.ConnectContext(ctx, postgres, slaveStr)
			if err != nil {
				// add log without return
				return err
			}
			slaves = append(slaves, slave)
		}
		databaseMap[db.Name] = database{
			name: db.Name,
			replication: databaseReplication{
				master: master,
				slaves: slaves,
			},
		}
	}
	globalLock.Unlock()

	return nil
}

var errInvalidParam = errors.New("Invalid parameter input")
var errDBNotExist = errors.New("Requested DB does not exist")

func GetDB(dbName, replication string) (*sqlx.DB, error) {
	var err error
	if dbName == "" || replication == "" {
		return nil, fmt.Errorf("%w: dbName or replication is empty", errInvalidParam)
	}

	if db, ok := databaseMap[dbName]; ok {
		if replication == "master" {
			return db.replication.master, nil
		}
		if strings.Contains(replication, "slave") {
			s := strings.Split(replication, "_")
			idx := 1
			if len(s) > 1 {
				idx, err = strconv.Atoi(s[1])
				if err != nil {
					return nil, fmt.Errorf("%w: replication param invalid", errInvalidParam)
				}
			}
			if len(db.replication.slaves) > idx-1 {
				return db.replication.slaves[idx-1], nil
			}
			return nil, errDBNotExist
		}
		return nil, fmt.Errorf("%w: replication param invalid", errInvalidParam)
	}
	return nil, errDBNotExist
}
