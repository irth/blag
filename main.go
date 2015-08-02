package main

import (
	"flag"
	"github.com/flosch/pongo2"
	"path"
)

type Config struct {
	Input  *string
	Theme  *string
	Output *string
}

type BlagPost struct {
	Title     string
	Timestamp int
	Author    string
	Content   string
	Slug      string
}

type Theme struct {
	Page *pongo2.Template
	Post *pongo2.Template
}

func LoadTheme(theme_dir string) Theme {
	t := Theme{}
	t.Page = pongo2.Must(pongo2.FromFile(path.Join(theme_dir, "page.html")))
	t.Post = pongo2.Must(pongo2.FromFile(path.Join(theme_dir, "post.html")))
	return t
}

func main() {
	var config Config
	config.Input = flag.String("input", "./input/", "Directory where blog posts are stored (in markdown format)")
	config.Output = flag.String("output", "./output/", "Directory where generated html should be stored")
	config.Theme = flag.String("theme", "./theme/", "Directory containing theme files (templates)")
	flag.Parse()

	var theme Theme
	theme = LoadTheme(*config.Theme)
}
