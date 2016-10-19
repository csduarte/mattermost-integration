package platform

import (
	"encoding/json"
	"io"
)

// ClientError for requests
type ClientError struct {
	ID            string                 `json:"id"`
	Message       string                 `json:"message"`               // Message to be display to the end user without debugging information
	DetailedError string                 `json:"detailed_error"`        // Internal error string to help the developer
	RequestID     string                 `json:"request_id,omitempty"`  // The RequestId that's also set in the header
	StatusCode    int                    `json:"status_code,omitempty"` // The http status code
	Where         string                 `json:"-"`                     // The function where it happened in the form of Struct.Func
	IsOAuth       bool                   `json:"is_oauth,omitempty"`    // Whether the error is OAuth specific
	Params        map[string]interface{} `json:"params"`
}

// ClientErrorFromJSON translates json error to client error
func ClientErrorFromJSON(data io.Reader) *ClientError {
	decoder := json.NewDecoder(data)
	var er ClientError
	err := decoder.Decode(&er)
	if err == nil {
		return &er
	}
	return NewClientError("ClientErrorFromJson", "model.utils.decode_json.app_error", nil, err.Error())
}
