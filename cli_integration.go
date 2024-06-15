package gosalesforce

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
	return FileSystemSfConnectionsManager{}.GetAllConnections()
}

type aliasesStore struct {
	Orgs map[string]string `json:"orgs"`
}

type SfConnectionsManager interface {
	GetAllConnections() ([]Connection, error)
	GetConnection(usernameOrAlias string) (Connection, error)
	GetDefaultUsername(filePath string) (string, error)
	GetDefaultConnection(filePath string) (Connection, error)
}
