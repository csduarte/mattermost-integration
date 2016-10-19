package platform

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

// MapToJSON will encode map to json
func MapToJSON(objmap map[string]string) string {
	b, err := json.Marshal(objmap)
	if err != nil {
		return ""
	}
	return string(b)
}

// MapFromJSON will decode the key/value pair map
func MapFromJSON(data io.Reader) map[string]string {
	decoder := json.NewDecoder(data)

	var objmap map[string]string
	if err := decoder.Decode(&objmap); err != nil {
		return make(map[string]string)
	}
	return objmap
}

// NewClientError sets ClientError
func NewClientError(where string, id string, params map[string]interface{}, details string) *ClientError {
	ap := &ClientError{}
	ap.ID = id
	ap.Params = params
	ap.Message = id
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = 500
	ap.IsOAuth = false
	return ap
}

// IsValidHTTPUrl checks for http or https url
func IsValidHTTPUrl(rawURL string) bool {
	if strings.Index(rawURL, "http://") != 0 && strings.Index(rawURL, "https://") != 0 {
		return false
	}

	if _, err := url.ParseRequestURI(rawURL); err != nil {
		return false
	}

	return true
}
