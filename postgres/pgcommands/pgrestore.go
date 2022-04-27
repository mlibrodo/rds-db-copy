package pgcommands

import (
	"fmt"
)

var (
	// PGRestoreCmd is the path to the `pg_restore` executable
	PGRestoreCmd = "pg_restore"
)

type PGRestore struct {
	*Conn
	// Verbose mode
	Verbose bool
	// File: Input file name
	File string

	JobCount int
}

func NewPGRestore(pgConnInfo *Conn, file string) *PGRestore {
	return &PGRestore{
		Conn: pgConnInfo,
		File: file,
	}
}

// Exec `pg_restore` for specified DB
func (x *PGRestore) Exec() Result {

	execFn := GenericExec(PGRestoreCmd, x.Conn, x.ParseArgs)

	return execFn()
}

func (x *PGRestore) ParseArgs() []string {

	var args []string

	if y := x.DBHostAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	if y := x.DBPortAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	if y := x.UsernameAsCmdArg(); y != nil {
		args = append(args, *y)
	}

	dbArgKey := "--dbname"
	if y := x.DBNameAsCmdArg(&dbArgKey); y != nil {
		args = append(args, *y)
	}
	if x.Verbose {
		args = append(args, "-v")
	}

	if x.JobCount != 0 {
		args = append(args, fmt.Sprintf("--jobs=%v", x.JobCount))
	}

	args = append(args, "--no-owner")
	args = append(args, "--no-acl")
	args = append(args, "--exit-on-error")
	args = append(args, x.File)

	return args
}
