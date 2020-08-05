package azureds

import (
	"net/url"
	"fmt"
	"log"
	"context"
	"os"
	"strings"
	"bytes"

	"github.com/Azure/azure-storage-blob-go/azblob"

	ds "github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
)

const (
	// listMax is the largest amount of objects you can request from S3 in a list
	// call.
	listMax = 1000

	// deleteMax is the largest amount of objects you can delete from S3 in a
	// delete objects call.
	deleteMax = 1000

	defaultWorkers = 100
)

// AzureStorage is a storage representation
type AzureStorage struct {
	Config
}

// Config representation for all info needed
type Config struct {
	accountName string
	accountKey string
	containerName string
	folderName string
}

// NewAzureDatastore creates an AzureDatastore
func NewAzureDatastore(conf Config) (*AzureStorage, error) {
	return &AzureStorage{
		Config: conf,
	}, nil
}

// Put adds a key value pair to the storage
func (storage *AzureStorage) Put(k ds.Key, value []byte) error {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName := storage.Config.accountName
	accountKey := storage.Config.accountKey
	containerName := storage.Config.containerName
	folderName := storage.Config.folderName

	// Create a ContainerURL object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s%s", accountName, containerName, folderName, k))
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
			log.Fatal(err)
			return err
	}
	blobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	ctx := context.Background() // This example uses a never-expiring context

	// Create a blob with metadata (string key/value pairs)
	// NOTE: Metadata key names are always converted to lowercase before being sent to the Storage Service.
	// Therefore, you should always use lowercase letters; especially when querying a map for a metadata key.
	creatingApp, _ := os.Executable()
	_, err = blobURL.Upload(ctx, bytes.NewReader(value), azblob.BlobHTTPHeaders{},
	azblob.Metadata{"author": "Jeffrey", "app": creatingApp}, azblob.BlobAccessConditions{})
	if err != nil {
			log.Fatal(err)
			return err
	}

	return nil
}

func (s *AzureStorage) Sync(prefix ds.Key) error {
	return nil
}

func (s *AzureStorage) Get(k ds.Key) ([]byte, error) {
}

func (s *AzureStorage) Has(k ds.Key) (exists bool, err error) {
}

func (s *AzureStorage) GetSize(k ds.Key) (size int, err error) {
}

func (s *AzureStorage) Delete(k ds.Key) error {
}

func (s *AzureStorage) Query(q dsq.Query) (dsq.Results, error) {
}

func (s *AzureStorage) Batch() (ds.Batch, error) {
	return &AzureBatch{}, nil
}

func (s *AzureStorage) Close() error {
	return nil
}

type AzureBatch struct {
}

func (b *AzureBatch) Put(k ds.Key, val []byte) error {
}

func (b *AzureBatch) Delete(k ds.Key) error {
}

func (b *AzureBatch) Commit() error {
}

var _ ds.Batching = (*AzureStorage)(nil)