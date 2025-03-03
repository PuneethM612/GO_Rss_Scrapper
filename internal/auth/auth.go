package auth

import (
	"net/http"
	"strings"
)

// GET API KEY
// extracts api key from header
// Example:
// Authorization: ApiKey {insert api key here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", http.ErrNoCookie
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", http.ErrNoCookie
	}
	if vals[0] != "ApiKey" {
		return "", http.ErrNoCookie
	}
	return vals[1], nil

}
