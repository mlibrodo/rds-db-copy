package provisioner

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/mlibrodo/rds-db-copy/config"
	"time"

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

	var awsRegion = in.AWSRegionOfDB
	var awsAccountId = in.AWSAccountIdOfDB

	// Get the lastest info after the instance is up
	if db, err = rds.DescribeInstance(&in.DBInstanceID); err != nil {
		return nil, nil, err
	}

	// Need to remove any non-alphanumeric since PG wont accept it
	DBUserName := removeNonAlphaNumeric(in.IAMUser)
	policy := iam.AssignDBToUserPolicy{
		Region:        awsRegion,
		AccountID:     awsAccountId,
		DbiResourceId: db.DBIResourceId,
		DBUserName:    DBUserName,
		AWSUser:       in.IAMUser,
		DBInstanceID:  in.DBInstanceID,
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

	//Figure out how to catch duplicate user creation
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

func AssignCopyToUser(c context.Context, dbConn pgx.Tx, dbCopy DBCopy, iamUser string) (*DBCopy, error) {

	var d *rds.RDSInstanceDescriptor
	var err error

	d, err = rds.DescribeInstance(&dbCopy.RDSInstanceId)

	if err != nil {
		return nil, err
	}

	var dbArn arn.ARN
	// parse the ARN of the db to extract the region and account id
	dbArn, err = arn.Parse(d.DBInstanceARN)

	if err != nil {
		return nil, err
	}

	accountId := dbArn.AccountID
	region := dbArn.Region

	assign := AssignRDSToIAMUser{
		IAMUser:          iamUser,
		DBInstanceID:     d.DBInstanceId,
		DBName:           dbCopy.DbName,
		DBMasterUser:     config.GetConfig().AWS.RDS.MasterUsername,
		DBMasterPassword: config.GetConfig().AWS.RDS.MasterPassword,
		AWSRegionOfDB:    region,
		AWSAccountIdOfDB: accountId,
	}

	d, _, err = assign.AssignToIAMUser()

	if err != nil {
		return nil, err
	}

	_, err = dbConn.Exec(c,
		`
UPDATE db_copies SET claimed_by=$1, claimed_dt=$2, expires_dt=$3
`, iamUser, time.Now(), expireIn(Days30))

	if err != nil {
		return nil, err
	}

	var updated *DBCopy
	// refetch
	updated, err = GetDBCopy(c, dbConn, dbCopy.ID)

	if err != nil {
		return nil, err
	}

	return updated, nil

}

func GetAllUnassignedCopies(c context.Context, dbConn pgx.Tx) ([]*DBCopy, error) {
	var err error
	var copies []*DBCopy

	err = pgxscan.Select(c, dbConn, &copies, `
SELECT * FROM db_copies WHERE claimed_by is null
`)
	if err != nil {
		return nil, err
	}

	return copies, nil
}

// Daily delete any unused db_instances

// every day we create RDS postgres- using a pgdump script
// create a postgres instance
// go to s3 and load the DB
// have a pg_dump of the whole DB everyday in s3
// create X-instances daily- pg_restore the instance

// Daily delete any unused db_instances
