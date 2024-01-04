package test

import (
	"github.com/PeerIslands/aci-fx-go/model/entity"
	"github.com/stretchr/testify/mock"
)

type MockDbService struct {
	mock.Mock
}

func (m *MockDbService) Init(credentials ...string) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDbService) Get(filter any) ([]entity.ForexData, error) {
	args := m.Called(filter)
	return args.Get(0).([]entity.ForexData), nil
}

func (m *MockDbService) CreateOne(document entity.ForexData) (entity.ForexData, error) {
	args := m.Called(document)
	return args.Get(0).(entity.ForexData), nil
}

func (m *MockDbService) UpdateOne(document any, filter any) (any, error) {
	args := m.Called(filter, document)
	return args.Get(0).(entity.ForexData), nil
}

func (m *MockDbService) UpdateOneById(id any) (any, error) {
	args := m.Called(id)
	return args.Get(0).(entity.ForexData), nil
}
func (m *MockDbService) DeleteOne(filter any) (int64, error) {
	args := m.Called(filter)
	return args.Get(0).(int64), nil
}

func (m *MockDbService) GetOne(filter any) (entity.ForexData, error) {
	args := m.Called(filter)
	return args.Get(0).(entity.ForexData), nil
}

//
//func TestGetConvertedRate(t *testing.T) {
//	mockRepo := new(MockDbService)
//	mockRepo.On("GetOne", mock.Anything).Return(entity.ForexData{
//		ID:                           primitive.ObjectID{},
//		Tier:                         "1",
//		DirectIndirectFlag:           "Y",
//		Multiplier:                   1,
//		BuyRate:                      2,
//		SellRate:                     3,
//		TolerancePercentage:          0,
//		EffectiveDate:                nil,
//		ExpirationDate:               nil,
//		ContractRequirementThreshold: "",
//		TenantID:                     1,
//		BankID:                       1,
//		BaseCurrency:                 "USD",
//		TargetCurrency:               "EUR",
//		CreatedDate:                  time.Time{},
//		DocVersion:                   0,
//		UpdatedDate:                  time.Time{},
//	}, nil)
//	service := bal.Fx_service{DbService: mockRepo}
//
//	res := service.GetConvertedRate(1, 1, 1000, "USD", "EUR", "1")
//
//	assert.Equal(t, 2000.00, res.Data.ConvertedAmount)
//}
