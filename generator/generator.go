package generator

import (
	"fmt"
	"github.com/go-delve/delve/pkg/config"
	"html/template"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Meta struct {
	Title 		string
	Short 		string
	Date 		string
	Tags 		[]string
	ParsedDate 	time.Time
}

type IndexData struct {
	HTMLTitle 		string
	PageTitle 		string
	Content 		template.HTML
	Year 			int
	Name 			string
	CanonicalLink 	string
	MetaDescription string
	HighlightCSS 	template.CSS
}

type Generator interface {
	Generate() error
}

type SiteGenerator struct {
	Config *SiteConfig
}

type SiteConfig struct {
	Sources []string
	Destination string
	Config *config.Config
}

func (g *SiteGenerator) Generate() error {
	templatePath := filepath.Join("static", "template.html")
	fmt.Printf("Generating Site...")
	sources := g.Config.Sources
	destination := g.Config.Destination
	if err := clearAndCreateDestination(destination); err != nil {
		return err
	}
	if err := clearAndCreateDestination(filepath.Join(destination, "archive")); err != nil {
		return nil
	}
	t, err := getTemlate(templatePath)
	if err != nil {
		return err
	}
	var posts []*Posts
	for _, path := range sources {
		post, err := newPost(path, g.Config.Config.Blog.DateFormat)
		if err != nil {
			return err
		}
		posts = append(posts, post)
	}
	sort.Sort(ByDateDesc(posts))
	if err := runTasks(posts, t, destination, g.Config.Config); err != nil {
		return err
	}
	fmt.Print("Finished generating Site...")
	return nil
}

func runTasks(posts []*Posts, t *template.Template, destination string, cfg *config.Config) error {
	var wg sync.WaitGroup
	finished := make(chan bool, 1)
	errors := make(chan error, 1)
	pool := make(chan struct{}, 50)
	generators := []Generator{}

	indexWriter := &IndexWriter{
		BlogURL: cfg.Blog.URL,
		BlogTitle: cfg.Blog.Title,
		BlogDescription: cfg.Blog.Description,
		BlogAuthor: cfg.Blog.Author,
	}

	for _, post := range posts {
		pg := PostGenerator{&PostConfig{
			Post: post,
			Destination: destination,
			Template: t,
			Writer: indexWriter,
		}}
		generators = append(generators, &pg)
	}
	tagPostsMap := createTagPostsMap(posts)
	// frongpage
	fg := ListingGenerator{&ListingConfig{
		Posts:			posts[:getNumOfPagesOnFrontpage(posts, cfg.Blog.Frontpageposts)],
		Template:		t,
		Destination:	destination,
		PageTitle:		"",
		IsIndex:		true,
		Writer:			indexWriter,
	}}
	// archive
	ag := ListingGenerator{&ListingConfig{
		Posts:			posts,
		Template:		t,
		Destination:	path.Join(destination, "archive"),
		PageTitle:		"Archive",
		IsIndex:		false,
		Writer:			indexWriter,
	}}
	// tags
	tg := TagsGenerator{&TagsConfig{
		TagPostsMap:	tagPostsMap,
		Template:		t,
		Destination:	destination,
		Writer:			indexWriter,
	}}

	staticURLs := []string{}
	for _, staticURL := range cfg.Blog.Statics.Templates {
		staticURLs = append(staticURLs, staticURL.Dest)
	}
}