package notes

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListCmd(t *testing.T) {
	old := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = old }()

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg := &Config{HomePath: filepath.Join(cwd, "testdata", "list", "normal")}

	format := func(s string) string {
		ss := strings.Split(strings.TrimPrefix(s, "\n"), "\n")
		for i, s := range ss {
			ss[i] = strings.Replace(filepath.FromSlash(strings.TrimLeft(s, "\t")), "HOME", cfg.HomePath, 1)
		}
		return strings.Join(ss, "\n")
	}

	for _, tc := range []struct {
		what string
		cmd  *ListCmd
		want string
	}{
		{
			what: "default",
			cmd:  &ListCmd{},
			want: format(`
			HOME/a/4.md
			HOME/a/1.md
			HOME/c/5.md
			HOME/b/2.md
			HOME/c/3.md
			HOME/b/6.md
			`),
		},
		{
			what: "sort by created",
			cmd: &ListCmd{
				SortBy: "created",
			},
			want: format(`
			HOME/a/4.md
			HOME/a/1.md
			HOME/c/5.md
			HOME/b/2.md
			HOME/c/3.md
			HOME/b/6.md
			`),
		},
		{
			what: "sort by filename",
			cmd: &ListCmd{
				SortBy: "filename",
			},
			want: format(`
			HOME/a/1.md
			HOME/b/2.md
			HOME/c/3.md
			HOME/a/4.md
			HOME/c/5.md
			HOME/b/6.md
			`),
		},
		{
			what: "sort by category",
			cmd: &ListCmd{
				SortBy: "category",
			},
			want: format(`
			HOME/a/1.md
			HOME/a/4.md
			HOME/b/2.md
			HOME/b/6.md
			HOME/c/3.md
			HOME/c/5.md
			`),
		},
		{
			what: "relative paths",
			cmd: &ListCmd{
				Relative: true,
			},
			want: format(`
			a/4.md
			a/1.md
			c/5.md
			b/2.md
			c/3.md
			b/6.md
			`),
		},
		{
			what: "relative paths sorted by file name",
			cmd: &ListCmd{
				Relative: true,
				SortBy:   "filename",
			},
			want: format(`
			a/1.md
			b/2.md
			c/3.md
			a/4.md
			c/5.md
			b/6.md
			`),
		},
		{
			what: "oneline",
			cmd: &ListCmd{
				Oneline: true,
			},
			want: format(`
			a/4.md a bar        4
			a/1.md a foo,bar    this is title
			c/5.md c a-bit-long this is title
			b/2.md b foo        this is title
			c/3.md c            this is title
			b/6.md b future     text from future
			`),
		},
		{
			what: "oneline sorted by category",
			cmd: &ListCmd{
				Oneline: true,
				SortBy:  "category",
			},
			want: format(`
			a/1.md a foo,bar    this is title
			a/4.md a bar        4
			b/2.md b foo        this is title
			b/6.md b future     text from future
			c/3.md c            this is title
			c/5.md c a-bit-long this is title
			`),
		},
		{
			what: "filter by category",
			cmd: &ListCmd{
				Category: "a",
			},
			want: format(`
			HOME/a/4.md
			HOME/a/1.md
			`),
		},
		{
			what: "filter by category with regex sorted by filename",
			cmd: &ListCmd{
				Category: "^(b|c)$",
				SortBy:   "filename",
			},
			want: format(`
			HOME/b/2.md
			HOME/c/3.md
			HOME/c/5.md
			HOME/b/6.md
			`),
		},
		{
			what: "filter by unknown category",
			cmd: &ListCmd{
				Category: "unknown-category-who-know",
			},
			want: format(`
			`),
		},
		{
			what: "filter by tag",
			cmd: &ListCmd{
				Tag: "foo",
			},
			want: format(`
			HOME/a/1.md
			HOME/b/2.md
			`),
		},
		{
			what: "filter by tag with regex sorted by filename",
			cmd: &ListCmd{
				Tag:    "^(foo|future)$",
				SortBy: "filename",
			},
			want: format(`
			HOME/a/1.md
			HOME/b/2.md
			HOME/b/6.md
			`),
		},
		{
			what: "filter by unknown tag",
			cmd: &ListCmd{
				Tag: "unknown-category-who-know",
			},
			want: format(`
			`),
		},
		{
			what: "filter by category and tag",
			cmd: &ListCmd{
				Category: "a",
				Tag:      "foo",
			},
			want: format(`
			HOME/a/1.md
			`),
		},
		{
			what: "full",
			cmd: &ListCmd{
				Full: true,
			},
			want: format(`
			HOME/a/4.md
			Category: a
			Tags:     bar
			Created:  2017-10-30T11:37:45+09:00
			
			4
			=
			
			this
			is
			old text
			
			HOME/a/1.md
			Category: a
			Tags:     foo, bar
			Created:  2018-10-30T11:37:45+09:00
			
			this is title
			=============
			
			this
			is
			test
			
			HOME/c/5.md
			Category: c
			Tags:     a-bit-long
			Created:  2018-10-30T11:37:45+09:00
			
			this is title
			=============
			
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			is
			test
			this
			
			HOME/b/2.md
			Category: b
			Tags:     foo
			Created:  2018-11-01T11:37:45+09:00
			
			this is title
			=============
			
			Lorem ipsum dolor sit amet, his no stet volumus sententiae. Usu id postea animal consetetur. Eum repudiare laboramus conclusionemque et, veritus tractatos dignissim duo ut. Ex sed quod admodum indoctu
			
			HOME/c/3.md
			Category: c
			Tags:     
			Created:  2018-12-30T11:37:45+09:00
			
			this is title
			=============
			
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			this
			is
			newer
			
			HOME/b/6.md
			Category: b
			Tags:     future
			Created:  2118-10-30T11:37:45+09:00
			
			text from future
			================
			
			Lorem ipsum dolor sit amet, his no stet volumus sententiae. Usu id postea animal
			consetetur. Eum repudiare laboramus conclusionemque et, veritus tractatos dignissim
			duo ut. Ex sed quod admodum indoctu
			
			`),
		},
		{
			what: "full with filter",
			cmd: &ListCmd{
				Full:     true,
				Category: "a",
				Tag:      "foo",
			},
			want: format(`
			HOME/a/1.md
			Category: a
			Tags:     foo, bar
			Created:  2018-10-30T11:37:45+09:00
			
			this is title
			=============
			
			this
			is
			test
			
			`),
		},
	} {
		t.Run(tc.what, func(t *testing.T) {
			var buf bytes.Buffer
			tc.cmd.Config = cfg
			tc.cmd.Out = &buf

			if err := tc.cmd.Do(); err != nil {
				t.Fatal(err)
			}

			have := buf.String()
			if tc.want != have {
				ls := strings.Split(tc.want, "\n")
				hint := ""
				for i, l := range strings.Split(have, "\n") {
					if l != ls[i] {
						hint = fmt.Sprintf("first mismatch at line %d: want:%#v v.s. have:%#v", i+1, ls[i], l)
						break
					}
				}
				t.Fatalf("have:\n%s\n\n%s", have, hint)
			}
		})
	}
}

