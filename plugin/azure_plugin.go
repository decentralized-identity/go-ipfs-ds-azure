package plugin

import (
	"fmt"

	"github.com/decentralized-identity/go-ipfs-ds-azure/azureds"
	"github.com/ipfs/go-ipfs/plugin"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

// Plugins contains all the plugins in this package
var Plugins = []plugin.Plugin{
	&AzurePlugin{},
}

// AzurePlugin struct
type AzurePlugin struct{}

// Name returns the name of the plugin
func (azurePlugin AzurePlugin) Name() string {
	return "azure-datastore-plugin"
}

// Version returns the version of the plugin
func (azurePlugin AzurePlugin) Version() string {
	return "0.0.1"
}

// Init is not implemented
func (azurePlugin AzurePlugin) Init(env *plugin.Environment) error {
	return nil
}

// DatastoreTypeName returns the name of the data store
func (azurePlugin AzurePlugin) DatastoreTypeName() string {
	return "azure-data-store"
}

// DatastoreConfigParser parses the config map and returns a config struct
func (azurePlugin AzurePlugin) DatastoreConfigParser() fsrepo.ConfigFromMap {
	return func(m map[string]interface{}) (fsrepo.DatastoreConfig, error) {
		// parse the m map here
		accountName, ok := m["accountName"].(string)
		if !ok {
			return nil, fmt.Errorf("no accountName specified")
		}

		accountKey, ok := m["accountKey"].(string)
		if !ok {
			return nil, fmt.Errorf("no accountKey specified")
		}

		containerName, ok := m["containerName"].(string)
		if !ok {
			return nil, fmt.Errorf("no containerName specified")
		}

		folderName, ok := m["folderName"].(string)
		if !ok {
			return nil, fmt.Errorf("no folderName specified")
		}

		return &AzureConfig{
			cfg: azureds.Config{
				AccountName: accountName,
				AccountKey: accountKey,
				ContainerName: containerName,
				FolderName: folderName,
			},
		}, nil
	}
}

// AzureConfig contains the config values
type AzureConfig struct {
	cfg azureds.Config
}

// DiskSpec represents the characteristics that represent a data store, 
// if 2 data stores have different values for any of keys, they will be seen as different data stores
// if 2 data stores have the same values for all of these keys, they will be seen as the same data store
func (azureConfig *AzureConfig) DiskSpec() fsrepo.DiskSpec {
	return fsrepo.DiskSpec{
		"accountName": azureConfig.cfg.AccountName,
		"containerName": azureConfig.cfg.ContainerName,
		"folderName": azureConfig.cfg.FolderName,
	}
}

// Create returns a new azure data store
func (azureConfig *AzureConfig) Create(path string) (repo.Datastore, error) {
	return azureds.NewAzureDatastore(azureConfig.cfg)
}