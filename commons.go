package gosalesforce

import "net/http"

const DefaultAPIVersion = "67.0"

type sfCommandResult[T any] struct {
	Code   string `json:"code"`
	Status int    `json:"status"`
	Result T
}
type SfError struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type SfConnectionLight struct {
	InstanceURL string `json:"instanceUrl,omitempty"`
	APIVersion  string `json:"apiVersion,omitempty"`
	Username    string `json:"username,omitempty"`
	Alias       string `json:"alias,omitempty"`
}

// Authorized connection
type SfConnectionWithToken struct {
	InstanceURL string `json:"instanceUrl,omitempty"`
	APIVersion  string `json:"apiVersion,omitempty"`
	HTTPClient  *http.Client
	Username    string `json:"username,omitempty"`
	Alias       string `json:"alias,omitempty"`
	Token       string `json:"token,omitempty"`
}
