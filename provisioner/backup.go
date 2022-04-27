package provisioner

import (
	"fmt"
	"github.com/mikel-at-tatari/tatari-dev-db/aws/s3"
	"github.com/mikel-at-tatari/tatari-dev-db/log"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres/pgcommands"
	"os"
	"time"
)

type BackupToS3Input struct {
	PGConnInfo *postgres.PGConnInfo
	S3         s3.S3Object
}

func BackupToS3(in BackupToS3Input) error {
	f, err := backup(in.PGConnInfo)

	if err != nil {
		return err
	}

	return s3.Upload(in.S3, f)
}

func backup(pgConnInfo *postgres.PGConnInfo) (*string, error) {
	filePattern := fmt.Sprintf(`%v_%v-*.sql.tar.gz`, pgConnInfo.DBName, time.Now().Unix())
	tempFile, err := os.CreateTemp(TEMP_DIR, filePattern)

	if err != nil {
		return nil, err
	}

	dump := pgcommands.NewPGDump(pgConnInfo, tempFile.Name())
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
