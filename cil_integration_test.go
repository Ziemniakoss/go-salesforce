package gosalesforce_test

import (
	"os"
	"os/user"
	"testing"

	gosalesforce "github.com/ziemniakoss/go-salesforce"
)

func TestUserInfoFetching(t *testing.T) {
	currentUser, error := user.Current()
	if error != nil {
		t.Fatalf("Could not retrieve user info on current platform")
	}

	realHome, _ := os.LookupEnv("HOME")
	if currentUser.HomeDir != realHome {
		t.Fatalf("Wrong house, got %s, expected %s", currentUser.HomeDir, realHome)
	}

}

func TestGetAllConnectionFromHomeDir(t *testing.T) {
	// TODO something usefull
	gosalesforce.GetAllConnectionFromHomeDir()
}

/*
Test that guarantees that results returned by all connection-fetching
methods return same list of connections
*/
func TestGettingConnectionsIntegrity(t *testing.T) {
	fsConnections, fsErr := gosalesforce.GetAllConnectionFromHomeDir()
	if fsErr != nil {
		t.Fatalf("Could not fetch connections using FS method: %s", fsErr.Error())
	}
	sfConnections, sfError := gosalesforce.GetAllConnectionsUsingSf()
	if sfError != nil {
		t.Fatalf("Could not retrieve connections using SF method: %s", sfError.Error())
	}
	if len(fsConnections) != len(sfConnections) {
		t.Fatalf("List of connections retrieved using SF (%d) and FS (%d) method is different", len(sfConnections), len(fsConnections))

	}

}
