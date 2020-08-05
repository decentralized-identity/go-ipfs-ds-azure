package azureds

import (
	"net/url"
	"fmt"
	"log"
	"context"
	"os"
	"bytes"

	"github.com/Azure/azure-storage-blob-go/azblob"

	ds "github.com/ipfs/go-datastore"
	dsq "github.com/ipfs/go-datastore/query"
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

// Sync is unimplemented
func (storage *AzureStorage) Sync(prefix ds.Key) error {
	return nil
}

// Get gets the data from the desired key
func (storage *AzureStorage) Get(k ds.Key) ([]byte, error) {
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
			return nil, err
	}
	blobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	ctx := context.Background() // This example uses a never-expiring context

	response, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)
	if err != nil {
			log.Fatal(err)
			return nil, err
	}
	blobData := &bytes.Buffer{}
	reader := response.Body(azblob.RetryReaderOptions{})
	blobData.ReadFrom(reader)
	reader.Close() // The client must close the response body when finished with it
	return blobData.Bytes(), nil
}

// Has checks if the given key exists
func (storage *AzureStorage) Has(k ds.Key) (exists bool, err error) {
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
				return false, err
		}
		blobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	
		ctx := context.Background() // This example uses a never-expiring context
		_, err = blobURL.GetBlockList(ctx, azblob.BlockListCommitted, azblob.LeaseAccessConditions{})
		if err != nil {
			if stgErr, ok := err.(azblob.StorageError); ok &&
			stgErr.ServiceCode() == azblob.ServiceCodeBlobNotFound {
				return false, nil
	 		}
			return false, err
		}
		return true, nil
}

// GetSize gets the size of the specified key
func (storage *AzureStorage) GetSize(k ds.Key) (size int, err error) {
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
				return 0, err
		}
		blobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
	
		ctx := context.Background() // This example uses a never-expiring context
		blockList, err := blobURL.GetBlockList(ctx, azblob.BlockListCommitted, azblob.LeaseAccessConditions{})
		if err != nil {
			return 0, err
		}
		return int(blockList.BlobContentLength()), nil
}

// Delete deletes the specified key
func (storage *AzureStorage) Delete(k ds.Key) error {
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
		_, err = blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
		if err != nil {
			return err
		}
		return nil
}

// Query returns a dsq result
func (storage *AzureStorage) Query(q dsq.Query) (dsq.Results, error) {
	return nil, fmt.Errorf("Azure storage query not supported")
}

// Close is not implemented
func (storage *AzureStorage) Close() error {
	return nil
}
