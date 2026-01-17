package metadata_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"

	gosalesforce "github.com/Ziemniakoss/go-salesforce"
)

func (client CmdMetadataClient) DeployFiles(
	sfProject gosalesforce.SfProject,
	connection gosalesforce.Connection,
	filePaths []string,
	options DeploymentOptions,
) (string, error) {
	additionalParameters := []string{}
	for _, path := range filePaths {
		additionalParameters = append(additionalParameters, "-p", path)
	}
	return client.runDeployment(
		sfProject,
		connection,
		options,
		additionalParameters,
	)
}

type scheduledDeploymentCmdResult struct {
	Result  scheduledDeploymentCmdResultResult
	Status  int    `json:"status"`
	Message string `json:"message"`
	Name    string `json:"name"`
}

type scheduledDeploymentCmdResultResult struct {
	Id string `json:"id"`
}

func (client CmdMetadataClient) DeployPackage(
	sfProject gosalesforce.SfProject,
	connection gosalesforce.Connection,
	sfPackage SfPackage,
	options DeploymentOptions,
) (string, error) {
	additionalParameters := []string{}
	for _, members := range sfPackage.Members {
		for _, apiName := range members.ApiNames {
			additionalParameters = append(additionalParameters, "-m", members.MetadataType+":"+apiName)
		}
	}

	return client.runDeployment(
		sfProject,
		connection,
		options,
		additionalParameters,
	)
}

func (client CmdMetadataClient) runDeployment(
	sfProject gosalesforce.SfProject,
	connection gosalesforce.Connection,
	options DeploymentOptions,
	additionalArguments []string,
) (string, error) {
	command := []string{
		"force",
		"source",
		"deploy",
		"--json",
		"-w",
		"0",
		"-u",
		connection.Username,
	}
	command = append(command, additionalArguments...)
	if options.CheckOnlyDeployment {
		command = append(command, "-c")
	}

	cmd := exec.Command(client.BinaryPath, command...)
	cmd.Dir = sfProject.RootFolder
	// TODO Environmental variables
	// cmd.Env = os.Environ()
	// cmd.Env = append(cmd.Env, "MY_VAR=some_value")
	var buffer bytes.Buffer
	cmd.Stdout = &buffer

	cmd.Run()
	out := buffer.String()
	var cmdOut scheduledDeploymentCmdResult
	err := json.Unmarshal([]byte(out), &cmdOut)
	if err != nil {
		return "", err
	}
	if cmdOut.Message != "" {
		return "", errors.New(cmdOut.Message)
	}
	return cmdOut.Result.Id, nil

}
