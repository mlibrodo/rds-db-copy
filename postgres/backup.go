package postgres

import (
	"fmt"
	"github.com/mikel-at-tatari/tatari-dev-db/aws/s3"
	"github.com/mikel-at-tatari/tatari-dev-db/log"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres/conn"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres/pgcommands"
	"os"
	"time"
)

type BackupToS3 struct {
	PGConnInfo *conn.PGConnInfo
}

func (in BackupToS3) Exec(S3 s3.S3Object) error {
	f, err := backup(in.PGConnInfo)

	if err != nil {
		return err
	}

	return s3.Upload(S3, f)
}

func backup(pgConnInfo *conn.PGConnInfo) (*string, error) {
	filePattern := fmt.Sprintf(`%v_%v-*.sql.tar.gz`, pgConnInfo.DBName, time.Now().Unix())
	tempFile, err := os.CreateTemp(pgcommands.TEMP_DIR, filePattern)

	if err != nil {
		return nil, err
	}

	dump := pgcommands.NewPGDump(&pgcommands.Conn{pgConnInfo}, tempFile.Name())
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
			"FullPath": fullPath,
			"Command":  dumpExec.FullCommand,
		},
	).Debug("Backup success")

	return &fullPath, nil
}
