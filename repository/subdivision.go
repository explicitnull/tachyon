package repository

import (
	"github.com/aerospike/aerospike-client-go"
	"github.com/sirupsen/logrus"
)

// 	GetSubdivisionID returns subdivision ID for DB normalization
func GetSubdivisionID(le *logrus.Entry, aClient *aerospike.Client, div string) (int, error) {
	return 1, nil
}
