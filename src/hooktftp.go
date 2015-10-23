package main

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

	"github.com/tftp-go-team/hooktftp/src/config"
	"github.com/tftp-go-team/hooktftp/src/hooks"
	"github.com/tftp-go-team/hooktftp/src/logger"
	"github.com/tftp-go-team/libgotftp/src"
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

	var outReader, errReader io.ReadCloser
	var len int
	for _, hook := range HOOKS {
		var err error
		outReader, errReader, len, err = hook(res.Request.Path, *res.Request)
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
		defer func() {
			err := outReader.Close()
			if err != nil {
				logger.Err("Failed to close reader for %s: %s", res.Request.Path, err)
			}
		}()
		break
	}

	if errReader != nil {
		go func() {
			defer func() {
				if err := errReader.Close(); err != nil {
					logger.Err("Failed to close error reader for %s: %s", res.Request.Path, err)
				}
			}()

			b := make([]byte, 4096)

			var bytesRead int
			var err error
			for ; err != io.EOF; bytesRead, err = errReader.Read(b) {

				if err != nil {
					logger.Err("Error while reading error reader: %s", err)
					return
				} else {
					logger.Warning("Hook error: %s", b[:bytesRead])
				}

			}
		}()
	}

	if outReader == nil {
		res.WriteError(tftp.NOT_FOUND, "No hook matches")
		return
	}

	if res.Request.TransferSize != -1 {
		res.TransferSize = len
	}

	if err := res.WriteOACK(); err != nil {
		logger.Err("Failed to write OACK", err)
		return
	}

	b := make([]byte, res.Request.Blocksize)

	totalBytes := 0

	for {
		bytesRead, err := outReader.Read(b)
		totalBytes += bytesRead

		if err == io.EOF {
			if _, err := res.Write(b[:bytesRead]); err != nil {
				logger.Err("Failed to write last bytes of the reader: %s", err)
				return
			}
			res.End()
			break
		} else if err != nil {
			logger.Err("Error while reading %s: %s", outReader, err)
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

func main() {

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
		return
	}

	conf, err := config.ParseYaml(configData)
	if err != nil {
		logger.Crit("Failed to parse config: %s", err)
		return
	}

	for _, hookDef := range conf.HookDefs {
		logger.Notice("Compiling hook %s", hookDef)

		// Create new hookDef variable for the hookDef pointer for each loop
		// iteration. Go reuses the hookDef variable and if we pass pointer to
		// that terrible things happen.
		newPointer := hookDef
		hook, err := hooks.CompileHook(&newPointer)
		if err != nil {
			logger.Crit("Failed to compile hook %s: %s", hookDef, err)
			return
		}
		HOOKS = append(HOOKS, hook)
	}

	if conf.Port == "" {
		conf.Port = "69"
	}

	addr, err := net.ResolveUDPAddr("udp", conf.Host+":"+conf.Port)
	if err != nil {
		logger.Crit("Failed to resolve address: %s", err)
		return
	}

	server, err := tftp.NewTFTPServer(addr)
	if err != nil {
		logger.Crit("Failed to listen: %s", err)
		return
	}

	logger.Notice("Listening on %s:%d", conf.Host, conf.Port)

	if conf.User != "" {
		err := DropPrivileges(conf.User)
		if err != nil {
			logger.Crit("Failed to drop privileges to '%s' error: %v", conf.User, err)
			return
		}
		currentUser, _ := user.Current()
		logger.Notice("Dropped privileges to %s", currentUser)
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

}
