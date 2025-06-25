package endpoints

import (
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"strings"

	"ama/api/constants"
	"ama/api/interfaces"
	"ama/api/logging"
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

// Generates a random string based on the base64 encoded output of a random
// byte slice defined by length. length does not determine the final string
// length. The final string will be at the longest the same size as length.
func generateRandomString(length int) string {
	if length <= 0 {
		length = 1
	}
	randomBytes := make([]byte, length)
	l, err := rand.Read(randomBytes)
	if err != nil || l != length {
		logger := logging.GetLogger()
		logger.Error("Error generating random string", "length", length, "error", err)
		return ""
	}
	// Firestore document names cannot contain forward slashes, so we replace them with
	// an empty string to prevent errors.
	randomizedValue := strings.ReplaceAll(base64.RawStdEncoding.EncodeToString(randomBytes), "/", "")
	if len(randomizedValue) < length {
		return randomizedValue
	}
	return randomizedValue[:length]
}


// Parses the API parameters for reading questions with defaults.
func GetReadQuestionsParamsWithDefaults(c interfaces.APIContext) (limit int, finalId string, tags []string) {
	logger := logging.GetLogger()
	limit = GetQueryParamToInt(c, constants.LimitParam, constants.DefaultLimit)
	finalId = c.DefaultQuery(constants.FinalIdParam, generateRandomString(8))
	tags, hasTags := c.GetQueryArray(constants.TagParam)
	if !hasTags {
		tags = []string{}
	}
	logger.Debug(
		"Read question params",
		constants.LimitParam, limit,
		constants.FinalIdParam, finalId,
		constants.TagParam, tags,
	)
	return limit, finalId, tags
}
