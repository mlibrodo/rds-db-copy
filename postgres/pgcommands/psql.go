package pgcommands

import "fmt"

var (
	// PGDumpCmd is the path to the `pg_dump` executable
	PSQLCmd = "psql"
)

type PSQLQuery struct {
	*Conn
	Query string
}

func (x *PSQLQuery) Exec() Result {
	execFn := GenericExec(PSQLCmd, x.Conn, x.ParseArgs)

	return execFn()

}

func (x *PSQLQuery) ParseArgs() []string {
	var args []string

	if y := x.DBHostAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	if y := x.DBPortAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	args = append(args, fmt.Sprintf("--command=%s", x.Query))

	return args
}