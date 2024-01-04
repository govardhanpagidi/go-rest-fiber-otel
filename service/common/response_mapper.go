package common

import "github.com/PeerIslands/aci-fx-go/model/dto/response"

func GetSimpleResponse[T any](data *T, statusCode response.StatusCode, errors *[]response.Error) response.ResponseWithSimpleData[T] {
	return response.ResponseWithSimpleData[T]{
		Data:   data,
		Status: statusCode,
		Errors: errors,
	}
}

func GetArrayResponse[T any](data *[]T, statusCode response.StatusCode, errors *[]response.Error) response.ResponseWithArrayData[T] {
	return response.ResponseWithArrayData[T]{
		Data:   data,
		Status: statusCode,
		Errors: errors,
	}
}
