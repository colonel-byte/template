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

package docs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// printExcludedNote marks the sections that docs/css/print.css hides from print.html, so the
// omission is visible to anyone reading SUMMARY.md.
const printExcludedNote = "<!-- Excluded from the print page (print.html) by docs/css/print.css. -->"

// summarySections is the book's chapter order, and it is load-bearing beyond the sidebar:
// docs/css/print.css drops everything from one heading to the end of the document out of the
// print page, so the sections carrying printExcludedNote have to stay last. Anything added
// after them is excluded from print too.
//
// A section whose folder does not exist is skipped rather than failing the run, so this list
// can name a directory before the repository has one.
var summarySections = []struct {
	// title is the section's H1 in SUMMARY.md.
	title string

	// folder is the directory under docs/ whose pages become the section's entries. Empty
	// for a section that supplies its own.
	folder string

	// pattern matches the filenames to include, with one capture group for the link text.
	pattern string

	// extra is written verbatim ahead of the folder walk's entries.
	extra string

	// note is a comment written under the title.
	note string

	// indent steps a bullet in one level per underscore-separated ancestor in the filename,
	// which is how a Cobra subcommand page nests under its group.
	indent bool
}{
	{
		// Both entries are prefix chapters: mdBook gives a line before the first list item
		// no number and places it ahead of every numbered chapter, which is where the front
		// page and the security policy belong.
		title: "Index",
		extra: "[readme](index.md)\n[security](security.md)",
	},
	{
		title:   "Guides",
		folder:  "guides",
		pattern: `(.+)\.md`,
	},
	// TEMPLATE: a "Commands" section goes here once something generates docs/commands/, and it
	// is the one that needs indent. Cobra names a subcommand's page after its whole command
	// path joined with underscores, so the pattern captures the tail:
	//
	//   {title: "Commands", folder: "commands", pattern: `<project>_(.+)\.md`, indent: true,
	//    extra: "- [<project>](commands/<project>.md)"},
	{
		title:   "Development",
		folder:  "dev",
		pattern: `(.+)\.md`,
		note:    printExcludedNote,
	},
	{
		title:   "Workflows",
		folder:  "workflows",
		pattern: `(.+)\.md`,
		note:    printExcludedNote,
	},
	{
		title:   "Agent",
		folder:  "agent",
		pattern: `(.+)\.md`,
		note:    printExcludedNote,
	},
}

// generateSummary writes docs/SUMMARY.md, mdBook's table of contents, from the tree the rest of
// the generators just produced.
func generateSummary() error {
	fmt.Println("docs/SUMMARY.md")

	var b strings.Builder
	for i, section := range summarySections {
		b.WriteString("# " + section.title + "\n\n")
		if section.note != "" {
			b.WriteString(section.note + "\n\n")
		}
		if section.extra != "" {
			b.WriteString(section.extra + "\n")
		}
		if section.folder != "" && section.pattern != "" {
			entries, err := sectionEntries(section.folder, section.pattern, section.indent)
			if err != nil {
				return err
			}
			for _, e := range entries {
				b.WriteString(e + "\n")
			}
		}
		if i != len(summarySections)-1 {
			b.WriteString("\n-----------\n\n")
		}
	}

	return os.WriteFile(filepath.Join(docsDir, "SUMMARY.md"), []byte(b.String()), 0o644)
}

// sectionEntries lists one section's pages as Markdown bullets, sorted by filename so that the
// sidebar does not reorder itself with the filesystem.
//
// A folder that does not exist yields no entries and no error: the section list names
// directories a repository may not have yet, and an empty chapter is the right answer for one
// that is not there. Any other read error is returned.
func sectionEntries(folder, pattern string, indent bool) ([]string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(docsDir, folder)
	dirEntries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var list []string
	for _, entry := range dirEntries {
		if entry.IsDir() || !re.MatchString(entry.Name()) {
			continue
		}
		name := re.FindStringSubmatch(entry.Name())[1]
		if indent {
			path := strings.Split(name, "_")
			list = append(list, fmt.Sprintf("%s- [%s](%s/%s)",
				strings.Repeat("  ", len(path)), path[len(path)-1], folder, entry.Name()))
			continue
		}
		list = append(list, fmt.Sprintf("- [%s](%s/%s)", name, folder, entry.Name()))
	}

	sort.Strings(list)
	return list, nil
}