func TestListNoNote(t *testing.T) {
	dir := "test-for-list-empty"
	cfg := &Config{HomePath: dir}
	if err := os.Mkdir(dir, 0755); err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	for _, c := range []*ListCmd{
		&ListCmd{},
		&ListCmd{Oneline: true},
		&ListCmd{Relative: true},
		&ListCmd{Full: true},
	} {
		var b bytes.Buffer
		c.Config = cfg
		c.Out = &b
		if err := c.Do(); err != nil {
			t.Fatal(err)
		}

		out := b.String()
		if out != "" {
			t.Fatalf("should not output anything: %#v", out)
		}
	}
}

func TestListNoHome(t *testing.T) {
	cfg := &Config{HomePath: "/path/to/unknown/directory"}
	err := (&ListCmd{Config: cfg}).Do()
	if err == nil {
		t.Fatal("Error did not occur")
	}
	if !strings.Contains(err.Error(), "Cannot read note-cli home") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestListBrokenCategoryRegex(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg := &Config{HomePath: filepath.Join(cwd, "testdata", "list", "normal")}
	cmd := &ListCmd{
		Config:   cfg,
		Category: "(foo",
	}

	err = cmd.Do()
	if err == nil {
		t.Fatal("Error did not occur")
	}
	if !strings.Contains(err.Error(), "Regular expression for filtering categories is invalid") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestListBrokenTagRegex(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg := &Config{HomePath: filepath.Join(cwd, "testdata", "list", "normal")}
	cmd := &ListCmd{
		Config: cfg,
		Tag:    "(foo",
	}

	err = cmd.Do()
	if err == nil {
		t.Fatal("Error did not occur")
	}
	if !strings.Contains(err.Error(), "Regular expression for filtering tags is invalid") {
		t.Fatal("Unexpected error:", err)
	}
}

func TestListBrokenNote(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cfg := &Config{HomePath: filepath.Join(cwd, "testdata", "list", "fail")}
	cmd := &ListCmd{Config: cfg}
	err = cmd.Do()
	if err == nil {
		t.Fatal("Error did not occur")
	}
	if !strings.Contains(err.Error(), "Cannot parse created date time") {
		t.Fatal("Unexpected error:", err)
	}
}