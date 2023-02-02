package azureds

import (
	"net/url"
	"fmt"
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
	credential, err := azblob.NewSharedKeyCredential(conf.AccountName, conf.AccountKey)
	if err != nil {
		return nil, err
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	baseUrl, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net", conf.AccountName))
	serviceURL := azblob.NewServiceURL(*baseUrl, pipeline)
	containerURL := serviceURL.NewContainerURL(conf.ContainerName)
	containerURL.Create(context.Background(), azblob.Metadata{}, "")

	return &AzureStorage{
		Config: conf,
	}, nil
}

// GetBlockURL returns the block url of a given key
func (storage *AzureStorage) GetBlockURL(key string) (*azblob.BlockBlobURL, error) {
		// From the Azure portal, get your Storage account blob service URL endpoint.
		accountName := storage.Config.AccountName
		accountKey := storage.Config.AccountKey
		containerName := storage.Config.ContainerName
		folderName := storage.Config.FolderName
	
		// Create a ContainerURL object that wraps a soon-to-be-created blob's URL and a default pipeline.
		u, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s%s", accountName, containerName, folderName, key))
		credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
		if err != nil {
				return nil, err
		}
		blobURL := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))
		return &blobURL, nil
}

// Put adds a key value pair to the storage
func (storage *AzureStorage) Put(ctx context.Context, k ds.Key, value []byte) error {
	blobURL, err := storage.GetBlockURL(k.String())
	if err != nil {
		return err
	}

	// Create a blob with metadata (string key/value pairs)
	// NOTE: Metadata key names are always converted to lowercase before being sent to the Storage Service.
	// Therefore, you should always use lowercase letters; especially when querying a map for a metadata key.
	creatingApp, _ := os.Executable()
	_, err = blobURL.Upload(ctx, bytes.NewReader(value), azblob.BlobHTTPHeaders{},
	azblob.Metadata{"author": "ipfs", "app": creatingApp}, azblob.BlobAccessConditions{})
	if err != nil {
			return err
	}

	return nil
}

// Sync is unimplemented
func (storage *AzureStorage) Sync(ctx context.Context, prefix ds.Key) error {
	return nil
}

// Get gets the data from the desired key
func (storage *AzureStorage) Get(ctx context.Context, k ds.Key) ([]byte, error) {
	blobURL, err := storage.GetBlockURL(k.String())
	if err != nil {
		return nil, err
	}

	response, err := blobURL.Download(ctx, 0, 0, azblob.BlobAccessConditions{}, false)

	if err != nil {
		if stgErr, ok := err.(azblob.StorageError); ok {
			if stgErr.ServiceCode() == azblob.ServiceCodeBlobNotFound {
				return nil, ds.ErrNotFound
			}
		}
		return nil, err
	}

	blobData := &bytes.Buffer{}
	reader := response.Body(azblob.RetryReaderOptions{})
	blobData.ReadFrom(reader)
	reader.Close() // The client must close the response body when finished with it
	return blobData.Bytes(), nil
}

// Has checks if the given key exists
func (storage *AzureStorage) Has(ctx context.Context, k ds.Key) (exists bool, err error) {
	blobURL, err := storage.GetBlockURL(k.String())
	if err != nil {
		return false, err
	}

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
func (storage *AzureStorage) GetSize(ctx context.Context, k ds.Key) (size int, err error) {
	blobURL, err := storage.GetBlockURL(k.String())
	if err != nil {
		return 0, err
	}

	blockList, err := blobURL.GetBlockList(ctx, azblob.BlockListCommitted, azblob.LeaseAccessConditions{})
	if err != nil {
		if stgErr, ok := err.(azblob.StorageError); ok &&
		stgErr.ServiceCode() == azblob.ServiceCodeBlobNotFound {
			return 0, ds.ErrNotFound
		}
		return 0, err
	}
	return int(blockList.BlobContentLength()), nil
}

// Delete deletes the specified key
func (storage *AzureStorage) Delete(ctx context.Context, k ds.Key) error {
	blobURL, err := storage.GetBlockURL(k.String())
	if err != nil {
		return err
	}

	_, err = blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		if stgErr, ok := err.(azblob.StorageError); ok &&
		stgErr.ServiceCode() == azblob.ServiceCodeBlobNotFound {
			return ds.ErrNotFound
		}
		return err
	}
	return nil
}

// Query returns a dsq result
func (storage *AzureStorage) Query(ctx context.Context, q dsq.Query) (dsq.Results, error) {
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
func (storage *AzureStorage) Batch(ctx context.Context) (ds.Batch, error) {
	return &azureBatch{
		storage: storage,
		ops: make(map[string]batchOp),
	}, nil
}

func (batch *azureBatch) Put(ctx context.Context, key ds.Key, val []byte) error {
	batch.ops[key.String()] = batchOp{val: val, delete: false}
	return nil
}

func (batch *azureBatch) Delete(ctx context.Context, key ds.Key) error {
	batch.ops[key.String()] = batchOp{val: nil, delete: true}
	return nil
}

func (batch *azureBatch) Commit(ctx context.Context) error {
	var err error
	for k, op := range batch.ops {
		if op.delete {
			err = batch.storage.Delete(ctx, ds.NewKey(k))
		} else {
			err = batch.storage.Put(ctx, ds.NewKey(k), op.val)
		}
		if err != nil {
			break
		}
	}

	return err
}
