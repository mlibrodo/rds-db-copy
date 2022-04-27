package pgcommands

import (
	"fmt"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres"
)

var (
	// PGDumpCmd is the path to the `pg_dump` executable
	PGDumpCmd = "pg_dump"
)

type PGDump struct {
	*postgres.PGConnInfo
	// Verbose mode
	Verbose bool
	// File: output file name
	File string
}

func NewPGDump(pgConnInfo *postgres.PGConnInfo, file string) *PGDump {
	return &PGDump{
		PGConnInfo: pgConnInfo,
		File:       file,
	}
}

// Exec `pg_dump` for specified DB
func (x *PGDump) Exec() Result {

	execFn := GenericExec(PGDumpCmd, x.PGConnInfo, x.ParseArgs)

	return execFn()
}

func (x *PGDump) ParseArgs() []string {

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

	if x.Verbose {
		args = append(args, "-v")
	}

	args = append(args, "--format=c")
	args = append(args, "--no-owner")
	args = append(args, "--no-acl")
	args = append(args, "--blob")

	args = append(args, fmt.Sprintf(`--file=%v`, x.File))

	if y := x.DBNameAsCmdArg(nil); y != nil {
		args = append(args, *y)
	}

	return args
}
