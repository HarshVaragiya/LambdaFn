package storage

import (
	"os"
	"testing"
)

func TestLocalFileSystemStore_AllTests(t *testing.T) {
	t.Run("LocalFileSystemStoreAllTests", func(t *testing.T) {
		tempDir, err := os.MkdirTemp(os.TempDir(), "lambdaFn-test")
		defer os.RemoveAll(tempDir)
		if err != nil {t.Fatalf("error creating temp dir for testing. error = %v", err)}
		fs := LocalFileSystemStore{pathPrefix: tempDir}
		testObjectData := []byte("this is a test object.")
		testObjectKey := "folder/test_object.txt"
		if err = fs.SetObject(testObjectKey, testObjectData); err != nil {
			t.Fatalf("error saving key [%s]. error = %v", testObjectKey, err)

		}
		returnedObject, err := fs.GetObject(testObjectKey)
		if err != nil {
			t.Fatalf("error retrieving saved key [%s]. error = %v", testObjectKey, err)
		}
		if string(returnedObject) != string(testObjectData) {
			t.Fatalf("retrieved data does not match test data. ")
		}


	})
}
