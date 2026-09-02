package gosalesforce

import "net/http"

const DefaultAPIVersion = "67.0"

type SfError struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type SfConnection struct {
	InstanceURL string
	AccessToken string
	APIVersion  string
	HTTPClient  *http.Client
}

/*
Connection to Salesforce org, legacy
*/
type Connection struct {
	AccessToken  string `json:"accessToken"`
	InstanceUrl  string `json:"instanceUrl"`
	Alias        string `json:"alias"`
	Username     string `json:"username"`
	RefreshToken string `json:"refreshToken"`
	IsDevHub     bool   `json:"isDevHub"`
	OrgId        string `json:"orgId"`
	ApiVersion   string `json:"instanceApiVersion"`
}
