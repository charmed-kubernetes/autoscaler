// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

<<<<<<< HEAD
//go:build !go1.13
=======
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
// +build !go1.13

package http

import "net/http"

// clonedTransport returns the given RoundTripper as a cloned *http.Transport.
// For versions of Go <1.13, this is not supported, so return nil.
func clonedTransport(rt http.RoundTripper) *http.Transport {
	return nil
}
