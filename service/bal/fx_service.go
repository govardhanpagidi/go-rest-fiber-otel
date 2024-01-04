package bal

import (
	"context"
	"github.com/PeerIslands/aci-fx-go/model/dto/request"
	"github.com/PeerIslands/aci-fx-go/model/dto/response"
	"github.com/PeerIslands/aci-fx-go/model/entity"
	"github.com/PeerIslands/aci-fx-go/service/common"
	"github.com/PeerIslands/aci-fx-go/service/dal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel"
	"math/rand"
	"os"
	"time"
)

type Fx_service struct {
	DbService dal.DBService[entity.ForexData]
}

var dbSpanName = "db-call"
var tracerName = "OTEL_SERVICE_NAME"

func getForexDtoFromEntity(result entity.ForexData) *response.ForexDataResponse {
	return &response.ForexDataResponse{
		Id:                           result.ID,
		TenantId:                     result.TenantID,
		BankId:                       result.BankID,
		BaseCurrency:                 result.BaseCurrency,
		TargetCurrency:               result.TargetCurrency,
		Tier:                         result.Tier,
		DirectIndirectFlag:           result.DirectIndirectFlag,
		Multiplier:                   result.Multiplier,
		BuyRate:                      result.BuyRate,
		SellRate:                     result.SellRate,
		TolerancePercentage:          result.TolerancePercentage,
		EffectiveDate:                result.EffectiveDate,
		ExpirationDate:               result.ExpirationDate,
		ContractRequirementThreshold: result.ContractRequirementThreshold,
		DocVersion:                   result.DocVersion,
	}
}

func (s *Fx_service) CreateForexData(c *context.Context,
	forexData request.CreateForexDataRequest) response.ResponseWithSimpleData[response.CreateForexDataResponse] {
	var dbObject = entity.ForexData{
		ID:                           primitive.NewObjectID(),
		Tier:                         forexData.Tier,
		DirectIndirectFlag:           forexData.DirectIndirectFlag,
		Multiplier:                   forexData.Multiplier,
		BuyRate:                      forexData.BuyRate,
		SellRate:                     forexData.SellRate,
		TolerancePercentage:          forexData.TolerancePercentage,
		EffectiveDate:                forexData.EffectiveDate,
		ExpirationDate:               forexData.ExpirationDate,
		ContractRequirementThreshold: forexData.ContractRequirementThreshold,
		TenantID:                     forexData.TenantId,
		BankID:                       forexData.BankId,
		BaseCurrency:                 forexData.BaseCurrency,
		TargetCurrency:               forexData.TargetCurrency,
		CreatedDate:                  time.Now(),
		DocVersion:                   1,
		UpdatedDate:                  time.Now(),
	}
	common.Logger.Info("Create a forex record started")
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)
	result, err := s.DbService.CreateOne(dbObject)
	span.End()

	common.Logger.Info("Create a forex record ended")
	if err != nil {
		common.Logger.Errorf("Error in creating a new Record. Exception:%v", err)
		e := &[]response.Error{
			{Code: "FAILURE", Message: "Unable to create record", Details: "Unable to create record due to some exception."},
		}
		return common.GetSimpleResponse[response.CreateForexDataResponse](nil, response.InternalError, e)
	}

	return common.GetSimpleResponse[response.CreateForexDataResponse](&response.CreateForexDataResponse{Id: result.ID}, response.Success, nil)
}

