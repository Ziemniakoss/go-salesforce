package gosalesforce

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// SfdxAuthResult represents the "result" object in the sfdx JSON output.
type sfdxAuthResult struct {
	SfdxAuthUrl string `json:"sfdxAuthUrl"`
}

// SfdxAuthResponse represents the full JSON response from the sfdx CLI.
type sfdxAuthResponse struct {
	Status int            `json:"status"`
	Result sfdxAuthResult `json:"result"`
}

func (con SfConnectionWithToken) GetSfAuthUrl() (string, error) {
	return getSfAuthUrl(con.Username)
}

func (con SfConnectionLight) GetSfAuthUrl() (string, error) {
	return getSfAuthUrl(con.Username)
}

func getSfAuthUrl(targetOrg string) (string, error) {
	cmd := exec.Command(
		"sf",
		"org",
		"auth",
		"show-sfdx-auth-url",
		"--target-org",
		targetOrg,
		"--json",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		// Even on failure, sfdx often still prints valid JSON with error details,
		// so include the raw output in the error for debugging.
		return "", fmt.Errorf("sfdx command failed: %w\noutput: %s", err, string(out))
	}

	var resp sfdxAuthResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return "", fmt.Errorf("failed to parse sfdx output: %w\noutput: %s", err, string(out))
	}

	if resp.Status != 0 {
		return "", fmt.Errorf("sfdx returned non-zero status %d\noutput: %s", resp.Status, string(out))
	}

	if resp.Result.SfdxAuthUrl == "" {
		return "", fmt.Errorf("sfdxAuthUrl not found in response\noutput: %s", string(out))
	}

	return resp.Result.SfdxAuthUrl, nil
}
