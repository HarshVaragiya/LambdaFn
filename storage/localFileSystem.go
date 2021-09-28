package storage

import (
	"fmt"
	"io/ioutil"
	"path"
)

type LocalFileSystemStore struct {
	pathPrefix		string
}

func (objectStore *LocalFileSystemStore) SetObject(key string, data []byte) error {
	log.Debugf("attempting to store key [%s] on local file system", key)
	fullObjectPath := objectStore.getFullObjectPath(key)
	if err := ioutil.WriteFile(fullObjectPath, data, 0666); err != nil {
		log.Errorf("error saving key [%s] on disk. error = %v", key, err)
		return err
	}
	return nil
}

func (objectStore *LocalFileSystemStore) GetObject(key string)([]byte, error) {
	log.Debugf("attempting to load object with key [%s] from local file system", key)
	fullObjectPath := objectStore.getFullObjectPath(key)
	if data, err := ioutil.ReadFile(fullObjectPath); err != nil{
		log.Errorf("error reading key [%s] from disk. error = %v", key, err)
		return data, err
	} else {
		return data, nil
	}
}

func (objectStore *LocalFileSystemStore) CheckIfExists(string)(bool, error) {
	return false, fmt.Errorf("not implemented")
}

func (objectStore *LocalFileSystemStore) CopyKeyToDir(key, filename string) error {
	log.Debugf("copying key [%s] to [%s]", key, filename)
	fullObjectPath := objectStore.getFullObjectPath(key)
	fileBytes, err := ioutil.ReadFile(fullObjectPath)
	if err != nil {
		log.Errorf("error loading object with key [%s]. error = %v", key,  err)
		return err
	}
	if err = ioutil.WriteFile(filename, fileBytes, 0666); err != nil {
		log.Errorf("error writing to required file. error = %v", err)
		return err
	}
	return nil
}

func (objectStore *LocalFileSystemStore) getFullObjectPath(key string) string {
	fullPath := path.Join(objectStore.pathPrefix, key)
	log.Tracef("key [%s] will have full object path [%s]", key, fullPath)
	return fullPath
}