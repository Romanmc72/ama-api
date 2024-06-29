package endpoints

import (
	"strconv"
	"strings"

	"ama/api/interfaces"
	"ama/api/logging"
	"ama/api/paths"
)

// This only kinda works, see this for examples:
// https://gin-gonic.com/docs/examples/querystring-param/
// Retrieves a url query parameter and parses it to an int if it is set.
// If it is not set it will return the input default value.
func GetQueryParamToInt(c interfaces.APIContext, paramName string, defaultValue int) int {
	rawValue := c.DefaultQuery(paramName, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(rawValue)
	if err != nil {
		logger := logging.GetLogger()
		logger.Debug(
			"Error converting value for query param, using default",
			"param", paramName,
			"default", defaultValue,
			"error", err,
		)
		return defaultValue
	}
	return value
}

// Retrieves a url query parameter and parses it to an int64 if it is set.
// If it is not set it will return the input default value.
func GetQueryParamToInt64(c interfaces.APIContext, paramName string, defaultValue int64) int64 {
	rawValue := c.DefaultQuery(paramName, strconv.FormatInt(defaultValue, 10))
	value, err := strconv.ParseInt(rawValue, 10, 64)
	if err != nil {
		logger := logging.GetLogger()
		logger.Debug(
			"unable to parse query param",
			"param", paramName,
			"value", rawValue,
			"error", err,
		)
		return defaultValue
	}
	return value
}

// If a query parameter is supposed to be an array of values delimited by
// some separator, then this will split them out to that array.
func GetQueryParamToStringArray(c interfaces.APIContext, paramName string, defaultValue string) []string {
	rawValue := c.DefaultQuery(paramName, defaultValue)
	values := strings.Split(rawValue, paths.ArraySeparator)
	if strings.TrimSpace(values[0]) == "" {
		return []string{}
	}
	return values
}
