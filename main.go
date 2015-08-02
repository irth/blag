package main

import (
	"bufio"
	"flag"
	"github.com/flosch/pongo2"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Config struct {
	Input  *string
	Theme  *string
	Output *string
}

type BlagPostMeta struct {
	Title     string
	Timestamp int
	Author    string
	Slug      string
}

type BlagPost struct {
	BlagPostMeta
	Content string
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

func LoadPost(path string) BlagPost {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(file)
	yaml_meta := ""
	for !strings.HasSuffix(yaml_meta, "\n\n") {
		var s string
		s, err = buf.ReadString('\n')
		yaml_meta += s
	}

	var meta BlagPostMeta
	yaml.Unmarshal([]byte(yaml_meta), &meta)

	markdown, _ := ioutil.ReadAll(buf)
	html := string(blackfriday.MarkdownCommon(markdown))
	return BlagPost{
		meta,
		html,
	}
}

func LoadPosts(input_dir string) []BlagPost {
	var p []BlagPost
	filelist, err := ioutil.ReadDir(input_dir)
	if err != nil {
		panic(err)
	}
	for _, file := range filelist {
		p = append(p, LoadPost(path.Join(input_dir, file.Name())))
	}
	return p
}

func main() {
	var config Config
	config.Input = flag.String("input", "./input/", "Directory where blog posts are stored (in markdown format)")
	config.Output = flag.String("output", "./output/", "Directory where generated html should be stored")
	config.Theme = flag.String("theme", "./theme/", "Directory containing theme files (templates)")
	flag.Parse()

	var theme Theme
	theme = LoadTheme(*config.Theme)

	var posts []BlagPost
	posts = LoadPosts(*config.Input)
}
