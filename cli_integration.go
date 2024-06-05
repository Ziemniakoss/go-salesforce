package gosalesforce

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

/*
Retrieves org info using sf executable.
*/
func GetAllConnectionsUsingSf() ([]Connection, error) {
	return []Connection{}, nil
}

/*
Retrieves all connections stored in file system.
*/
func GetAllConnectionFromHomeDir() ([]Connection, error) {
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
		if entry.Name() == "key.json" || entry.Name() == "stash.json" {
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

type aliasesStore struct {
	Orgs map[string]string `json:"orgs"`
}
