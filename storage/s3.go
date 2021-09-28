package storage

import "fmt"

type S3ObjectStore struct {

}

func (objectStore *S3ObjectStore) SetObject(key string, data []byte) error {
	return fmt.Errorf("not implemented")
}

func (objectStore *S3ObjectStore) GetObject(key string)([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (objectStore *S3ObjectStore) CheckIfExists(key string)(bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (objectStore *S3ObjectStore) CopyKeyToDir(key, dir string) error {
	return fmt.Errorf("not implemented")
}

