package provisioner

import (
	"context"
	"github.com/jackc/pgx/v4"
)

func InitStorage(ctx context.Context, tx pgx.Tx) error {

	var err error

	_, err = tx.Exec(ctx, `
CREATE TABLE db_copies (
	id 			  	SERIAL PRIMARY KEY,
	db_name        	TEXT NOT NULL,
	rds_instance_id TEXT NOT NULL,
	claimed_by     	TEXT,
	created_dt     	TIMESTAMPTZ NOT NULL DEFAULT now(),
	claimed_dt     	TIMESTAMPTZ,
    s3_key			TEXT,
	expires_dt		TIMESTAMPTZ NOT NULL
)
`,
	)

	_, err = tx.Exec(ctx, `
CREATE TABLE rds_instances (
	id 					SERIAL PRIMARY KEY,
	rds_id				TEXT NOT NULL,
	instance_class   	TEXT NOT NULL,
	engine 				TEXT NOT NULL,
	version     		TEXT NOT NULL,
	size				INT NOT NULL,
	created_dt     		TIMESTAMPTZ NOT NULL DEFAULT now()
)
`,
	)

	return err
}
