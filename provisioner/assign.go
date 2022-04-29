package provisioner

import (
	"fmt"
	"github.com/mlibrodo/rds-db-copy/aws"
	"github.com/mlibrodo/rds-db-copy/aws/config"
	"github.com/mlibrodo/rds-db-copy/aws/iam"
	"github.com/mlibrodo/rds-db-copy/aws/rds"
	"github.com/mlibrodo/rds-db-copy/postgres/conn"
	"github.com/mlibrodo/rds-db-copy/postgres/pgcommands"
	"regexp"
)

type AssignRDSToIAMUser struct {
	IAMUser          string
	DBInstanceID     string
	DBName           string
	DBMasterUser     string
	DBMasterPassword string
	AWSRegionOfDB    string
	AWSAccountIdOfDB string
}

/*
	To enable a user to access a Postgres RDS database via IAM Authenticator we need to
	1. Attach a policy with action to `rds-db:connect` for the RDS resource
	2. Create a user in Postgres and grant that user `rds_iam`. This will be the DBUSER
	This is outlined in https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html

	That only gives the user access to the DB Instance but not the database itself. We need to add the following:
	3. GRANT ALL PRIVILEGES ON DATABASE "<DATABASE>" to <DBUSER>;
 	4. GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO <DBUSER>;

*/
func (in AssignRDSToIAMUser) AssignToIAMUser() (db *rds.RDSInstanceDescriptor, dbUser *string, err error) {

	var awsRegion string
	if in.AWSRegionOfDB == "" {
		awsRegion = config.AWSConfig.Region
	}

	var awsAccountId string
	if in.AWSAccountIdOfDB == "" {
		awsAccountId = aws.AWSAccountId
	}

	// Get the lastest info after the instance is up
	if db, err = rds.DescribeInstance(&in.DBInstanceID); err != nil {
		return nil, nil, err
	}

	// Need to remove any non-alphanumeric since PG wont accept it
	DBUserName := removeNonAlphaNumeric(in.IAMUser)
	policy := iam.AssignDBToUserPolicy{
		Region:        awsRegion,
		AccountID:     awsAccountId,
		DbiResourceId: *db.DBIResourceId,
		DBUserName:    DBUserName,

		AWSUser:      in.IAMUser,
		DBInstanceID: in.DBInstanceID,
	}

	// 1. Attach the policy for the user
	if err = policy.AttachPolicyToUser(); err != nil {
		return nil, nil, err
	}

	connInfo := &conn.PGConnInfo{
		DBHost:   db.DBHost,
		DBPort:   db.DBPort,
		DBName:   in.DBName,
		Username: in.DBMasterUser,
		Password: in.DBMasterPassword,
	}

	queries := []string{
		fmt.Sprintf("CREATE USER %s", DBUserName),
		fmt.Sprintf("GRANT rds_iam TO %s", DBUserName),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE \"%s\" to %s", in.DBName, DBUserName),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO %s;", DBUserName),
	}

	for _, q := range queries {
		run := pgcommands.PSQLQuery{
			Conn:  &pgcommands.Conn{PGConnInfo: connInfo},
			Query: q,
		}
		if result := run.Exec(); result.Error != nil {
			return nil, nil, result.Error.Err
		}
	}

	dbUser = &DBUserName
	return db, dbUser, nil
}

func removeNonAlphaNumeric(s string) string {
	// Make a Regex to say we only want letters and numbers
	var reg, _ = regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(s, "")
}

// Daily delete any unused db_instances

// every day we create RDS postgres- using a pgdump script
// create a postgres instance
// go to s3 and load the DB
// have a pg_dump of the whole DB everyday in s3
// create X-instances daily- pg_restore the instance

// Daily delete any unused db_instances
