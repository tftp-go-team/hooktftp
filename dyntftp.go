package main

import (
	"github.com/epeli/dyntftp/tftp"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var root = flag.String("root", "/var/lib/tftpboot/", "Serve files from")
var port = flag.String("port", "69", "Port to listen")
var badinternet = flag.Bool("simulate-bad-internet", false, "Simulate bad internet connection for testing purposes")
var configPath = flag.String("config", "/etc/dyntftp.json", "Config file")
var config *Config

func handleRRQ(req *tftp.Request, res *tftp.RRQresponse) {
	path := filepath.Join(*root, req.Path)

	started := time.Now()

	path, err := filepath.Abs(path)

	if err != nil {
		fmt.Println("Bad path", err)
		res.WriteError(tftp.UNKNOWN_ERROR, "Invalid file path:"+err.Error())
		return
	}

	fmt.Println("GET", path, "blocksize", req.Blocksize)

	if !strings.HasPrefix(path, *root) {
		fmt.Println("Path access violation", path)
		res.WriteError(tftp.ACCESS_VIOLATION, "Path access violation")
		return
	}

	if err := res.WriteOACK(); err != nil {
		fmt.Println("Failed to write OACK", err)
		return
	}

	for _, hook := range config.Hooks {
		if hook.Regexp.MatchString(path) {
			// TODO: execute command
			res.Write([]byte("customdata"))
			res.End()
			return
		}
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Failed to open file", path, err)
		res.WriteError(tftp.NOT_FOUND, err.Error())
		return
	}

	defer file.Close()

	b := make([]byte, req.Blocksize)

	totalBytes := 0

	for {
		bytesRead, err := file.Read(b)
		totalBytes += bytesRead

		if err == io.EOF {
			if _, err := res.Write(b[:bytesRead]); err != nil {
				fmt.Println("Failed to write last bytes of the file", err)
				return
			}
			res.End()
			break
		} else if err != nil {
			fmt.Println("Error while reading", file, err)
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
	flag.Parse()
	*root, _ = filepath.Abs(*root)
	var err error
	config, err = ParseConfigFile(*configPath)
	if err != nil {
		fmt.Println("Failed to parse", *configPath, err)
		return
	}

	fmt.Println("flags", *root, *port)
	addr, err := net.ResolveUDPAddr("udp", ":"+*port)
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
		req, res, err := server.WaitForRequest()
		if err != nil {
			fmt.Println("Bad tftp request", err)
			continue
		}

		go handleRRQ(req, res)
	}

}
