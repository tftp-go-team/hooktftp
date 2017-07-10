package hooks

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/google/shlex"
	"github.com/tftp-go-team/hooktftp/src/config"
	"github.com/tftp-go-team/libgotftp/src"
)

// Borrowed from Ruby
// https://github.com/ruby/ruby/blob/v1_9_3_429/lib/shellwords.rb#L82
var shellEscape = regexp.MustCompile("([^A-Za-z0-9_\\-.,:\\/@\n])")

var ShellHook = HookComponents{
	func(command string, request tftp.Request) (*HookResult, error) {

		if len(command) == 0 {
			return nil, errors.New("Empty shell command")
		}

		split, err := shlex.Split(command)
		if err != nil {
			return nil, err
		}

		cmd := exec.Command(split[0], split[1:]...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return nil, err
		}

		env := os.Environ()
		env = append(env, fmt.Sprintf("CLIENT_ADDR=%s", (*request.Addr).String()))
		cmd.Env = env

		err = cmd.Start()

		// cmd.Wait is used as Finalize to log the exit status.
		return newHookResult(stdout, stderr, -1, cmd.Wait), err
	},
	func(s string, _ config.HookExtraArgs) (string, error) {
		return shellEscape.ReplaceAllStringFunc(s, func(s string) string {
			return "\\" + s
		}), nil
	},
}
