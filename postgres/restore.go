package postgres

import (
	"github.com/mlibrodo/rds-db-copy/aws/s3"
	"github.com/mlibrodo/rds-db-copy/log"
	"github.com/mlibrodo/rds-db-copy/postgres/conn"
	"github.com/mlibrodo/rds-db-copy/postgres/pgcommands"
	"os"
)

type RestoreFromS3 struct {
	S3 s3.S3Object
}

func (in RestoreFromS3) Exec(pgConnInfo *conn.PGConnInfo) error {
	file, _ := os.CreateTemp(pgcommands.TEMP_DIR, "pg_dump-*.sql.tar.gz")
	defer file.Close()

	err := s3.Download(in.S3, file)

	if err != nil {
		return err
	}

	err = createDB(pgConnInfo)

	if err != nil {
		return err
	}

	return restoreFromFile(file, pgConnInfo)
}

func createDB(pgConnInfo *conn.PGConnInfo) error {
	createDB := pgcommands.NewCreateDB(&pgcommands.Conn{pgConnInfo})
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

func restoreFromFile(file *os.File, pgConnInfo *conn.PGConnInfo) error {

	pgRestore := pgcommands.NewPGRestore(&pgcommands.Conn{pgConnInfo}, file.Name())

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
