package metadata_api

import (
	"encoding/json"
	"io"
	"net/http"

	gosalesforce "github.com/Ziemniakoss/go-salesforce"
)

/*
Integrate with SF using direct API access or file system calls (for auth data).
Should be much faster than CLI integration BUT will be less complete
*/
type ApiMetadataClient struct {
}

/*
Fetch information about specific deployment

## Arguments

- connection: defines which org should be used
- deploymentId: Id of deployment
*/
func (ApiMetadataClient) GetDeploymentInfo(connection gosalesforce.Connection, deploymentId string) (DeploymentInfo, error) {
	url := connection.InstanceUrl + "/services/data/v60.0/metadata/deployRequest/" + deploymentId + "?includeDetails=true"
	request, error := http.NewRequest("GET", url, nil)
	if error != nil {
		return DeploymentInfo{}, nil
	}
	request.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + connection.AccessToken},
	}
	client := http.Client{}
	response, httpError := client.Do(request)
	if httpError != nil {
		return DeploymentInfo{}, httpError
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return DeploymentInfo{}, err
	}
	var deploymentInfo DeploymentInfo
	parsingError := json.Unmarshal(body, &deploymentInfo)
	if parsingError != nil {
		return DeploymentInfo{}, parsingError
	}
	return deploymentInfo, nil
}
