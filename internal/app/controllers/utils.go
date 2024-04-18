package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

func ConvToInt(param string) (int, error) {
	return strconv.Atoi(param)
}

func ConvToBool(param string) (bool, error) {
	return strconv.ParseBool(param)
}

func ParseQueryParam[T any](queryContext *gin.Context, name string, required bool, defaultVal T, convFunc func(param string) (T, error)) (convertedParam T, err error) {
	param := queryContext.Query(name)
	if param == "" {
		if required {
			return convertedParam, errors.New(name + " is required")
		}

		convertedParam = defaultVal
		return convertedParam, nil
	}

	convertedParam, err = convFunc(param)
	if err != nil {
		return convertedParam, errors.New(name + " is invalid")
	}

	return convertedParam, nil
}

func CheckAdminStatus(ctx *gin.Context) (isAdmin bool, err error) {
	admin, ok := ctx.Get("admin")
	if !ok {
		return false, errors.New("'admin' is not specified in context")
	}

	return admin.(bool), nil
}
