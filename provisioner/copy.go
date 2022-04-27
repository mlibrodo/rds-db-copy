package provisioner

import (
	"fmt"
	"github.com/mikel-at-tatari/tatari-dev-db/log"
	"github.com/xo/dburl"

	_ "github.com/lib/pq"
)

type CopyDBInput struct {
	DBURL             string
	SourceDBName      string
	DestinationDBName string
}

func CopyDB(in CopyDBInput) error {

	db, err := dburl.Open(in.DBURL)
	if err != nil {
		log.Fatal(err)
		return err
	}

	ddl := fmt.Sprintf("CREATE DATABASE %s WITH TEMPLATE %s", in.DestinationDBName, in.SourceDBName)

	if _, err := db.Query(ddl); err != nil {
		log.Fatal(err)
		return err
	} else {
		log.Debugf("Database %s copied from %s", in.DestinationDBName, in.SourceDBName)
	}

	return nil
}
