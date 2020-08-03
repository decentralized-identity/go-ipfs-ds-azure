package azureds

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"

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

type AzureStorage struct {
	Config
}

type Config struct {
	// config values used in the plugin
}

func NewAzureDatastore(conf Config) (*AzureStorage, error) {
}

func (s *AzureStorage) Put(k ds.Key, value []byte) error {
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