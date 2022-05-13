# rds-db-copy

A Golang project.

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ ./bin/rds-db-copy
```

### AWS Infrastructure Requirements
TODO: Elaborate on what infrastructure requirements
1) RDS subnet group name is created and correct ports for PG are configured in the security groups
2) Backup bucket for pgdumps

### Overriding via Environment Variables
To override the default configuration from the [default_config.toml](config/default_config.toml) use a `.` for nested values
i.e. override `server->port` set the `SERVER.PORT` environment variable 


### Logging into your assigned instance
1) Make sure to download the SSL certificates `global-bundle.pem`
   https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html
2) export PGPASSWORD="$(aws rds generate-db-auth-token --hostname <DBHOST> --port <DBPORT> --region <REGION> --username <DBUSER> )"
   PGSSLMODE=verify-full PGSSLROOTCERT=<PATH_TO_AWS_RDS_CERT> psql --host=<DBHOST> --port=<DBPORT>  --dbname=<DBNAME> --username=<DBUSER>

### Testing

``make test``

### TODO
1) make sure docker file has pg_tools
2) Make sure security group assigned in VPC subnets are able to access PG port
3) Increase number of DB instances allowed (default 40 per region) https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Limits.html
4) get aws account id from env var (or be injected somehow)
5) Delete the db instance command
   1) remove inline policy from user