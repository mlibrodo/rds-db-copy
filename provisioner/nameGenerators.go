package provisioner

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// time format
const yyyyMMddHHmm = "200601021504"
const fileExt = ".sql.tar.gz"

func GenerateInstanceId(instanceClass string, engine string, version string, size int32) string {
	prefix := fmt.Sprintf("%s-%s-%s-%dGB", instanceClass, engine, version, size)
	instanceId := fmt.Sprintf("%s-%s", prefix, time.Now().Format(yyyyMMddHHmm))

	// remove any .
	return strings.ReplaceAll(instanceId, ".", "")
}

// Create a DBName from an S3 key. Key should be the s3 object used for the restore
// For a given s3 key FooBar:200601021504.sql.tar.gz, This creates a dbname
// FooBar_200601021504_xyz, where xyz is a 3 letter random characters
func GenerateRestoreDBName(s3key string) string {

	dbNameTimePair := strings.ReplaceAll(s3key, fileExt, "")

	splits := strings.Split(dbNameTimePair, ":")

	db := splits[0]
	t := splits[1]

	randStr := randSeq(3)

	return fmt.Sprintf("%s_%s_%s", db, t, randStr)
}

// Create a S3 key. For a given dbName FooBar, creates a key of the form
// FooBar:200601021504.sql.tar.gz
func MakeS3BackupKey(dbName string) string {
	return fmt.Sprintf(`%s:%s%s`, dbName, time.Now().Format(yyyyMMddHHmm), fileExt)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	rand.Seed(time.Now().Unix())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
