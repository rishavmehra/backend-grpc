package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/rishavmehra/backend-grpc/util"
)

var validCurrency validator.Func = func(fieldlevel validator.FieldLevel) bool {
	if currency, ok := fieldlevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}
	return false
}
