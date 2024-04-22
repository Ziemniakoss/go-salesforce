package gosalesforce

/*
Connection to Salesforce org
*/
type Connection struct {
	AccessToken string `json:accessToken`
	InstanceUrl string `json:instaceUrl`
	Alias       string `json:alias`
	Username    string `json:username`
}
