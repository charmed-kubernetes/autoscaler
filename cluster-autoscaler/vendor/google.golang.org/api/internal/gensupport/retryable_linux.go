// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

<<<<<<< HEAD
//go:build linux
=======
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
// +build linux

package gensupport

import "syscall"

func init() {
	// Initialize syscallRetryable to return true on transient socket-level
	// errors. These errors are specific to Linux.
	syscallRetryable = func(err error) bool { return err == syscall.ECONNRESET || err == syscall.ECONNREFUSED }
}
