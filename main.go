package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/russross/blackfriday"
	"gopkg.in/yaml.v2"
)

// Config is a struct that will be used to store Blag's config.
type Config struct {
	Input            *string
	Theme            *string
	Output           *string
	PostsPerPage     *int
	StoryShortLength *int
}

// BlagPostMeta is a struct that will hold a blogpost metadata
type BlagPostMeta struct {
	Title     string
	Timestamp int
	Author    string
	Slug      string
}

// BlagPost is a struct that holds post's content (in html) and its metadata
type BlagPost struct {
	BlagPostMeta
	Content string
}

// Theme holds templates that will be used to render HTML
type Theme struct {
	Page *pongo2.Template
	Post *pongo2.Template
}

// LoadTheme loads pongo2 templates for both pages and posts.
// It will try to load templates from themeDir/page.html and
// themeDir/post.html, and it will panic if that will not succeed.
func LoadTheme(themeDir string) Theme {
	t := Theme{}
	t.Page = pongo2.Must(pongo2.FromFile(path.Join(themeDir, "page.html")))
	t.Post = pongo2.Must(pongo2.FromFile(path.Join(themeDir, "post.html")))
	return t
}

// LoadPost loads post file specified by path argument, and returns BlagPost
// object with data loaded from that file.
func LoadPost(fpath string) BlagPost {
	file, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	buf := bufio.NewReader(file)
	yamlMeta := ""
	for !strings.HasSuffix(yamlMeta, "\n\n") {
		var s string
		s, err = buf.ReadString('\n')
		yamlMeta += s
	}

	var meta BlagPostMeta
	yaml.Unmarshal([]byte(yamlMeta), &meta)

	if len(meta.Slug) == 0 {
		basename := filepath.Base(file.Name())
		meta.Slug = strings.TrimSuffix(basename, filepath.Ext(basename))
	}

	markdown, _ := ioutil.ReadAll(buf)
	html := string(blackfriday.MarkdownCommon(markdown))
	return BlagPost{
		meta,
		html,
	}
}

// LoadPosts loads all markdown files in inputDir (not recursive), and returns
// a slice []BlagPost, containing extracted metadata and HTML rendered from
// Markdown.
func LoadPosts(inputDir string) []BlagPost {
	var p []BlagPost
	filelist, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	for _, file := range filelist {
		p = append(p, LoadPost(path.Join(inputDir, file.Name())))
	}
	return p
}

// GenerateHTML generates page's static html and stores it in directory
// specified in config.
func GenerateHTML(config Config, theme Theme, posts []BlagPost) {
	for _, post := range posts {
		postFile, err := os.OpenFile(path.Join(*config.Output, post.Slug+".html"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		defer postFile.Close()
		if err != nil {
			panic(err)
		}
		theme.Post.ExecuteWriter(pongo2.Context{
			"title":     post.Title,
			"author":    post.Author,
			"timestamp": post.Timestamp,
			"content":   post.Content,
		}, postFile)
	}
}

func main() {
	var config Config
	config.Input = flag.String("input", "input", "Directory where blog posts are stored (in markdown format)")
	config.Output = flag.String("output", "output", "Directory where generated html should be stored")
	config.Theme = flag.String("theme", "theme", "Directory containing theme files (templates)")
	config.PostsPerPage = flag.Int("pps", 10, "Post count per page")
	config.StoryShortLength = flag.Int("short", 250, "Length of shortened versions of stories (-1 disables shortening)")
	flag.Parse()

	var theme Theme
	theme = LoadTheme(*config.Theme)

	var posts []BlagPost
	posts = LoadPosts(*config.Input)

	GenerateHTML(config, theme, posts)
}
