# tatari-dev-db

A Golang project.

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ ./bin/tatari-dev-db
```

### Testing

``make test``

### TODO
1) make sure docker file has pg_tools
2) Make sure security group assigned in VPC subnets are able to access PG port
3) Increase number of DB instances allowed (default 40 per region) https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_Limits.html

# when deleting db
1) delete the db instance
2) remove inline policy from user