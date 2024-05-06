package gosalesforce

/*
Connection to Salesforce org
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
