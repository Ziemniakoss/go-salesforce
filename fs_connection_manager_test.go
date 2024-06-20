package gosalesforce_test

import (
	"math/rand"
	"os"
	"path"
	"strconv"
	"testing"

	gosalesforce "github.com/ziemniakoss/go-salesforce"
)

// Creates random sfdx project and returns path to it
func createTestProject(t *testing.T) string {
	projectDir := path.Join(os.TempDir(), strconv.Itoa(rand.Int()))
	sfdxConfigFolder := path.Join(projectDir, ".sfdx")
	if err := os.MkdirAll(sfdxConfigFolder, 0777); err != nil {
		t.Error("Could not create sfdx folder", err.Error())
		t.FailNow()
	}
	sfdxProjectConfigPath := path.Join(projectDir, "sfdx-project.json")
	os.WriteFile(sfdxProjectConfigPath, []byte("{}"), 0777)

	return projectDir
}

func createTestProjectWithFolder(t *testing.T) (string, string) {
	projectDir := createTestProject(t)
	testedFolder := path.Join(projectDir, "some", "custom", "folder")
	if err := os.MkdirAll(testedFolder, 0777); err != nil {
		t.Error("Could not create test folder", err.Error())
		t.FailNow()
	}
	return projectDir, testedFolder
}

func createTestProjectWithFile(t *testing.T) (string, string) {
	projectDir, testFolder := createTestProjectWithFolder(t)
	testFile := path.Join(testFolder, "someTestFile.json")
	if err := os.WriteFile(testFile, []byte("aaa"), 0777); err != nil {
		t.Error("Could not create test file")
		t.FailNow()
	}
	return projectDir, testFile
}

func TestGetDefaultUsername_sfdxConfig_folder(t *testing.T) {
	projectDir, testFolder := createTestProjectWithFolder(t)
	sfdxConfigFolder := path.Join(projectDir, ".sfdx")
	if err := os.MkdirAll(sfdxConfigFolder, 0777); err != nil {
		t.Error("Could not create sfdx folder", err.Error())
		t.FailNow()
	}

	sfdxConfigContent := "{\"defaultusername\": \"someOrg\"}"
	if err := os.WriteFile(path.Join(sfdxConfigFolder, "sfdx-config.json"), []byte(sfdxConfigContent), 0777); err != nil {
		t.Error("Could not create sfdx config file", err.Error())
		t.FailNow()
	}

	result, err := gosalesforce.FileSystemSfConnectionsManager{}.GetDefaultUsername(testFolder)
	if err != nil {
		t.Error("Could not read sfdx username", err.Error())
		t.FailNow()
	}
	if result != "someOrg" {
		t.Error("Wrong username returned, expected", "someOrg", "got", result)
	}
}

func TestGetDefaultUsername_sfCnfig_folder(t *testing.T) {
	projectDir, testFolder := createTestProjectWithFolder(t)
	sfConfigFolder := path.Join(projectDir, ".sf")
	if err := os.MkdirAll(sfConfigFolder, 0777); err != nil {
		t.Error("Could not create sfdx folder", err.Error())
		t.FailNow()
	}

	sfConfigContent := "{\"target-org\": \"testOrg\"}"
	if err := os.WriteFile(path.Join(sfConfigFolder, "config.json"), []byte(sfConfigContent), 0777); err != nil {
		t.Error("Could not create sfdx config file", err.Error())
		t.FailNow()
	}

	result, err := gosalesforce.FileSystemSfConnectionsManager{}.GetDefaultUsername(testFolder)
	if err != nil {
		t.Error("Could not read sf username", err.Error())
		t.FailNow()
	}
	if result != "testOrg" {
		t.Error("Wrong username returned, expected testOrg got", result)
	}
}

func TestGetDefaultUsername_sfdxConfig_file(t *testing.T) {
	projectDir, testFile := createTestProjectWithFile(t)
	sfdxConfigFolder := path.Join(projectDir, ".sfdx")
	if err := os.MkdirAll(sfdxConfigFolder, 0777); err != nil {
		t.Error("Could not create sfdx folder", err.Error())
		t.FailNow()
	}

	sfdxConfigContent := "{\"defaultusername\": \"someOrg\"}"
	if err := os.WriteFile(path.Join(sfdxConfigFolder, "sfdx-config.json"), []byte(sfdxConfigContent), 0777); err != nil {
		t.Error("Could not create sfdx config file", err.Error())
		t.FailNow()
	}

	result, err := gosalesforce.FileSystemSfConnectionsManager{}.GetDefaultUsername(testFile)
	if err != nil {
		t.Error("Could not read sfdx username", err.Error())
		t.FailNow()
	}
	if result != "someOrg" {
		t.Error("Wrong username returned, expected", "someOrg", "got", result)
	}
}

// func TestGetDefaultUsername_sfConfig(t *testing.T) {
// 	projectDir := os.TempDir()
// }

func TestGetDefaultUsername_noDefaultOrg(t *testing.T) {
	_, testFile := createTestProjectWithFile(t)

	result, err := gosalesforce.FileSystemSfConnectionsManager{}.GetDefaultUsername(testFile)
	if err == nil {
		t.Error("In folders without sfdx or sf config, function should return error, returned ", result, "instead")
	}

}
