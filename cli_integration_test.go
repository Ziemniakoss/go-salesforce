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
