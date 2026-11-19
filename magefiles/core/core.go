// Copyright 2026 colonel-byte
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package main is the bootstrap for `go run ./magefiles/core`, which runs the same targets as
// the `mage` CLI without needing it installed. That is the form CI and the pre-commit hooks
// use, so it is the one to reach for in a script or a document.
//
// This file carries no build tag, unlike the namespace files one directory up: it has to be a
// normal main package for `go run` to find it.
package main

import (
	"os"

	"github.com/magefile/mage/mage"
)

func main() {
	os.Exit(mage.Main())
}
