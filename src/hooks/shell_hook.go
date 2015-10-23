package hooks

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"

	"github.com/google/shlex"
	"github.com/tftp-go-team/hooktftp/src/logger"
	"github.com/tftp-go-team/libgotftp/src"
)

// Borrowed from Ruby
// https://github.com/ruby/ruby/blob/v1_9_3_429/lib/shellwords.rb#L82
var shellEscape = regexp.MustCompile("([^A-Za-z0-9_\\-.,:\\/@\n])")

var ShellHook = HookComponents{
	func(command string, _ tftp.Request) (io.ReadCloser, int, error) {

		if len(command) == 0 {
			return nil, -1, errors.New("Empty shell command")
		}

		split, err := shlex.Split(command)
		if err != nil {
			return nil, -1, err
		}

		cmd := exec.Command(split[0], split[1:]...)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, -1, err
		}
		err = cmd.Start()

		// Buffering content to avoid Reader closing after cmd.Wait()
		// For more informations please see:
		//    http://stackoverflow.com/questions/20134095/why-do-i-get-bad-file-descriptor-in-this-go-program-using-stderr-and-ioutil-re
		// Note:
		//    This is not a perfect solution because of buffering. (Memory usage...)
		//    If you have better solution ...
		out, err := ioutil.ReadAll(stdout)
		if err != nil {
			logger.Err("Shell output buffering failed: %v", err)
		}

		// Use goroutine to log the exit status for debugging purposes.
		// XXX: It probably is bad practice to access variables from multiple
		// goroutines, but I hope it is ok in this case since the purpose is
		// only to read and log the status. If not we must remove this bit.
		// Please let me know if you know better.
		go func() {
			err := cmd.Wait()
			if err != nil {
				logger.Err("Command '%v' failed to execute: '%v'", command, err)
			}
		}()

		return ioutil.NopCloser(bytes.NewReader(out)), -1, err
	},
	func(s string) string {
		return shellEscape.ReplaceAllStringFunc(s, func(s string) string {
			return "\\" + s
		})
	},
}
