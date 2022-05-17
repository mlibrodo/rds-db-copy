package provisioner

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/mlibrodo/rds-db-copy/aws/rds"
	"github.com/mlibrodo/rds-db-copy/aws/s3"
	"github.com/mlibrodo/rds-db-copy/config"
	"github.com/mlibrodo/rds-db-copy/log"
	"github.com/mlibrodo/rds-db-copy/postgres"
	"github.com/mlibrodo/rds-db-copy/postgres/conn"
	"time"
)

const Days30 = time.Duration(24) * time.Hour * 30

func expireIn(d time.Duration) time.Time {
	return time.Now().Add(d)
}

type DBCopy struct {
	ID            int32
	DbName        string
	RDSInstanceId string
	ClaimedBy     *string
	ClaimedDT     *time.Time
	S3Key         string
	CreatedDt     time.Time
	ExpiresDt     time.Time
}

func NewDBCopy(c context.Context, dbConn pgx.Tx, s3Key string, instanceId string) (*DBCopy, error) {

	var db *rds.RDSInstanceDescriptor
	var err error
	var dbName string
	var id *int

	dbName = GenerateRestoreDBName(s3Key)

	s3Obj := s3.S3Object{
		Bucket: config.GetConfig().Backup.S3Bucket,
		Key:    s3Key,
	}

	if db, err = rds.DescribeInstance(&instanceId); err != nil {
		log.Error("Unable to find instance")
		return nil, err
	}
	// wait til the DB accepts connections
	if err = rds.WaitTilAvailable(&instanceId); err != nil {
		log.Error("Unable to check if instance is available")
		return nil, err
	}

	connInfo := conn.PGConnInfo{
		DBHost:   db.DBHost,
		DBPort:   db.DBPort,
		DBName:   dbName,
		Username: config.GetConfig().AWS.RDS.MasterUsername,
		Password: config.GetConfig().AWS.RDS.MasterPassword,
	}

	restore := postgres.RestoreFromS3{s3Obj}

	if err = restore.Exec(&connInfo); err != nil {
		return nil, err
	}

	r := dbConn.QueryRow(c,
		`
INSERT INTO db_copies(db_name, rds_instance_id, s3_key, created_dt, expires_dt)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`, dbName, instanceId, s3Key, time.Now(), expireIn(Days30))

	if err = r.Scan(&id); err != nil {
		return nil, err
	}

	dbCopy := DBCopy{
		ID:            int32(*id),
		DbName:        dbName,
		RDSInstanceId: instanceId,
		S3Key:         s3Key,
		ClaimedBy:     nil,
		ClaimedDT:     nil,
	}

	return &dbCopy, nil
}

func GetDBCopy(c context.Context, dbConn pgx.Tx, dbCopyId int32) (*DBCopy, error) {

	var dbCopy DBCopy
	var err error

	err = pgxscan.Get(c, dbConn, &dbCopy, `
SELECT id, db_name, rds_instance_id, claimed_by, created_dt, claimed_dt, s3_key, expires_dt FROM db_copies where id=$1
`, dbCopyId)

	if err != nil {
		return nil, err
	}
	return &dbCopy, nil
}
