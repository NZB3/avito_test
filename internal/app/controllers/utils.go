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

func ParseQueryParam[T any](params gin.Params, name string, required bool, convFunc func(param string) (T, error)) (convertedParam T, err error) {
	param := params.ByName(name)
	if param == "" {
		if required {
			return convertedParam, errors.New(name + " is required")
		}

		return convertedParam, nil
	}

	convertedParam, err = convFunc(param)
	if err != nil {
		return convertedParam, errors.New(name + " is invalid")
	}

	return convertedParam, nil
}
