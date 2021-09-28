package storage

import "github.com/sirupsen/logrus"

var log = logrus.New()

type ObjectStore interface {
	SetObject(string, []byte) error
	GetObject(string)([]byte, error)
	CheckIfExists(string)(bool, error)
	CopyKeyToDir(string, string) error
}