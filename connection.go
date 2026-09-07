package gosalesforce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

type sfAuthTokenCommandResult struct {
	AccessToken string `json:"accessToken"`
}

func (con SfConnectionLight) GetConnectionWithAccessToken() (SfConnectionWithToken, error) {
	cmd := exec.Command("sf", "org", "auth", "show-access-token", "--json", "--target-org", con.Username)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return SfConnectionWithToken{}, fmt.Errorf("command failed: %w, output: %s", err, string(output))
	}

	var parsed sfCommandResult[sfAuthTokenCommandResult]
	if err := json.Unmarshal(output, &parsed); err != nil {
		return SfConnectionWithToken{}, fmt.Errorf("failed to parse JSON: %w, raw output: %s", err, string(output))
	}

	if parsed.Status != 0 {
		return SfConnectionWithToken{}, fmt.Errorf("sf command returned non-zero status: %d", parsed.Status)
	}

	return SfConnectionWithToken{
		InstanceURL: con.InstanceURL,
		APIVersion:  con.APIVersion,
		HTTPClient:  http.DefaultClient,
		Username:    con.Username,
		Alias:       con.Alias,
		Token:       parsed.Result.AccessToken,
	}, nil
}

type sFOrgListResult struct {
	Other          []SfConnectionLight `json:"other"`
	SandboxOrgs    []SfConnectionLight `json:"sandboxes"`
	ScratchOrgs    []SfConnectionLight `json:"scratchOrgs"`
	NonScratchOrgs []SfConnectionLight `json:"nonScratchOrgs"`
}

func ListOrgs() ([]SfConnectionLight, error) {
	cmd := exec.Command("sf", "org", "list", "--json")

	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var response sfCommandResult[sFOrgListResult]
	if err := json.Unmarshal(out, &response); err != nil {
		return nil, fmt.Errorf("failed to parse sf org list output: %w\noutput: %s", err, string(out))
	}

	orgs := []SfConnectionLight{}
	orgs = append(orgs, response.Result.NonScratchOrgs...)
	orgs = append(orgs, response.Result.SandboxOrgs...)
	orgs = append(orgs, response.Result.ScratchOrgs...)
	orgs = append(orgs, response.Result.Other...)
	return orgs, nil
}
