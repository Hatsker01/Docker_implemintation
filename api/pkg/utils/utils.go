package utils

import (
	"strconv"
	"strings"
)

type QueryParams struct {
	Filters  map[string]string
	Page     int64
	Limit    int64
	Ordering []string
	Search   string
}

type QueryParamLog struct {
	Filters  map[string]string
	Email    string
	Password string
}

func ParseQueryParams(queryParams map[string][]string) (*QueryParams, []string) {
	params := QueryParams{
		Filters:  make(map[string]string),
		Page:     1,
		Limit:    10,
		Ordering: []string{},
		Search:   "",
	}
	var errStr []string
	var err error

	for key, value := range queryParams {
		if key == "page" {
			params.Page, err = strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				errStr = append(errStr, "Invalid `page` param")
			}
			continue
		}

		if key == "limit" {
			params.Limit, err = strconv.ParseInt(value[0], 10, 64)
			if err != nil {
				errStr = append(errStr, "Invalid `limit` param")
			}
			continue
		}

		if key == "search" {
			params.Search = value[0]
			continue
		}

		if key == "ordering" {
			params.Ordering = strings.Split(value[0], ",")
			continue
		}
		params.Filters[key] = value[0]
	}

	return &params, errStr
}

func ParseQueryParamsLog(queryParams map[string][]string) (*QueryParamLog, []string) {
	params := QueryParamLog{
		Filters:  make(map[string]string),
		Email:    "",
		Password: "",
	}
	var errStr []string

	for key, value := range queryParams {
		if key == "email" {
			params.Email = value[0]

		}

		if key == "password" {
			params.Password = value[0]

		}

	}

	return &params, errStr
}
