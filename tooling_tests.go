package gosalesforce

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

/*
Sends request to schedule async apex test run using `/services/data/ver/tooling/runTestsAsynchronous` endpoint.

Returns Id of apex test run job
*/
func RunTestsAsync(connection Connection, request AsyncTestRunRequest) (string, error) {
	url := connection.InstanceUrl + "/services/data/v" + connection.ApiVersion + "/tooling/runTestsAsynchronous"
	httpRequest, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", err
	}
	httpRequest.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + connection.AccessToken},
	}
	client := http.Client{}
	response, httpError := client.Do(httpRequest)
	if httpError != nil {
		return "", httpError
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if response.StatusCode == 200 {
		var testRunId string
		parsingError := json.Unmarshal(body, &testRunId)
		if parsingError != nil {
			return "", parsingError
		}
		return testRunId, nil
	}
	return "", errors.New(string(body))
}

// TODO test result

type AsyncTestRunRequest struct {
	/*
		Tests to run as part of test run.
		Single class can only be referenced in one instance of AsyncTestRunSingleItem.
		Duplicates can cause SF error
	*/
	Tests []AsyncTestRunSingleItem `json:"tests"`
	/*
		When true, code coverage is not calculated, which in
		big orgs can speed up test execution.
	*/
	SkipCodeCoverage bool `json:"skipCodeCoverage"`
}

type AsyncTestRunSingleItem struct {
	ClassName   string   `json:"className"`
	TestMethods []string `json:"testMethods"`
}
