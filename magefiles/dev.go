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

//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/OWNER/REPO/magefiles/pkg/util"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type (
	Dev mg.Namespace
)

// Vendor re-resolves the module graph and rewrites vendor/ from it. This repository builds from
// vendor/ and is expected to build with no network at all, so a dependency change is not done
// until this has run and the result is committed - see AGENTS.md.
func (Dev) Vendor() error {
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}
	return sh.RunV("go", "mod", "vendor")
}

// VerifyVendor fails if go.mod, go.sum, or vendor/ is not what the module graph implies - which
// is what a dependency added without re-vendoring looks like. `go mod verify` checks the module
// cache rather than the vendor tree, so this runs the consistency check the toolchain does for
// a vendored build and reports it as a failure rather than as a build error somewhere later.
func (Dev) VerifyVendor() error {
	fmt.Println("Verifying vendor/ is consistent with go.mod")
	if err := sh.RunV("go", "mod", "verify"); err != nil {
		return err
	}
	// -mod=vendor with no network: this is the invocation a release has to survive, so it is
	// the one to fail on.
	return sh.RunWithV(map[string]string{"GOFLAGS": "-mod=vendor", "GOPROXY": "off"},
		"go", "build", "./...")
}

// CheckTables checks that every Markdown table in the tree is padded the way the style guide
// in AGENTS.md requires, and prints the expected form of each line that is not. The root
// AGENTS.md (and CLAUDE.md, the symlink to it) are exempt.
func (Dev) CheckTables() error {
	fmt.Println("Checking markdown table padding")
	return util.CheckRepoMarkdownTables()
}
