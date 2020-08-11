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
	AccountName string
	AccountKey string
	ContainerName string
	FolderName string
}

// NewAzureDatastore creates an AzureDatastore
func NewAzureDatastore(conf Config) (*AzureStorage, error) {
	return &AzureStorage{
		Config: conf,
	}, nil
}

// Put adds a key value pair to the storage
func (storage *AzureStorage) Put(k ds.Key, value []byte) error {
	fmt.Printf("Put is called on %s", k.String())
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName := storage.Config.AccountName
	accountKey := storage.Config.AccountKey
	containerName := storage.Config.ContainerName
	folderName := storage.Config.FolderName

	// Create a ContainerURL object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s%s", accountName, containerName, folderName, k.String()))
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
	fmt.Printf("Get is called on %s", k.String())
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName := storage.Config.AccountName
	accountKey := storage.Config.AccountKey
	containerName := storage.Config.ContainerName
	folderName := storage.Config.FolderName

	// Create a ContainerURL object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s%s", accountName, containerName, folderName, k.String()))
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
	fmt.Printf("Has is called on %s", k.String())
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName := storage.Config.AccountName
	accountKey := storage.Config.AccountKey
	containerName := storage.Config.ContainerName
	folderName := storage.Config.FolderName

	// Create a ContainerURL object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s%s", accountName, containerName, folderName, k.String()))
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
			fmt.Println("not found")
			return false, nil
		}
		return false, err
	}
	fmt.Println("found")
	return true, nil
}

// GetSize gets the size of the specified key
func (storage *AzureStorage) GetSize(k ds.Key) (size int, err error) {
	// From the Azure portal, get your Storage account blob service URL endpoint.
	accountName := storage.Config.AccountName
	accountKey := storage.Config.AccountKey
	containerName := storage.Config.ContainerName
	folderName := storage.Config.FolderName

	// Create a ContainerURL object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s%s", accountName, containerName, folderName, k.String()))
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
	accountName := storage.Config.AccountName
	accountKey := storage.Config.AccountKey
	containerName := storage.Config.ContainerName
	folderName := storage.Config.FolderName

	// Create a ContainerURL object that wraps a soon-to-be-created blob's URL and a default pipeline.
	u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s%s", accountName, containerName, folderName, k.String()))
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

type op struct {
	delete bool
	value  []byte
}

// basicBatch implements
type batchOp struct {
	val    []byte
	delete bool
}

type azureBatch struct {
	storage          *AzureStorage
	ops        map[string]batchOp
}

// Batch returns a batch struct that can take more ops or be committed
func (storage *AzureStorage) Batch() (ds.Batch, error) {
	return &azureBatch{
		storage: storage,
		ops: make(map[string]batchOp),
	}, nil
}

func (batch *azureBatch) Put(key ds.Key, val []byte) error {
	batch.ops[key.String()] = batchOp{val: val, delete: false}
	return nil
}

func (batch *azureBatch) Delete(key ds.Key) error {
	batch.ops[key.String()] = batchOp{val: nil, delete: true}
	return nil
}

func (batch *azureBatch) Commit() error {
	var err error
	for k, op := range batch.ops {
		if op.delete {
			err = batch.storage.Delete(ds.NewKey(k))
		} else {
			err = batch.storage.Put(ds.NewKey(k), op.val)
		}
		if err != nil {
			break
		}
	}

	return err
}