func (s *Fx_service) BulkInsertForexData(c *context.Context,
	forexData []request.CreateForexDataRequest) response.ResponseWithSimpleData[response.ForexDataResponse] {
	var dbObjects []entity.ForexData
	for _, item := range forexData {
		dbObjects = append(dbObjects, entity.ForexData{
			ID:                           primitive.NewObjectID(),
			Tier:                         item.Tier,
			DirectIndirectFlag:           item.DirectIndirectFlag,
			Multiplier:                   item.Multiplier,
			BuyRate:                      item.BuyRate,
			SellRate:                     item.SellRate,
			TolerancePercentage:          item.TolerancePercentage,
			EffectiveDate:                item.EffectiveDate,
			ExpirationDate:               item.ExpirationDate,
			ContractRequirementThreshold: item.ContractRequirementThreshold,
			TenantID:                     item.TenantId,
			BankID:                       item.BankId,
			BaseCurrency:                 item.BaseCurrency,
			TargetCurrency:               item.TargetCurrency,
			CreatedDate:                  time.Now(),
			DocVersion:                   randomInt(),
			UpdatedDate:                  time.Now(),
		})
	}
	common.Logger.Info("Bulk insert started")
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)

	_, err := s.DbService.BulkInsert(dbObjects)
	span.End()
	common.Logger.Info("Bulk insert ended")
	if err != nil {
		common.Logger.Errorf("Error in creating a new Record. Exception:%v", err)
		e := &[]response.Error{
			{Code: "FAILURE", Message: "Unable to create record", Details: "Unable to create record due to some exception."},
		}
		return common.GetSimpleResponse[response.ForexDataResponse](nil, response.InternalError, e)
	}

	return common.GetSimpleResponse[response.ForexDataResponse](nil, response.Success, nil)
}

func (s *Fx_service) GetForexRateById(c *context.Context,
	id string) response.ResponseWithSimpleData[response.ForexDataResponse] {
	objectId, _ := primitive.ObjectIDFromHex(id)
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)

	result, err := s.DbService.GetOne(bson.D{{"_id", objectId}})
	span.End()

	if err != nil {
		common.Logger.Errorf("Error in retriving forex rate by id. Exception:%v", err)
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		return common.GetSimpleResponse[response.ForexDataResponse](nil, response.NotFound, e)
	}

	return common.GetSimpleResponse[response.ForexDataResponse](getForexDtoFromEntity(result), response.Success, nil)
}

func (s *Fx_service) DeleteForexRateById(c *context.Context,
	id string) response.ResponseWithSimpleData[response.ForexDataResponse] {
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)

	_, err := s.DbService.DeleteOne(id)
	span.End()
	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record Deleted", Details: "No record Deleted"},
		}
		common.Logger.Errorf("Error in Deleting forex rate by id. Exception:%v", err)
		return common.GetSimpleResponse[response.ForexDataResponse](nil, response.NotFound, e)
	}

	return common.GetSimpleResponse[response.ForexDataResponse](nil, response.Success, nil)
}

func (s *Fx_service) GetForexRateByFilter(c *context.Context,
	tenantId int, bankId int, baseCurrency string, targetCurrency string) response.ResponseWithArrayData[response.ForexDataResponse] {
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)

	result, err := s.DbService.Get(bson.D{
		{"tenantId", tenantId},
		{"bankId", bankId},
		{"baseCurrency", baseCurrency},
		{"targetCurrency", targetCurrency},
	})
	span.End()

	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		common.Logger.Errorf("Error in retriving forex rate by filters. Exception:%v", err)
		return common.GetArrayResponse[response.ForexDataResponse](nil, response.NotFound, e)
	}

	var data []response.ForexDataResponse

	for _, item := range result {
		data = append(data, *getForexDtoFromEntity(item))
	}

	return common.GetArrayResponse[response.ForexDataResponse](&data, response.Success, nil)

}

func (s *Fx_service) UpdateForexRateById(c *context.Context,
	id string, body request.UpdateForexDataRequest) response.ResponseWithSimpleData[response.ForexDataResponse] {

	objectId, _ := primitive.ObjectIDFromHex(id)

	updateDocument := bson.D{
		{"$set", bson.D{
			{"sellRate", body.SellRate},
			{"buyRate", body.BuyRate},
			{"expirationDate", body.ExpirationDate},
			{"effectiveDate", body.EffectiveDate},
			{"tolerancePercentage", body.TolerancePercentage},
			{"multiplier", body.Multiplier},
			{"directIndirectFlag", body.DirectIndirectFlag},
			{"contractRequirementThreshold", body.ContractRequirementThreshold},
			{"updatedDate", time.Now()},
		}},
	}
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)

	result, err := s.DbService.UpdateOne(updateDocument, bson.D{{"_id", objectId}})
	span.End()

	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record Updated", Details: "No record Updated"},
		}
		common.Logger.Errorf("Error in updating forex rate by id. Exception:%v", err)
		return common.GetSimpleResponse[response.ForexDataResponse](nil, response.NotFound, e)
	}

	return common.GetSimpleResponse[response.ForexDataResponse](getForexDtoFromEntity(result.(entity.ForexData)), response.Success, nil)
}

