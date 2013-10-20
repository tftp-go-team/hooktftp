package main

import (
	"os/user"
	"strconv"
	"syscall"
)

func DropPrivileges(username string) error {
	userInfo, err := user.Lookup(username)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(userInfo.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(userInfo.Gid)
	if err != nil {
		return err
	}

	// TODO: should set secondary groups too
	err = syscall.Setgroups([]int{gid})
	if err != nil {
		return err
	}

	err = syscall.Setgid(gid)
	if err != nil {
		return err
	}

	err = syscall.Setuid(uid)
	if err != nil {
		return err
	}

	return nil
}

