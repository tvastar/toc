// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

// Command main generates a TOC markdown representation of
// the provided markdown file.
//
// The header option specifies what the header should be called.
//
// The output option can be used to rewrite a similarly named
// section in the same markdown file
//
// Usage
//
//     $ toc README.md > TOC.md
//     $ toc -h "Walkthrough" -o README.md README.md
//
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/russross/blackfriday.v2"
)

var output = flag.String("o", "", "optional output markdown file to write the ToC into")
var header = flag.String("h", "Table of Contents", "The header to use for the ToC")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s [options] file.md:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.Arg(0) == "" {
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Arg(0)
	fmt.Println(Generate(&input, output, header))
}

// Generate does bulk of the work and is exposed here for testing purposes
func Generate(input, output, header *string) string {
	toc := "## " + *header + "\n"
	counts := []int{0}

	src, err := ioutil.ReadFile(*input)
	must(err, *input)

	opts := blackfriday.WithExtensions(blackfriday.CommonExtensions | blackfriday.AutoHeadingIDs)
	blackfriday.New(opts).Parse(src).Walk(
		func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			if entering && n.Type == blackfriday.Heading {
				toc, counts = processHeading(n, toc, counts, header)
			}
			return blackfriday.GoToNext
		})

	// must end in \n\n so we can search for section easily
	toc += "\n"

	if *output == "" {
		return toc
	}

	data, err := ioutil.ReadFile(*output)
	must(err, *output)
	s := string(data)
	idx := strings.Index(s, "## "+*header+"\n")
	if idx < 0 {
		panic(fmt.Errorf("Could not locate '## %s' line in %s.", *header, *output))

	}
	before, after := s[:idx], s[idx:]
	idx = strings.Index(after, "\n\n")
	if idx > 0 {
		after = after[idx+2:]
	} else {
		after = ""
	}
	s = before + toc + after
	must(ioutil.WriteFile(*output, []byte(s), 0644), *output)
	return ""
}

func processHeading(n *blackfriday.Node, toc string, counts []int, header *string) (string, []int) {
	if n.HeadingData.Level == 1 {
		return toc, counts
	}
	l := n.HeadingData.Level - 1

	if l < len(counts) {
		counts = counts[:l]
	}

	for l > len(counts) {
		counts = append(counts, 0)
	}

	name := string(n.Literal)
	text := n.FirstChild
	for name == "" && text != nil {
		if text.Type == blackfriday.Text {
			name = string(text.Literal)
		}
		text = text.Next
	}

	if l == 1 && name == *header {
		return toc, counts
	}

	counts[l-1]++
	id := n.HeadingData.HeadingID
	prefix := strings.Repeat("    ", l-1)
	toc += fmt.Sprintf("%s%d. [%s](#%s)\n", prefix, counts[l-1], name, id)
	return toc, counts
}

var errlog = log.New(os.Stderr, "toc ", 0)

func must(err error, msg string) {
	if err != nil {
		errlog.Println(msg)
		panic(err)
	}
}
