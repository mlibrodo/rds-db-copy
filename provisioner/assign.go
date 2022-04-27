package provisioner

type AssignToIAMUserInput struct {
	IAMUser string
	DBName  string
}

func AssignToIAMUser(in AssignToIAMUserInput) error {

	return nil
}

func createDbUser(user string) error {
	return nil
}

// Daily delete any unused db_instances

// every day we create RDS postgres- using a pgdump script
// create a postgres instance
// go to s3 and load the DB
// have a pg_dump of the whole DB everyday in s3
// create X-instances daily- pg_restore the instance

// Daily delete any unused db_instances

// User request
// rename instance to user_expiration_date
// create the user in the db if it doesnt exist https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.DBAccounts.html
// GRANT rds_iam TO db_userx;
//
// GRANT ALL PRIVILEGES ON DATABASE database_name TO db_userx;
