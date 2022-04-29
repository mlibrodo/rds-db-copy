package rds

type RDSInstanceDescriptor struct {
	// There should not be any passwords/creds in this structure

	DBHost        string
	DBPort        int32
	DBInstanceId  *string
	DBInstanceARN *string
	DBIResourceId *string
	Region        string
}
