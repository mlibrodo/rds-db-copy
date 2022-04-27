package pgcommands

import (
	"bufio"
	"fmt"
	"github.com/mikel-at-tatari/tatari-dev-db/postgres"
	"io"
	"os"
	"os/exec"
	"strings"
)

func streamExecOutput(out io.ReadCloser) string {
	output := ""
	reader := bufio.NewReader(out)
	line, err := reader.ReadString('\n')
	output += line
	for err == nil {
		fmt.Printf(line)
		line, err = reader.ReadString('\n')
		output += line
	}
	return output
}

func GenericExec(pgCommand string, pgConnInfo *postgres.PGConnInfo, parseArgFn func() []string) func() Result {

	return func() Result {

		result := Result{}

		cmdArgs := parseArgFn()
		result.FullCommand = strings.Join(cmdArgs, " ")
		cmd := exec.Command(pgCommand, cmdArgs...)

		if y := pgConnInfo.PasswordAsEnv(); y != nil {
			cmd.Env = append(os.Environ(), *y)
		}

		stderrIn, _ := cmd.StderrPipe()
		go func() {
			result.Output = streamExecOutput(stderrIn)
		}()
		cmd.Start()
		err := cmd.Wait()

		if exitError, ok := err.(*exec.ExitError); ok {
			result.Error = &ResultError{Err: err, ExitCode: exitError.ExitCode(), CmdOutput: result.Output}
		}
		return result
	}
}
