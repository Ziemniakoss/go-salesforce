package gosalesforce

import (
	"encoding/json"
	"errors"
	"os"
	"path"
)

type SfProject struct {
	RootFolder  string                      `json:"root"`
	Name        string                      `json:"name"`
	Directories []SfProjectPackageDirectory `json:"packageDirectories"`
	ApiVersion  string                      `json:"sourceApiVersion"`
}

type SfProjectPackageDirectory struct {
	IsDefault bool   `json:"default"`
	Path      string `json:"path"`
}

func GetSfProject(filePath string) (SfProject, error) {
	currentPath, error := os.Stat(filePath)
	if error != nil {
		return SfProject{}, errors.New("FILE_DOES_NOT_EXIST")
	}

	currentDir := filePath
	if !currentPath.IsDir() {
		currentDir = path.Dir(currentDir)
	}
	for {
		projectDefinitionFile := path.Join(currentDir, "sfdx-project.json")
		config_content, error := os.ReadFile(projectDefinitionFile)
		if error == nil {
			var project SfProject
			if parsingError := json.Unmarshal(config_content, &project); parsingError != nil {
				return SfProject{}, parsingError
			}
			project.RootFolder = currentDir
			return project, nil
		}
		nextDir := path.Dir(currentDir)
		if nextDir == currentDir {
			return SfProject{}, errors.New("NOT_IN_SF_PROJECT")
		}
		currentDir = nextDir
	}

}
