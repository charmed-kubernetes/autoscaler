<<<<<<< HEAD
//go:build go1.13
=======
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
// +build go1.13

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package autorest

import (
	"errors"

	"github.com/Azure/go-autorest/autorest/adal"
)

// IsTokenRefreshError returns true if the specified error implements the TokenRefreshError interface.
func IsTokenRefreshError(err error) bool {
	var tre adal.TokenRefreshError
	return errors.As(err, &tre)
}
