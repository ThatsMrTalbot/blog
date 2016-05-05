package blog

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
)

// Article contains information on an article
type Article struct {
	Name string
	Mod  time.Time
	Data []byte
}

// Preview generates a preview for the article listing
func (a *Article) Preview(baseURL *url.URL) template.HTML {
	data, _ := sanitize.HTMLAllowing(string(a.Data), []string{"h1", "a"})
	data = strings.Replace(data, "<h1>", "<h3>", -1)
	data = strings.Replace(data, "</h1>", "</h3>", -1)
	if len(data) > 500 {
		url, _ := baseURL.Parse("article/" + a.Name)
		data = fmt.Sprintf(`%s... <a href="%s">(Read more)</a>`, data[:500], url.String())
	}
	return template.HTML(data)
}

// Full returns the full article
func (a *Article) Full() template.HTML {
	return template.HTML(a.Data)
}
