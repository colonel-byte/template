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
	"github.com/OWNER/REPO/magefiles/pkg/gen/docs"
	"github.com/magefile/mage/mg"
)

type (
	Generate mg.Namespace
)

// Document regenerates every part of the book that is written from something else: docs/index.md
// from README.md, docs/security.md from .github/SECURITY.md, and docs/SUMMARY.md from the tree.
func (Generate) Document() error {
	return docs.GenerateDocument()
}

// TEMPLATE: the generators that belong here once there is something to generate from, each a
// method delegating straight into a package under magefiles/pkg/gen/ per the rules in AGENTS.md:
//
//   - Schema, reflecting the config structs into JSON schemas and their reference pages.
//   - Completion, writing shell completion scripts for the CLI.
//   - a docs/commands/ pass over the Cobra tree with doc.GenMarkdownTreeCustom, which is what
//     makes the "Commands" section of SUMMARY.md non-empty.
//
// Each one added here needs its output marked in .gitattributes and its path listed in the
// generated-pages table in docs/AGENTS.md.
