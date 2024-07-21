package metadata_api

import gosalesforce "github.com/ziemniakoss/go-salesforce"

type DeploymentScheduler interface {
	DeployFiles(
		sfProject gosalesforce.SfProject,
		connection gosalesforce.Connection,
		filePaths []string,
		options DeploymentOptions,
	) (string, error)

	DeployPackage(
		sfProject gosalesforce.SfProject,
		connection gosalesforce.Connection,
		sfPackage SfPackage,
		options DeploymentOptions,
	) (string, error)
}

type DeploymentOptions struct {
	CheckOnlyDeployment bool
	// Could be used for providing SF replacements
	AdditionalEnvVariables map[string]string
}

type SfPackage struct {
	Members []SfPackageMembers
}

type SfPackageMembers struct {
	MetadataType string
	ApiNames     []string
}
