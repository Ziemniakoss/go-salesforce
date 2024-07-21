package metadata_api

import (
	gosalesforce "github.com/ziemniakoss/go-salesforce"
)

type DeploymentInfoProvider interface {
	/*
		Fetch information about specific deployment

		## Arguments

		- connection: defines which org should be used
		- deploymentId: Id of deployment
	*/
	GetDeploymentInfo(connection gosalesforce.Connection, deploymentId string) (DeploymentInfo, error)
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
