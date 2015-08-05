## Welcome to Blag.
`Blag` (see [xkcd://148](https://xkcd.com/148/)) is yet another simple static page (blog in this case) generator.

## Demo!
[Click here to read some posts by the famous bestseller writer, Markov Chain.](https://irth.pl/blagdemo/)  
Command line used to generate it:

	blag -baseurl="https://irth.pl/blagdemo/" -input=blagdemo-in -output=blagdemo -title="Blag demo." -disqus="blagdemo" -theme="$GOPATH/src/github.com/irth/blag/theme" -short=500 -pps=2

## Why?
Because a simple static page generator isn't that hard to write, yet not easy enough to be boring to make. I wanted something to code, and this seemed like a nice idea. So, for fun. Like [gopher Hacker News proxy](https://github.com/irth/gophernews) (yep, Internet Gopher, that one with menus.)
Maybe I'll write a blog. Just kidding. I don't have anything wise to say.

## Ok. How do I use it?

Install it:

	go get github.com/irth/blag

then check out command line arguments (config is done only via cmdline args):

	blag --help

Argument name | Default               | Description
--------------|-----------------------|------------
baseurl       | "/"                   | Base URL of your website. It will be used in `<base href="{{ baseurl }}">`.
dateformat    | "2006-01-02 15:04:05" | Format of date. See [Golang docs](http://golang.org/pkg/time/#Time.Format).
disqus        | ""                    | If you want to use Disqus comments, set this to your shortname. If you don't leave it with default value, which is an empty string.
input         | "input"               | This is path to the directory where your `blagposts` are. They will be read in alphabetical order and appear on page in reverse alphabetical order. So, your post names should probably look similar to `01-firstpost.md` or `2015-05-12_20:22_whatever.md`. (`.md` extension is obligatory)
output        | "output"              | This is the directory where `blag` will store generated static HTML. **It will be deleted recursively (using `os.RemoveAll`). So, be careful.**
pps           | 10                    | Number of posts per page. Positive integer, please.
short         | 250                   | Length to which articles may be shortened when compact version is needed.
theme         | "theme"               | Directory containing the desired theme. Default simple theme probably exists in `$GOPATH/src/github.com/irth/blag/theme`. Feel free to use it.
title         | "Blag."               | Title of your `blag`.


Store your posts in your input directory, with names ending with `.md`.  
Post files should look like that:

	title: "Post title"
	timestamp: 112342151
	slug: "this_will_be_in_article_url"
	author: "John Doe"

	markdownmarkdownmarkdown **MARKDOWN**

	even more markdown

Top part is YAML. Bottom part is Markdown.  
If you omit `title` or `author`, a dragon will come and burn your house.  
If you omit `timestamp`, file modification date will be used.  
If you omit `slug`, file name without extension will be used.


If you want to write your own theme, make a directory looking like that

	your_directory/
	 |- static/
	 |- templates/
	     |- page.html
	     |- post.html

Files from `static` directory will be copied as they are to output directory.  
Templates use pongo2 syntax, which should be compatible with Django's template syntax.  

`templates/page.html` is template used to display multiple posts. First page will be also copied to `index.html`

Variables available in that template:

	"disqus_shortname": *config.DisqusShortname,  // --disqus cmdline argument value
	"base":             *config.BaseURL,          // --baseurl
	"title":            *config.Title,            // --title
	"posts":            v,                        // slice (list) of posts.
	                                              // You can iterate over it.
	"current_page":     k,                        // Current page number.
	                                              // Do something cool with it, like "next/previous page" links.
	"page_count":       pageCount,                // Count of all pages.
	"shortlen":         *config.StoryShortLength, // --short

`templates/post.html` is template used to generate single post page.
Variables:

	"disqus_shortname": *config.DisqusShortname, // --disqus cmdline argument value
	"base":             *config.BaseURL,         // --baseurl
	"title":            *config.Title,           // --title
	"post":             post,                    // post struct


Posts in `post` and `posts` look like that:

	Title     string
	Timestamp int64
	Time      string // date and time formatted according to --dateformat
	Author    string
	Slug      string
    Content   string // html rendered from markdown

You can refer to them in templates using `.`, e.g. `post.Title`.  
Definitely have a look at the template in this repo, it's pretty short when you don't count CSS, and you don't need to read it, and it allows you to see template inheritance, `for` loops, `if` statements, using those variables, etc.


Have fun.
