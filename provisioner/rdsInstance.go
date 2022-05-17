package provisioner

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/mlibrodo/rds-db-copy/aws/rds"
	"time"
)

type RDSInstance struct {
	ID            int32
	RDSID         string `db:"rds_id"`
	InstanceClass string
	Engine        string
	Version       string
	Size          int32
	CreatedDt     time.Time
}

func NewRDSInstance(
	c context.Context,
	dbConn pgx.Tx,
	instanceClass string,
	subnetGroup string,
	engine string,
	version string,
	username string,
	password string,
	size int32) (*RDSInstance, error) {

	var d *rds.RDSInstanceDescriptor
	var err error
	var id *int

	rdsId := GenerateInstanceId(instanceClass, engine, version, size)
	createInstance := &rds.CreateInstance{
		DBInstanceID:       rdsId,
		InstanceClass:      instanceClass,
		SubnetGroupName:    subnetGroup,
		PubliclyAccessible: true,
		Engine:             engine,
		EngineVersion:      version,
		MasterUser:         username,
		MasterPassword:     password,
		StorageSize:        size,
	}

	if d, err = createInstance.Exec(); err != nil {
		return nil, err
	}
	createdDt := time.Now()

	r := dbConn.QueryRow(c,
		`
INSERT INTO rds_instances(rds_id, instance_class, engine, version, created_dt, size)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id
`, rdsId, instanceClass, engine, version, createdDt, size)

	if err = r.Scan(&id); err != nil {
		return nil, err
	}

	inst := RDSInstance{
		ID:            int32(*id),
		RDSID:         d.DBInstanceId,
		InstanceClass: d.InstanceClass,
		Engine:        d.Engine,
		Version:       d.EngineVersion,
		Size:          size,
		CreatedDt:     createdDt,
	}

	return &inst, nil
}

func GetAllInstances(c context.Context, dbConn pgx.Tx) ([]*RDSInstance, error) {
	var err error
	var instances []*RDSInstance

	err = pgxscan.Select(c, dbConn, &instances, `
SELECT * FROM rds_instances
`)
	if err != nil {
		return nil, err
	}

	return instances, nil
}

func DeleteRDSInstance(c context.Context, dbConn pgx.Tx, dbCopy *RDSInstance) error {
	return nil
}
