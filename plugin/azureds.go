package plugin

import (
	"fmt"

	"github.com/ipfs/go-ipfs/plugin"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

var Plugins = []plugin.Plugin{
	&AzurePlugin{},
}

type AzurePlugin struct{}

func (azurePlugin AzurePlugin) Name() string {
	return "azure-datastore-plugin"
}

func (azurePlugin AzurePlugin) Version() string {
	return "0.0.1"
}

func (azurePlugin AzurePlugin) Init(env *plugin.Environment) error {
	return nil
}

func (azurePlugin AzurePlugin) DatastoreTypeName() string {
	return "azureds"
}

func (azurePlugin AzurePlugin) DatastoreConfigParser() fsrepo.ConfigFromMap {
	return func(m map[string]interface{}) (fsrepo.DatastoreConfig, error) {
		// parse the m map here

		// return
		return &AzureConfig{
			cfg: azureds.Config{
				// key value pairs of the config
			},
		}, nil
	}
}

type AzureConfig struct {
	cfg azureds.Config
}

func (azureConfig *AzureConfig) DiskSpec() fsrepo.DiskSpec {
	return fsrepo.DiskSpec{
		// key value pair of the minimal representation of a data store
	}
}

func (azureConfig *AzureConfig) Create(path string) (repo.Datastore, error) {
	return azureds.NewAzureDatastore(azureConfig.cfg)
}