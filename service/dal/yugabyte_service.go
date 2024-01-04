package dal

import (
	"context"
	"errors"
	"github.com/PeerIslands/aci-fx-go/model/dto/request"
	"github.com/go-pg/pg/v10"
	"log"
	"strconv"
)

type YugaByteDbService[T any] struct {
	YbDB *pg.DB
}

var ybDB *pg.DB

func (y *YugaByteDbService[T]) Init(credentials ...string) {
	poolSize, _ := strconv.Atoi(credentials[4])
	ybDB = pg.Connect(&pg.Options{
		User:     credentials[0],
		Password: credentials[1],
		Database: credentials[2],
		Addr:     credentials[3],
		PoolSize: poolSize,
		//TLSConfig: &tls.Config{
		//	InsecureSkipVerify: true,
		//},
	})

	// Check the connection
	ctx := context.Background()
	if err := ybDB.Ping(ctx); err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
	log.Println("Connected to database", credentials[2], "on", credentials[3], "with pool size", poolSize)
	y.YbDB = ybDB
}

func (y *YugaByteDbService[T]) getDatabase() *pg.DB {
	return ybDB
}

func (y *YugaByteDbService[T]) GetOne(filter any) (T, error) {
	var fxRequest = filter.(request.FxDataRequest)
	var data T
	err := y.YbDB.Model(&data).
		Where("tenant_id = ?", fxRequest.TenantId).
		Where("bank_id = ?", fxRequest.BankId).
		Where("base_currency = ?", fxRequest.BaseCurrency).
		Where("target_currency = ?", fxRequest.TargetCurrency).
		Where("tier = ?", fxRequest.Tier).First()
	if err != nil {
		return data, err
	}
	return data, nil
}
func (y *YugaByteDbService[T]) GetOneById(id int) (T, error) {

	var data T
	err := y.YbDB.Model(&data).
		Where("id = ?", id).
		First()
	if err != nil {
		return data, err
	}
	return data, nil
}

func (y *YugaByteDbService[T]) CreateOne(record T) (T, error) {
	_, err := y.YbDB.Model(&record).Insert()
	if err != nil {
		return record, err
	}
	return record, nil
}

func (y *YugaByteDbService[T]) UpdateOne(record any, filter any) (any, error) {
	var fxRequest = filter.(request.FxDataRequest)

	var rec T
	_, err := y.YbDB.Model(&rec).
		Set("buy_rate = buy_rate + ?", 0.001).
		Where("tenant_id = ?", fxRequest.TenantId).
		Where("bank_id = ?", fxRequest.BankId).
		Where("base_currency = ?", fxRequest.BaseCurrency).
		Where("target_currency = ?", fxRequest.TargetCurrency).
		Where("tier = ?", fxRequest.Tier).
		Returning("*").
		Update()
	if err != nil {
		return record, err
	}
	return record, nil
}

func (y *YugaByteDbService[T]) UpdateOneById(id any) (any, error) {
	var rowId = id.(int)

	var rec T
	_, err := y.YbDB.Model(&rec).
		Set("buy_rate = buy_rate + ?", 0.001).
		Where("id = ?", rowId).
		Returning("*").
		Update()
	if err != nil {
		return rec, err
	}
	return rec, nil
}

func (y *YugaByteDbService[T]) DeleteOne(id any) (int64, error) {
	var data []T
	result, err := y.YbDB.Model(&data).Where("id = ?", id.(string)).Delete()
	if err != nil {
		return 0, err
	}
	if result.RowsAffected() == 0 {
		return 0, errors.New("no record found")
	}
	return 1, nil
}

func (y *YugaByteDbService[T]) Get(filter any) ([]T, error) {
	var data []T
	err := y.YbDB.Model(&data).Where(filter.(string)).Select()
	if err != nil {
		return data, err
	}
	return data, nil
}

func (y *YugaByteDbService[T]) BulkInsert(documents []T) (T, error) {
	_, err := y.YbDB.Model(&documents).Insert()
	if err != nil {
		return documents[0], err
	}
	return documents[0], nil
}
