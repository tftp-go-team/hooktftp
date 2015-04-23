package hooks

import (
	"fmt"
	"io"
	"os/exec"
	"regexp"
)

// Borrowed from Ruby
// https://github.com/ruby/ruby/blob/v1_9_3_429/lib/shellwords.rb#L82
var shellEscape = regexp.MustCompile("([^A-Za-z0-9_\\-.,:\\/@\n])")

var ShellHook = HookComponents{
	func(command string) (io.ReadCloser, error) {
		cmd := exec.Command("sh", "-c", command)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return nil, err
		}
		err = cmd.Start()

		// Use goroutine to log the exit status for debugging purposes.
		// XXX: It probably is bad practice to access variables from multiple
		// goroutines, but I hope it is ok in this case since the purpose is
		// only to read and log the status. If not we must remove this bit.
		// Please let me know if you know better.
		go func() {
			err := cmd.Wait()
			if err != nil {
				fmt.Printf("Command '%v' failed to execute: '%v'", command, err)
			}
		}()

		return stdout, err
	},
	func(s string) string {
		return shellEscape.ReplaceAllStringFunc(s, func(s string) string {
			return "\\" + s
		})
	},
}
