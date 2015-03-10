package main

import (
	"os/user"
	"strconv"
	"golang.org/x/sys/unix"
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
	err = unix.Setgroups([]int{gid})
	if err != nil {
		return err
	}

	err = unix.Setregid(gid, gid)
	if err != nil {
		return err
	}

	err = unix.Setreuid(uid, uid)
	if err != nil {
		return err
	}

	return nil
}

