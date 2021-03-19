package hooktftp

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/user"
	"syscall"
	"time"

	"github.com/tftp-go-team/hooktftp/internal/config"
	"github.com/tftp-go-team/hooktftp/internal/hooks"
	"github.com/tftp-go-team/hooktftp/internal/logger"
	tftp "github.com/tftp-go-team/libgotftp/src"
)

const (
	CONFIG_ERROR = 1
	HOOK_ERROR = 2
	NET_ERROR = 3
	SYS_ERROR = 4
)

var HOOKS []hooks.Hook
var CONFIG_PATH string = "/etc/hooktftp.yml"

func handleRRQ(res *tftp.RRQresponse) {

	started := time.Now()

	path := res.Request.Path

	logger.Info(fmt.Sprintf(
		"GET %s blocksize %d from %s",
		path,
		res.Request.Blocksize,
		*res.Request.Addr,
	))

	defer res.End()

	var hookResult *hooks.HookResult

	for _, hook := range HOOKS {

		var err error

		hookResult, err = hook(res.Request.Path, *res.Request)
		if err == hooks.NO_MATCH {
			continue
		} else if err != nil {

			if err, ok := err.(*os.PathError); ok {
				res.WriteError(tftp.NOT_FOUND, err.Error())
				return
			}

			logger.Err("Failed to execute hook for '%v' error: %v", res.Request.Path, err)
			res.WriteError(tftp.UNKNOWN_ERROR, "Hook failed: "+err.Error())
			return
		}
		break
	}

	if hookResult == nil {
		res.WriteError(tftp.NOT_FOUND, "No hook matches")
		return
	}

	// Consume stderr in a go routine, and defer the call to hook.Finalize.
	if hookResult.Stderr != nil {

		stderrReady := make(chan bool)

		go func() {
			defer func() {
				if err := hookResult.Stderr.Close(); err != nil {
					logger.Err("Failed to close error reader for %s: %s", res.Request.Path, err)
				}
				stderrReady <- true
			}()

			var bytesRead int
			var err error
			b := make([]byte, 4096)

			for ; err != io.EOF; bytesRead, err = hookResult.Stderr.Read(b) {
				if bytesRead > 0 {
					logger.Warning("Hook error: %s", b[:bytesRead])
				}
				if err != nil {
					logger.Err("Error while reading error reader: %s", err)
					return
				}
			}
		}()

		if hookResult.Finalize != nil {
			defer func() {
				<-stderrReady
				err := hookResult.Finalize()
				if err != nil {
					logger.Err("Hook for %v failed to finalize: %v", res.Request.Path, err)
				}
			}()
		}

	}

	// Close stdout before calling Finalize.
	defer func() {
		err := hookResult.Stdout.Close()
		if err != nil {
			logger.Err("Failed to close reader for %s: %s", res.Request.Path, err)
		}
	}()

	if res.Request.TransferSize != -1 {
		res.TransferSize = hookResult.Length
	}

	if err := res.WriteOACK(); err != nil {
		logger.Err("Failed to write OACK: %v", err)
		return
	}

	b := make([]byte, res.Request.Blocksize)

	totalBytes := 0

	for {
		bytesRead, err := hookResult.Stdout.Read(b)
		totalBytes += bytesRead

		if err == io.EOF {
			if _, err := res.Write(b[:bytesRead]); err != nil {
				logger.Err("Failed to write last bytes of the reader: %s", err)
				return
			}
			break
		} else if err != nil {
			logger.Err("Error while reading %s: %s", hookResult.Stdout, err)
			res.WriteError(tftp.UNKNOWN_ERROR, err.Error())
			return
		}

		if _, err := res.Write(b[:bytesRead]); err != nil {
			logger.Err("Failed to write bytes for %s: %s", path, err)
			return
		}
	}

	took := time.Since(started)

	speed := float64(totalBytes) / took.Seconds() / 1024 / 1024

	logger.Info("Sent %v bytes in %v %f MB/s\n", totalBytes, took, speed)
}

func HookTFTP() int {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [-v] [config]\n", os.Args[0])
	}
	verbose := flag.Bool("v", false, "a bool")
	flag.Parse()

	if !*verbose {
		e := logger.Initialize("hooktftp")
		if e != nil {
			log.Fatal("Failed to initialize logger")
		}
	}

	if len(flag.Args()) > 0 {
		CONFIG_PATH = flag.Args()[0]
	}

	logger.Info("Reading hooks from %s", CONFIG_PATH)

	configData, err := ioutil.ReadFile(CONFIG_PATH)

	if err != nil {
		logger.Crit("Failed to read config: %s", err)
		return CONFIG_ERROR
	}

	conf, err := config.ParseYaml(configData)
	if err != nil {
		logger.Crit("Failed to parse config: %s", err)
		return CONFIG_ERROR
	}

	for _, hookDef := range conf.HookDefs {
		logger.Notice("Compiling hook %v", hookDef)

		// Create new hookDef variable for the hookDef pointer for each loop
		// iteration. Go reuses the hookDef variable and if we pass pointer to
		// that terrible things happen.
		newPointer := hookDef
		hook, err := hooks.CompileHook(&newPointer)
		if err != nil {
			logger.Crit("Failed to compile hook %v: %v", hookDef, err)
			return HOOK_ERROR
		}
		HOOKS = append(HOOKS, hook)
	}

	if conf.Port == "" {
		conf.Port = "69"
	}

	addr, err := net.ResolveUDPAddr("udp", conf.Host+":"+conf.Port)
	if err != nil {
		logger.Crit("Failed to resolve address: %s", err)
		return NET_ERROR
	}

	server, err := tftp.NewTFTPServer(addr)
	if err != nil {
		logger.Crit("Failed to listen: %s", err)
		return NET_ERROR
	}

	logger.Notice("Listening on %v:%v", conf.Host, conf.Port)

	if conf.User != "" {
		err := DropPrivileges(conf.User)
		if err != nil {
			logger.Crit("Failed to drop privileges to '%s' error: %v", conf.User, err)
			return SYS_ERROR
		}
		currentUser, _ := user.Current()
		logger.Notice("Dropped %v privileges", currentUser.Username)
	}

	if conf.User == "" && syscall.Getuid() == 0 {
		logger.Warning("Running as root and 'user' is not set in %s", CONFIG_PATH)
	}

	for {
		res, err := server.Accept()
		if err != nil {
			logger.Err("Bad tftp request: %s", err)
			continue
		}

		go handleRRQ(res)
	}

	logger.Close()

	return 0

}
