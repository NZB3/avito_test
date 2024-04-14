package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"testing"
)

func TestController_ParseParam(t *testing.T) {
	params := gin.Params{
		{Key: "tag_id", Value: "12"},
	}

	tagID, err := ParseQueryParam(params, "tag_id", true, func(param string) (int, error) {
		return strconv.Atoi(param)
	})

	if err != nil {
		t.Error(err)
	}

	notRequiredBool, err := ParseQueryParam(params, "use", false, ConvToBool)
	if err != nil {
		t.Error(err)
	}

	t.Log(tagID)
	t.Log(notRequiredBool)
}
