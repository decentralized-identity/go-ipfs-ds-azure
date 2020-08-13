package plugin

import (
	"testing"
)

func TestAzurePluginName(t *testing.T) {
	resultingName := AzurePlugin{}.Name()
	expectedName := "azure-datastore-plugin"
	if resultingName != expectedName {
		t.Errorf("Expected %s but got %s", expectedName, resultingName)
	}
}

func TestAzurePluginVersion(t *testing.T) {
	resultingVersion := AzurePlugin{}.Version()
	expectedVersion := "0.0.1"

	if resultingVersion != expectedVersion {
		t.Errorf("Expected %s but got %s", expectedVersion, resultingVersion)
	}
}

func TestAzurePluginDatastoreTypeName(t *testing.T) {
	resultingDatastoreTypeName := AzurePlugin{}.DatastoreTypeName()
	expectedDatastoreTypeName := "azure-data-store"

	if resultingDatastoreTypeName != expectedDatastoreTypeName {
		t.Errorf("Expected %s but got %s", expectedDatastoreTypeName, resultingDatastoreTypeName)
	}
}

func TestAzurePluginDatastoreConfigParser(t *testing.T) {
	configMap := map[string]interface{}{
		"accountName": "accountNameValue",
		"accountKey": "accountKeyValue",
		"containerName": "containerNameValue",
		"folderName": "folderNameValue",
	}
	_, err := AzurePlugin{}.DatastoreConfigParser()(configMap)
	if err != nil {
		t.Errorf("unexpected error %s", err)
	} 
}