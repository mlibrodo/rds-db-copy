package provisioner

import (
	"github.com/mikel-at-tatari/tatari-dev-db/aws/s3"
	"github.com/mikel-at-tatari/tatari-dev-db/log"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres/pgcommands"
	"os"
)

type RestoreFromS3Input struct {
	S3         s3.S3Object
	PGConnInfo *postgres.PGConnInfo
}

func RestoreFromS3(in RestoreFromS3Input) error {
	file, _ := os.CreateTemp(TEMP_DIR, "pg_dump-*.sql.tar.gz")
	err := s3.Download(in.S3, file)

	if err != nil {
		return err
	}

	pgConnInfo := in.PGConnInfo

	if err != nil {
		return err
	}

	createDB(pgConnInfo)

	if err != nil {
		return err
	}

	return restoreFromFile(file, pgConnInfo)
}

func createDB(pgConnInfo *postgres.PGConnInfo) error {
	createDB := pgcommands.NewCreateDB(pgConnInfo)
	createDBExec := createDB.Exec()

	if createDBExec.Error != nil {
		log.WithFields(
			log.Fields{
				"Command": createDBExec.FullCommand,
				"Error":   createDBExec.Error.Err,
			},
		).Error("CreateDB failed")

		log.Error(createDBExec.Output)

		return createDBExec.Error.Err

	}
	log.WithFields(
		log.Fields{
			"Command": createDBExec.FullCommand,
		},
	).Debug("CreateDB success")

	log.Debug(createDBExec.Output)

	return nil
}

func restoreFromFile(file *os.File, pgConnInfo *postgres.PGConnInfo) error {

	pgRestore := pgcommands.NewPGRestore(pgConnInfo, file.Name())

	restoreExec := pgRestore.Exec()

	if restoreExec.Error != nil {
		log.WithFields(
			log.Fields{
				"Command": restoreExec.FullCommand,
				"Error":   restoreExec.Error.Err,
			},
		).Error("Restore failed")

		log.Error(restoreExec.Output)

		return restoreExec.Error.Err

	}
	log.WithFields(
		log.Fields{
			"Command": restoreExec.FullCommand,
		},
	).Debug("Restore success")

	log.Debug(restoreExec.Output)

	return nil
}
