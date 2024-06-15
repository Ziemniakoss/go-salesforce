package gosalesforce

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path"
	"strings"
)

// Manager that uses only file system to retrieve data about
// Salesforce connections.
// A lot faster but might break if Salesforce changes way in which
// auth data is stored
type FileSystemSfConnectionsManager struct {
}

// GetAllConnections implements SfConnectionsManager.
func (f FileSystemSfConnectionsManager) GetAllConnections() ([]Connection, error) {
	currentUser, error := user.Current()
	if error != nil {
		return []Connection{}, error
	}
	folderWithAuth := path.Join(currentUser.HomeDir, ".sfdx")
	// We don't care about fs errors
	entries, _ := os.ReadDir(folderWithAuth)

	usernameToConnection := make(map[string]Connection)
	aliases := make(map[string]string)
	for _, entry := range entries {
		if entry.Name() == "key.json" || entry.Name() == "stash.json" || strings.HasSuffix(entry.Name(), ".sandbox.json") {
			// we are ignoring this as no one knows what this is
		} else if entry.Name() == "alias.json" {
			fileContent, fsError := os.ReadFile(path.Join(folderWithAuth, entry.Name()))
			if fsError != nil {
				continue
			}
			var aliasesStore aliasesStore
			jsonError := json.Unmarshal(fileContent, &aliasesStore)
			if jsonError != nil {
				continue
			}
			for orgAlias, username := range aliasesStore.Orgs {
				aliases[username] = orgAlias
			}
		} else {
			var con Connection
			fileContent, fsError := os.ReadFile(path.Join(folderWithAuth, entry.Name()))
			if fsError != nil {
				continue
			}
			jsonError := json.Unmarshal(fileContent, &con)
			if jsonError != nil {
				continue
			}
			usernameToConnection[con.Username] = con
			// Real parsing
		}
	}
	connections := make([]Connection, 0, len(usernameToConnection))
	for _, connection := range usernameToConnection {
		alias, hasAlias := aliases[connection.Username]
		if hasAlias {
			connection.Alias = alias
		}
		connections = append(connections, connection)
	}
	return connections, nil
}

// GetConnection implements SfConnectionsManager.
func (f FileSystemSfConnectionsManager) GetConnection(usernameOrAlias string) (Connection, error) {
	// TODO possible fix optimization here because we don't need all connections here
	allConnections, error := f.GetAllConnections()
	if error != nil {
		return Connection{}, error
	}
	for _, connection := range allConnections {
		if connection.Username == usernameOrAlias || connection.Alias == usernameOrAlias {
			return connection, nil
		}
	}
	return Connection{}, errors.New("NO_CONNECTION")
}

// GetDefaultConnection implements SfConnectionsManager.
func (f FileSystemSfConnectionsManager) GetDefaultConnection(filePath string) (Connection, error) {
	defaultUsername, error := f.GetDefaultUsername(filePath)
	if error != nil {
		return Connection{}, error
	}
	return f.GetConnection(defaultUsername)
}

// GetDefaultUsername implements SfConnectionsManager.
func (f FileSystemSfConnectionsManager) GetDefaultUsername(filePath string) (string, error) {
	currentPath, error := os.Stat(filePath)
	if error != nil {
		return "", errors.New("FILE_DOES_NOT_EXIST")
	}

	currentDir := filePath
	if !currentPath.IsDir() {
		currentDir = path.Dir(currentDir)
	}

	for currentDir != "/" {
		if sfUsername, _ := f.readUsernameFromSfConfig(currentDir); sfUsername != "" {
			return sfUsername, nil
		}
		if sfdxUsername, _ := f.readUsernameFromSfdxConfig(currentDir); sfdxUsername != "" {
			return sfdxUsername, nil
		}

		currentDir = path.Dir(currentDir)
	}
	return "", errors.New("NO_DEFAULT_ORG")
}

func (f FileSystemSfConnectionsManager) readUsernameFromSfConfig(directory string) (string, error) {
	sfConfigFile := path.Join(directory, ".sf", "config.json")
	configContent, error := os.ReadFile(sfConfigFile)
	if error != nil {
		return "", error
	}

	var parsedConfig sfConfig
	parsingError := json.Unmarshal(configContent, &parsedConfig)
	if parsingError != nil {
		return "", error
	}
	if parsedConfig.TargetOrg == "" {
		return "", errors.New("NO_DEFAULT_ORG_SET_UP")
	}
	return parsedConfig.TargetOrg, nil
}

type sfConfig struct {
	TargetOrg string `json:"target-org"`
}

type sfdxConfig struct {
	Username string `json:"defaultusername"`
}

func (f FileSystemSfConnectionsManager) readUsernameFromSfdxConfig(directory string) (string, error) {
	sfdxConfigFile := path.Join(directory, ".sfdx", "sfdx-config.json")
	configContent, error := os.ReadFile(sfdxConfigFile)
	if error != nil {
		return "", error
	}

	var parsedConfig sfdxConfig
	parsingError := json.Unmarshal(configContent, &parsedConfig)
	if parsingError != nil {
		return "", error
	}
	if parsedConfig.Username == "" {
		return "", errors.New("NO_DEFAULT_ORG_SET_UP")
	}
	return parsedConfig.Username, nil
}
