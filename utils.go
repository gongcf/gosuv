package main

import (
	"errors"
	"os"
	"runtime"
	"time"

	"github.com/qiniu/log"
)

var ErrGoTimeout = errors.New("GoTimeoutFunc")

func GoFunc(f func() error) chan error {
	ch := make(chan error)
	go func() {
		ch <- f()
	}()
	return ch
}

func GoTimeoutFunc(timeout time.Duration, f func() error) chan error {
	ch := make(chan error)
	go func() {
		var err error
		select {
		case err = <-GoFunc(f):
			ch <- err
		case <-time.After(timeout):
			log.Debugf("timeout: %v", f)
			ch <- ErrGoTimeout
		}
	}()
	return ch
}

func IsDir(dir string) bool {
	fi, err := os.Stat(dir)
	return err == nil && fi.IsDir()
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}
