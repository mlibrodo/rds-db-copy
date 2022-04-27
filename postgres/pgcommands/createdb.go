package pgcommands

var (
	// PGCreateDBCmd is the path to the `createdb` executable
	PGCreateDBCmd = "createdb"
)

type CreateDB struct {
	*Conn
}

func NewCreateDB(pgConnInfo *Conn) *CreateDB {
	return &CreateDB{Conn: pgConnInfo}
}

// Exec `createdb` for specified DB
func (x *CreateDB) Exec() Result {
	execFn := GenericExec(PGCreateDBCmd, x.Conn, x.ParseArgs)
	return execFn()
}

func (x *CreateDB) ParseArgs() []string {

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

	if y := x.DBNameAsCmdArg(nil); y != nil {
		args = append(args, *y)
	}

	return args
}
