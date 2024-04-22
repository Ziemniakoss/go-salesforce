package gosalesforce

// TODO test for this
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetDeploymentInfo(connection Connection, deploymentId string) (DeploymentInfo, error) {
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
	fmt.Println(string(body))
	parsingError := json.Unmarshal(body, &deploymentInfo)
	if parsingError != nil {
		return DeploymentInfo{}, parsingError
	}
	return deploymentInfo, nil

}

type DeploymentInfo struct {
	Id     string               `json:"id"`
	Result DeploymentResultInfo `json:"deployResult"`
}

type DeploymentResultInfo struct {
	Details   DeploymentInfoDetails `json:"details"`
	Status    string                `json:"status"`
	Success   bool                  `json:"success"`
	CheckOnly bool                  `json:"checkOnly"`
}

type DeploymentInfoDetails struct {
	Successes []DeploymentInfoItem `json:"componentSuccesses"`
	Failures  []DeploymentInfoItem `json:"componentFailures"`
}

type DeploymentInfoItem struct {
	Type              string `json:"componentType"`
	FullName          string `json:"fullName"`
	Id                string `json:"id"`
	IsCreated         bool   `json:"created"`
	IsChanged         bool   `json:"changed"`
	IsDeleted         bool   `json:"deleted"`
	IsSuccess         bool   `json:"success"`
	ErrorLineNumber   int    `json:"lineNumber"`
	ErrorColumnNumber int    `json:"columnNumber"`
}

type DeploymentInfoFailureItem struct {
}
