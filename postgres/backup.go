package postgres

import (
	"fmt"
	"github.com/mlibrodo/rds-db-copy/aws/s3"
	"github.com/mlibrodo/rds-db-copy/log"
	"github.com/mlibrodo/rds-db-copy/postgres/conn"
	"github.com/mlibrodo/rds-db-copy/postgres/pgcommands"
	"os"
	"time"
)

type BackupToS3 struct {
	DB *conn.PGConnInfo
}

func (in BackupToS3) Exec(S3 s3.S3Object) error {
	db := in.DB
	backupFile, err := backup(db)

	if err != nil {
		return err
	}

	if err = s3.Upload(S3, backupFile); err != nil {
		return err
	}

	log.WithFields(
		log.Fields{
			"Host":        db.DBHost,
			"Port":        db.DBPort,
			"DB":          db.DBName,
			"Backup File": backupFile,
			"S3 Bucket":   S3.Bucket,
			"S3 Key":      S3.Key,
		},
	).Debug("Backup success")

	return nil
}

func backup(db *conn.PGConnInfo) (*string, error) {
	filePattern := fmt.Sprintf(`%v_%v-*.sql.tar.gz`, db.DBName, time.Now().Unix())
	tempFile, err := os.CreateTemp(pgcommands.TEMP_DIR, filePattern)

	if err != nil {
		return nil, err
	}

	dump := pgcommands.NewPGDump(&pgcommands.Conn{db}, tempFile.Name())
	dump.Verbose = true

	dumpExec := dump.Exec()

	if dumpExec.Error != nil {
		log.WithFields(
			log.Fields{
				"Command": dumpExec.FullCommand,
				"Error":   dumpExec.Error.Err,
			},
		).Debug("Backup failed")

		log.Error(dumpExec.Output)

		return nil, dumpExec.Error.Err

	}
	fullPath := tempFile.Name()

	log.WithFields(
		log.Fields{
			"PGDump Flags": dumpExec.FullCommand,
		},
	).Debug("Backup success")

	return &fullPath, nil
}
