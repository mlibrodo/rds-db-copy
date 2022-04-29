package postgres

import (
	"fmt"
	"time"
)

func MakeS3Key(dbName string) string {
	return fmt.Sprintf(`%v_%v.sql.tar.gz`, dbName, time.Now().Unix())
}
