package provisioner

import (
	"errors"
	"fmt"
	"github.com/mlibrodo/rds-db-copy/aws/rds"
	"github.com/mlibrodo/rds-db-copy/aws/s3"
	"github.com/mlibrodo/rds-db-copy/log"
	"github.com/mlibrodo/rds-db-copy/postgres"
	"github.com/mlibrodo/rds-db-copy/postgres/conn"
	"github.com/mlibrodo/rds-db-copy/postgres/pgcommands"
	"time"
)

const DefaultInstanceClass = "db.t2.micro"
const DefaultEngineVersion = "12.10"
const MaxWaitTilReady = time.Duration(60) * time.Second

type RDSDB struct {
	*rds.CreateInstance
}

func (rdsDB *RDSDB) Launch(s3 s3.S3Object, dbName string) (*rds.RDSInstanceDescriptor, error) {

	start := time.Now()
	var db *rds.RDSInstanceDescriptor
	var err error

	if db, err = rdsDB.Exec(); err != nil {
		return nil, err
	}

	// wait unit RDS instance is ready. This doesn't necessarily mean the instance is ready to accept connections
	if err = rds.WaitTilAvailable(db.DBInstanceId); err != nil {
		return nil, err
	}

	// Get the lastest info after the instance is up
	if db, err = rds.DescribeInstance(db.DBInstanceId); err != nil {
		return nil, err
	}

	connInfo := conn.PGConnInfo{
		DBHost:   db.DBHost,
		DBPort:   db.DBPort,
		DBName:   dbName,
		Username: rdsDB.MasterUser,
		Password: rdsDB.MasterPassword,
	}
	// wait til the DB accepts connections
	if err = rdsDB.waitTilDBAcceptConnections(connInfo); err != nil {
		return nil, err
	}

	// after the DB is up then Restore the DB
	restore := postgres.RestoreFromS3{s3}

	if err = restore.Exec(&connInfo); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"Host":             connInfo.DBHost,
		"Port":             connInfo.DBPort,
		"DB":               connInfo.DBName,
		"ElapsedTimeInSec": time.Now().Sub(start) * time.Second,
	}).Info("DB Launched")

	return db, nil
}

func (rdsDB *RDSDB) waitTilDBAcceptConnections(info conn.PGConnInfo) error {

	startTime := time.Now()
	c := &pgcommands.Conn{&info}

	for duration := time.Now().Sub(startTime); duration < MaxWaitTilReady; duration = time.Now().Sub(startTime) {
		isReady := pgcommands.IsReady{Conn: c}
		result := isReady.Exec()

		if result.Error == nil {
			return nil
		}
		time.Sleep(10 * time.Second)

	}

	log.WithFields(log.Fields{
		"Host":          info.DBHost,
		"Port":          info.DBPort,
		"WaitTimeInSec": MaxWaitTilReady,
	}).Error("DB was not ready")

	return errors.New(fmt.Sprintf("DB was not ready in time "))

}
