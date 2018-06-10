package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/russross/blackfriday"
	"github.com/termie/go-shutil"
	"gopkg.in/yaml.v2"
)

// Config is a struct that will be used to store Blag's config.
type Config struct {
	Input             *string
	Theme             *string
	Output            *string
	PostsPerPage      *int
	StoryShortLength  *int
	Title             *string
	DateFormat        *string
	BaseURL           *string
	DisqusShortname   *string
	GoogleAnalyticsID *string
	CookieWarning     *bool
}

// BlagPostMeta is a struct that will hold a blogpost metadata
type BlagPostMeta struct {
	Title     string
	Timestamp int64
	Time      string
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
	t.Page = pongo2.Must(pongo2.FromFile(path.Join(themeDir, "templates", "page.html")))
	t.Post = pongo2.Must(pongo2.FromFile(path.Join(themeDir, "templates", "post.html")))
	return t
}

// LoadPost loads post file specified by path argument, and returns BlagPost
// object with data loaded from that file.
func LoadPost(config Config, fpath string) BlagPost {
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

	if meta.Timestamp <= 0 {
		stat, err := file.Stat()
		if err != nil {
			panic(err)
		}
		meta.Timestamp = stat.ModTime().Unix()
	}

	time := time.Unix(meta.Timestamp, 0)
	meta.Time = time.Format(*config.DateFormat)

	if len(meta.Slug) == 0 {
		basename := filepath.Base(file.Name())
		meta.Slug = strings.TrimSuffix(basename, filepath.Ext(basename))
	}

	markdown, _ := ioutil.ReadAll(buf)
	markdown = []byte(strings.Trim(string(markdown), " \r\n"))
	html := string(blackfriday.MarkdownCommon(markdown))
	return BlagPost{
		meta,
		html,
	}
}

// LoadPosts loads all markdown files in inputDir (not recursive), and returns
// a slice []BlagPost, containing extracted metadata and HTML rendered from
// Markdown.
func LoadPosts(config Config) []BlagPost {
	inputDir := *config.Input
	var p []BlagPost
	filelist, err := ioutil.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}
	for _, file := range filelist {
		if strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
			p = append(p, LoadPost(config, path.Join(inputDir, file.Name())))
		}
	}
	return p
}

// GenerateHTML generates page's static html and stores it in directory
// specified in config.
func GenerateHTML(config Config, theme Theme, posts []BlagPost) {
	err := os.RemoveAll(*config.Output)
	if err != nil {
		panic(err)
	}
	shutil.CopyTree(path.Join(*config.Theme, "static"), *config.Output, &shutil.CopyTreeOptions{
		Symlinks:               true,
		IgnoreDanglingSymlinks: true,
		CopyFunction:           shutil.Copy,
		Ignore:                 nil,
	})

	os.MkdirAll(*config.Output, 0755)

	os.MkdirAll(path.Join(*config.Output, "post"), 0755)
	for _, post := range posts {
		postFile, err := os.OpenFile(
			path.Join(*config.Output, "post", fmt.Sprintf("%s.html", post.Slug)),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		defer postFile.Close()
		if err != nil {
			panic(err)
		}
		theme.Post.ExecuteWriter(pongo2.Context{
			"disqus_shortname":    *config.DisqusShortname,
			"google_analytics_id": *config.GoogleAnalyticsID,
			"cookie_warning":      *config.CookieWarning,
			"base":                *config.BaseURL,
			"title":               *config.Title,
			"post":                post,
		}, postFile)
	}

	postCount := len(posts)
	pageCount := int(math.Ceil(float64(postCount) / float64(*config.PostsPerPage)))

	os.MkdirAll(path.Join(*config.Output, "page"), 0755)

	pagePosts := make(map[int][]BlagPost)

	for i := postCount - 1; i >= 0; i-- {
		pageNo := int(math.Floor(float64(postCount-i-1)/float64(*config.PostsPerPage))) + 1
		pagePosts[pageNo] = append(pagePosts[pageNo], posts[i])
	}

	if postCount == 0 {
		pagePosts[1] = make([]BlagPost, 0)
	}

	for k, v := range pagePosts {
		pageFile, err := os.OpenFile(
			path.Join(*config.Output, "page", fmt.Sprintf("%d.html", k)),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		defer pageFile.Close()
		if err != nil {
			panic(err)
		}
		theme.Page.ExecuteWriter(pongo2.Context{
			"disqus_shortname":    *config.DisqusShortname,
			"google_analytics_id": *config.GoogleAnalyticsID,
			"cookie_warning":      *config.CookieWarning,
			"base":                *config.BaseURL,
			"title":               *config.Title,
			"posts":               v,
			"current_page":        k,
			"page_count":          pageCount,
			"shortlen":            *config.StoryShortLength,
		}, pageFile)
		if k == 1 {
			indexFile, err := os.OpenFile(
				path.Join(*config.Output, "index.html"),
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			defer pageFile.Close()
			if err != nil {
				panic(err)
			}
			theme.Page.ExecuteWriter(pongo2.Context{
				"disqus_shortname":    *config.DisqusShortname,
				"google_analytics_id": *config.GoogleAnalyticsID,
				"cookie_warning":      *config.CookieWarning,
				"base":                *config.BaseURL,
				"title":               *config.Title,
				"posts":               v,
				"current_page":        k,
				"page_count":          pageCount,
				"shortlen":            *config.StoryShortLength,
			}, indexFile)
		}
	}
}

func main() {
	var config Config
	config.Input = flag.String("input", "input", "Directory where blog posts are stored (in markdown format)")
	config.Output = flag.String("output", "output", "Directory where generated html should be stored (IT WILL REMOVE ALL FILES INSIDE THAT DIR)")
	config.Theme = flag.String("theme", "theme", "Directory containing theme files (templates)")
	config.Title = flag.String("title", "Blag.", "Blag title")
	config.DateFormat = flag.String("dateformat", "2006-01-02 15:04:05", "Time layout, as used in Golang's time.Time.Format()")
	config.BaseURL = flag.String("baseurl", "/", "URL that will be used in <base href=\"\"> element.")
	config.DisqusShortname = flag.String("disqus", "", "Your Disqus shortname. If empty, comments will be disabled.")
	config.GoogleAnalyticsID = flag.String("google", "", "Your Google Analytics Tracker ID. If empty, analytics will be disabled.")
	config.CookieWarning = flag.Bool("cookies", false, "If enabled and supported by the theme, an EU cookie law warning will be shown.")
	config.PostsPerPage = flag.Int("pps", 10, "Post count per page")
	config.StoryShortLength = flag.Int("short", 250, "Length of shortened versions of stories (-1 disables shortening)")
	flag.Parse()

	pongo2.RegisterFilter("trim", func(in *pongo2.Value, param *pongo2.Value) (out *pongo2.Value, err *pongo2.Error) {
		out = pongo2.AsValue(strings.Trim(in.String(), "\r\n"))
		err = nil
		return out, err
	})

	var theme Theme
	theme = LoadTheme(*config.Theme)

	var posts []BlagPost
	posts = LoadPosts(config)

	GenerateHTML(config, theme, posts)
}
