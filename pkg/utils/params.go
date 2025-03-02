package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetQueryParamInt(c echo.Context, key string) (*int, error) {
	intStr := c.QueryParam(key)

	if intStr == "" {
		return nil, nil
	}

	num, err := strconv.Atoi(intStr)
	if err != nil {
		return nil, err
	}

	return &num, nil
}

func GetParamInt(c echo.Context, key string) (*int, error) {
	intStr := c.Param(key)

	if intStr == "" {
		return nil, nil
	}

	num, err := strconv.Atoi(intStr)
	if err != nil {
		return nil, err
	}

	return &num, nil
}
