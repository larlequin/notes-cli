package notes

import (
	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
)

func TestParseArgs(t *testing.T) {
	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(
			CategoriesCmd{},
			ConfigCmd{},
			ListCmd{},
			NewCmd{},
			SaveCmd{},
			TagsCmd{},
		),
		cmpopts.IgnoreTypes(&Config{}),
		cmpopts.IgnoreFields(ListCmd{}, "Out"),
		cmpopts.IgnoreFields(ConfigCmd{}, "Out"),
		cmpopts.IgnoreFields(TagsCmd{}, "Out"),
		cmpopts.IgnoreFields(CategoriesCmd{}, "Out"),
	}

	for _, tc := range []struct {
		args []string
		want Cmd
	}{
		{
			args: []string{"config", "home"},
			want: &ConfigCmd{
				Name: "home",
			},
		},
		{
			args: []string{"save", "--message", "hello"},
			want: &SaveCmd{
				Message: "hello",
			},
		},
		{
			args: []string{"tags", "dog"},
			want: &TagsCmd{
				Category: "dog",
			},
		},
		{
			args: []string{"categories"},
			want: &CategoriesCmd{},
		},
		{
			args: []string{"cats"},
			want: &CategoriesCmd{},
		},
		{
			args: []string{"list", "--category", "dog", "--tag", "cat", "--oneline"},
			want: &ListCmd{
				Category: "dog",
				Tag:      "cat",
				Oneline:  true,
			},
		},
		{
			args: []string{"ls", "--category", "dog", "--tag", "cat", "--oneline"},
			want: &ListCmd{
				Category: "dog",
				Tag:      "cat",
				Oneline:  true,
			},
		},
		{
			args: []string{"new", "dog", "filename", "cat,bird", "--no-inline-input"},
			want: &NewCmd{
				Category: "dog",
				Tags:     "cat,bird",
				NoInline: true,
				Filename: "filename",
			},
		},
	} {
		t.Run(tc.args[0]+" command", func(t *testing.T) {
			have, err := ParseCmd(tc.args)
			if err != nil {
				t.Fatal(err)
			}

			if !cmp.Equal(tc.want, have, opts...) {
				t.Fatal(cmp.Diff(tc.want, have, opts...))
			}
		})
	}
}

func TestParseGlobalFlags(t *testing.T) {
	old := color.NoColor
	defer func() {
		color.NoColor = old
	}()

	if _, err := ParseCmd([]string{"--no-color", "config"}); err != nil {
		t.Fatal(err)
	}

	if !color.NoColor {
		t.Fatal("Color was not disabled")
	}
}

func TestParseFailure(t *testing.T) {
	if _, err := ParseCmd([]string{"unknown-command"}); err == nil {
		t.Fatal("Unknown command did not cause an error")
	}
}