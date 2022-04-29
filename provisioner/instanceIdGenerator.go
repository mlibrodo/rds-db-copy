package provisioner

import (
	"fmt"
	"time"
)

func GenerateDBInstanceId(prefix string) string {

	return fmt.Sprintf("%s-%v", prefix, time.Now().Unix())
}

func GenerateUnclaimedInstanceId() string {
	s := "unclaimed"
	return GenerateDBInstanceId(s)
}
