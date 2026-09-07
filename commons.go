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
	InstanceURL  string `json:"instanceUrl,omitempty"`
	APIVersion   string `json:"instanceApiVersion,omitempty"`
	Username     string `json:"username,omitempty"`
	Alias        string `json:"alias,omitempty"`
	OrgEdition   string `json:"orgEdition,omitempty"`
	IsDevOrg     bool   `json:"isDevOrg,omitempty"`
	IsSandbox    bool   `json:"isSandbox,omitempty"`
	IsScratchOrg bool   `json:"isScratchOrg,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
}

// Authorized connection
type SfConnectionWithToken struct {
	InstanceURL  string `json:"instanceUrl,omitempty"`
	APIVersion   string `json:"instanceApiVersion,omitempty"`
	HTTPClient   *http.Client
	Username     string `json:"username,omitempty"`
	Alias        string `json:"alias,omitempty"`
	Token        string `json:"token,omitempty"`
	OrgEdition   string `json:"orgEdition,omitempty"`
	IsDevOrg     bool   `json:"isDevOrg,omitempty"`
	IsSandbox    bool   `json:"isSandbox,omitempty"`
	IsScratchOrg bool   `json:"isScratchOrg,omitempty"`
	InstanceName string `json:"instanceName,omitempty"`
}
