package provisioner

import (
	"errors"
	"fmt"
	"github.com/mikel-at-tatari/tatari-dev-db/aws/rds"
	"github.com/mikel-at-tatari/tatari-dev-db/aws/s3"
	"github.com/mikel-at-tatari/tatari-dev-db/log"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres/conn"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres/pgcommands"
	"time"
)

const DefaultInstanceClass = "db.t2.micro"
const DefaultEngineVersion = "12.10"
const MaxWaitTilReady = time.Duration(60) * time.Second

type RDSDB struct {
	*rds.CreateInstance
}

func (rdsDB *RDSDB) Launch(s3 s3.S3Object, dbName string) error {

	start := time.Now()
	var db *rds.RDSInstanceDescriptor
	var err error

	if db, err = rdsDB.Exec(); err != nil {
		return err
	}

	// wait unit RDS instance is ready. This doesn't necessarily mean the instance is ready to accept connections
	if err = rdsDB.waitTilInstanceIsUp(db); err != nil {
		return err
	}

	// Get the lastest info after the instance is up
	if db, err = rdsDB.getInstanceDescriptor(db); err != nil {
		return err
	}

	connInfo := conn.PGConnInfo{
		DBHost:   db.DBHost,
		DBPort:   db.DBPort,
		DBName:   dbName,
		Username: rdsDB.MasterUser,
		Password: rdsDB.MasterPassword,
	}
	// wait til the DB is up
	if err = rdsDB.waitTilDBAcceptConnections(connInfo); err != nil {
		return err
	}

	// after the DB is up then Restore the DB
	restore := postgres.RestoreFromS3{s3}

	if err = restore.Exec(&connInfo); err != nil {
		return nil
	}

	log.WithFields(log.Fields{
		"Host":             connInfo.DBHost,
		"Port":             connInfo.DBPort,
		"DB":               connInfo.DBName,
		"ElapsedTimeInSec": time.Now().Sub(start) * time.Second,
	}).Info("DB Launched")

	return nil
}

func (rdsDB *RDSDB) getInstanceDescriptor(db *rds.RDSInstanceDescriptor) (*rds.RDSInstanceDescriptor, error) {
	return rds.DescribeInstance(db.DBInstanceId)
}

func (rdsDB *RDSDB) waitTilInstanceIsUp(db *rds.RDSInstanceDescriptor) error {
	return rds.WaitTilAvailable(db.DBInstanceId)
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
		//time.Sleep()

	}

	log.WithFields(log.Fields{
		"Host":          info.DBHost,
		"Port":          info.DBPort,
		"WaitTimeInSec": MaxWaitTilReady,
	}).Error("DB was not ready")

	return errors.New(fmt.Sprintf("DB was not ready in time "))

}
