package main

import (
	"github.com/epeli/hooktftp/tftp"
	"github.com/epeli/hooktftp/config"
	"github.com/epeli/hooktftp/hooks"
	"io/ioutil"
	"fmt"
	"os"
	"io"
	"net"
	"time"
)

var HOOKS []hooks.Hook

func handleRRQ(res *tftp.RRQresponse) {

	started := time.Now()

	path := res.Request.Path

	fmt.Println("GET", path, "blocksize", res.Request.Blocksize)

	if err := res.WriteOACK(); err != nil {
		fmt.Println("Failed to write OACK", err)
		return
	}

	var reader io.Reader
	for _, hook := range HOOKS {
		var err error
		reader, err = hook(res.Request.Path)
		if err == hooks.NO_MATCH {
			continue
		} else if err != nil {

			if err, ok := err.(*os.PathError); ok {
				res.WriteError(tftp.NOT_FOUND, err.Error())
				return
			}

			fmt.Printf("Failed to execute hook for '%v' error: %v", res.Request.Path, err)
			res.WriteError(tftp.UNKNOWN_ERROR, "Hook exec failed")
			return
		}
		break
	}

	if reader == nil {
		res.WriteError(tftp.NOT_FOUND, "No hook matches")
		return
	}

	// TODO: close!!

	b := make([]byte, res.Request.Blocksize)

	totalBytes := 0

	for {
		bytesRead, err := reader.Read(b)
		totalBytes += bytesRead

		if err == io.EOF {
			if _, err := res.Write(b[:bytesRead]); err != nil {
				fmt.Println("Failed to write last bytes of the reader", err)
				return
			}
			res.End()
			break
		} else if err != nil {
			fmt.Println("Error while reading", reader, err)
			res.WriteError(tftp.UNKNOWN_ERROR, err.Error())
			return
		}

		if _, err := res.Write(b[:bytesRead]); err != nil {
			fmt.Println("Failed to write bytes for", path, err)
			return
		}
	}

	took := time.Since(started)

	speed := float64(totalBytes) / took.Seconds() / 1024 / 1024

	fmt.Printf("Sent %v bytes in %v %f MB/s\n", totalBytes, took, speed)
}

func main() {

	configData, err := ioutil.ReadFile("./config_test.yml")
	if err != nil {
		fmt.Println("Failed to read config", err)
		return
	}

	conf, err := config.ParseYaml(configData)
	if err != nil {
		fmt.Println("Failed to parse config", err)
		return
	}

	for _, hookDef := range conf.HookDefs {
		fmt.Println("Compiling hook", hookDef)
		hook, err := hooks.CompileHook(&hookDef)
		if err != nil {
			fmt.Println("Failed to compile hook", hookDef, err)
			return
		}
		HOOKS = append(HOOKS, hook)
	}

	addr, err := net.ResolveUDPAddr("udp", ":"+conf.Port)
	if err != nil {
		fmt.Println("Failed to resolve address", err)
		return
	}

	server, err := tftp.NewTFTPServer(addr)
	if err != nil {
		fmt.Println("Failed to listen", err)
		return
	}

	for {
		res, err := server.Accept()
		if err != nil {
			fmt.Println("Bad tftp request", err)
			continue
		}

		go handleRRQ(res)
	}

}
