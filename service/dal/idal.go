package dal

import (
	"github.com/PeerIslands/aci-fx-go/config"
	"github.com/PeerIslands/aci-fx-go/model/entity"
	"github.com/gofiber/fiber/v2/log"
)

type DBService[T any] interface {
	Init(credentials ...string)
	GetOne(filter any) (T, error)
	GetOneById(id int) (T, error)
	Get(filter any) ([]T, error)
	CreateOne(document T) (T, error)
	BulkInsert(documents []T) (T, error)
	UpdateOne(document any, filter any) (any, error)
	UpdateOneById(id any) (any, error)
	DeleteOne(filter any) (int64, error)
}

func GetDataAccess(config *config.Config) DBService[entity.ForexData] {

	if config == nil {
		log.Fatal("No configuration found")
		return nil
	}

	if config.Db.Mongo.Url != "" {
		var db = MongoDbService[entity.ForexData]{}
		db.Init(config.Db.Mongo.Url)
		return &db
	}

	if config.Db.Yugabyte.Address != "" {
		var ydb = YugaByteDbService[entity.ForexData]{}
		ydb.Init(config.Db.Yugabyte.Username, config.Db.Yugabyte.Password, config.Db.Yugabyte.Dbname, config.Db.Yugabyte.Address, config.Db.Yugabyte.PoolSize)
		return &ydb
	}

	log.Fatal("No database configuration found")
	return nil
}
