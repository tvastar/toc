// Copyright (C) 2019 rameshvk. All rights reserved.
// Use of this source code is governed by a MIT-style license
// that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/tvastar/test"
)

func TestListOfContents(t *testing.T) {
	test.File(t.Fatal, "README.input.md", "README.toc.md", func(input string) string {
		inf, outf, h := "testdata/README.input.md", "", "Contents"
		return Generate(&inf, &outf, &h)
	})
}

func TestInplaceModify(t *testing.T) {
	test.File(t.Fatal, "README.input.md", "README.output.md", func(input string) string {
		err := ioutil.WriteFile("testdata/temp.md", []byte(input), 0644)
		must(err, "temp file")
		defer func() {
			must(os.Remove("testdata/temp.md"), "remove temp")
		}()

		inf, h := "testdata/temp.md", "Contents"
		s := Generate(&inf, &inf, &h)
		if s != "" {
			return s
		}
		bytes, err := ioutil.ReadFile("testdata/temp.md")
		must(err, "read output")
		return string(bytes)
	})
}

func TestMissingHeaderInPlace(t *testing.T) {
	ignore := func(args ...interface{}) {}
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if strings.Contains(err.Error(), "Could not locate '##") {
					return
				}
			}
		}

		t.Fatal("Did not panic")
	}()

	test.File(ignore, "README.input.md", "README.toc.md", func(input string) string {
		inf := "testdata/README.input.md"
		h := "Missing Contents"
		return Generate(&inf, &inf, &h)
	})
}

func TestNonExistentInput(t *testing.T) {
	ignore := func(args ...interface{}) {}
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if strings.Contains(err.Error(), "no such file or directory") {
					return
				}
			}
		}

		t.Fatal("Did not panic")
	}()

	test.File(ignore, "README.input.md", "README.toc.md", func(input string) string {
		inf := "testdatax/boo.md"
		h := "Missing Contents"
		return Generate(&inf, &inf, &h)
	})
}
