package main

import "os"

import hooktftp "github.com/tftp-go-team/hooktftp/internal"

func main() {
	os.Exit(hooktftp.HookTFTP())
}
