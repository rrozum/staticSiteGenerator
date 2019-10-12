package generator

import (
	"fmt"
	"github.com/go-delve/delve/pkg/config"
	"html/template"
	"path/filepath"
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
	//templatePath := filepath.Join("static", "template.html")
	//fmt.Printf()
}