func (s *Fx_service) GetConvertedRate(c *context.Context,
	tenantId int, bankId int, amount float64, baseCurrency string, targetCurrency string, tier string) response.ResponseWithSimpleData[response.ConversionResponse] {
	ConvertRequest := request.FxDataRequest{
		TenantId:       tenantId,
		BankId:         bankId,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Tier:           tier,
	}
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)

	result, err := s.DbService.GetOne(ConvertRequest)
	span.End()
	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		common.Logger.Errorf("Error in retriving and converting forex rate. Exception:%v", err)
		return common.GetSimpleResponse[response.ConversionResponse](nil, response.NotFound, e)
	}

	resp := response.ConversionResponse{
		Amount:          amount,
		ConvertedAmount: amount * result.BuyRate,
		BaseCurrency:    baseCurrency,
		TargetCurrency:  targetCurrency,
		InitiatedOn:     int64(time.Nanosecond),
		Rate:            result.BuyRate,
	}
	return common.GetSimpleResponse[response.ConversionResponse](&resp, response.Success, nil)
}

func (s *Fx_service) GetConvertedRateById(c *context.Context,
	id int) response.ResponseWithSimpleData[response.ConversionResponse] {
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)

	result, err := s.DbService.GetOneById(id)
	span.End()

	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		common.Logger.Errorf("Error in retriving and converting forex rate. Exception:%v", err)
		return common.GetSimpleResponse[response.ConversionResponse](nil, response.NotFound, e)
	}

	resp := response.ConversionResponse{
		BaseCurrency:   result.BaseCurrency,
		TargetCurrency: result.TargetCurrency,
		InitiatedOn:    int64(time.Nanosecond),
		Rate:           result.BuyRate,
	}
	return common.GetSimpleResponse[response.ConversionResponse](&resp, response.Success, nil)
}

func (s *Fx_service) UpdateForexRate(c *context.Context,
	tenantId int, bankId int, baseCurrency string, targetCurrency string, tier string) response.ResponseWithSimpleData[response.ConversionResponse] {
	updateRequest := request.FxDataRequest{
		TenantId:       tenantId,
		BankId:         bankId,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Tier:           tier,
	}
	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)

	_, err := s.DbService.UpdateOne(updateRequest, updateRequest)
	span.End()

	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		common.Logger.Errorf("Error in retriving and converting forex rate. Exception:%v", err)
		return common.GetSimpleResponse[response.ConversionResponse](nil, response.NotFound, e)
	}

	//resp := response.ForexDataResponse{}
	return common.GetSimpleResponse[response.ConversionResponse](nil, response.Success, nil)
}

func (s *Fx_service) UpdateForexById(c *context.Context,
	id int) response.ResponseWithSimpleData[response.ConversionResponse] {

	tracer := otel.Tracer(os.Getenv(tracerName))
	_, span := tracer.Start(*c, dbSpanName)
	_, err := s.DbService.UpdateOneById(id)
	span.End()

	if err != nil {
		e := &[]response.Error{
			{Code: "DATA_NOT_FOUND", Message: "No record found", Details: "No record found"},
		}
		common.Logger.Errorf("Error in retriving and converting forex rate. Exception:%v", err)
		return common.GetSimpleResponse[response.ConversionResponse](nil, response.NotFound, e)
	}

	//resp := response.ForexDataResponse{}
	return common.GetSimpleResponse[response.ConversionResponse](nil, response.Success, nil)
}

func randomInt() int {
	// Define the range for the random integer (1 million to 2 million)
	minVal := 1000000
	maxVal := 2000000

	// Generate a random integer within the specified range
	return rand.Intn(maxVal-minVal+1) + minVal
}
