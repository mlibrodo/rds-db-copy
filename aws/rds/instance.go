package rds

type RDSInstanceDescriptor struct {
	DBHost        string
	DBPort        int32
	DBInstanceId  *string
	DBInstanceARN *string
}
