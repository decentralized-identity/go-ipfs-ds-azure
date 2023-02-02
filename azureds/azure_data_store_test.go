package azureds

import (
	"testing"
	"math/rand"
	"strconv"
	"context"
	"time"

	ds "github.com/ipfs/go-datastore"
)

func TestIntegration(t *testing.T) {
	azds, _ := NewAzureDatastore(Config{
		AccountName: "INSERT ACCOUNT NAME",
		AccountKey: "INSERT ACCOUNT KEY",
		ContainerName: "INSERT CONTAINER NAME",
		FolderName: "INSERT FOLDER NAME",
	})

	testKey := ds.NewKey(strconv.Itoa(rand.Int()))

	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Minute)

	// test has false
	has, err := azds.Has(ctx, testKey)
	if err != nil {
		t.Errorf("unexpected error when has expected false but got error %s", err)
	} else if has {
		t.Errorf("unexpected result from has, expected false but got true")
	}

	// test put
	err = azds.Put(ctx, testKey, []byte("test string"))
	if err != nil {
		t.Errorf("unexpected error when put %s", err)
	}

	// test has true 
	has, err = azds.Has(ctx, testKey)
	if err != nil {
		t.Errorf("unexpected error when has expected true but got error %s", err)
	} else if !has {
		t.Errorf("unexpected result from has, expected true but got false")
	}

	// test get
	_, err = azds.Get(ctx, testKey)
	if err != nil {
		t.Errorf("unexpected error when get, got error %s", err)
	}

	// test delete
	err = azds.Delete(ctx, testKey)
	if err != nil {
		t.Errorf("unexpected error when delete got error %s", err)
	}

	// test that delete actually happened
	has, err = azds.Has(ctx, testKey)
	if err != nil {
		t.Errorf("unexpected error when has expected false after delete but got error %s", err)
	} else if has {
		t.Errorf("unexpected result from has after delete, expected false but got true")
	}
}
