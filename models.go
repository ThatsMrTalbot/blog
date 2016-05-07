package blog

import (
	"fmt"
	"html/template"
	"net/url"
)

// IndexModel is the model passed to the index template
type IndexModel struct {
	GitURL   string
	Page     int
	Count    int
	Articles []Article
	BaseURL  *url.URL
}

// Pagination creates pagination for the index
func (i *IndexModel) Pagination() template.HTML {
	p := ""
	for e := i.Page - 3; e < i.Page+3; e++ {
		if e < 0 || e >= i.Count {
			continue
		}
		page := fmt.Sprintf(`pages/%d`, e)
		url, _ := i.BaseURL.Parse(page)
		classes := "pagination__page"
		if e == i.Page {
			classes += " pagination__page--active"
		}
		p += fmt.Sprintf(`<a href="%s" class="%s">%d</a> `, url, classes, e)
	}

	if i.Page-3 > 0 {
		url, _ := i.BaseURL.Parse("page/0")
		p = fmt.Sprintf(`<a href="%s">%d</a> ...`, url, 0) + p
	}

	if i.Page+2 < i.Count {
		page := fmt.Sprintf(`pages/%d`, i.Count-1)
		url, _ := i.BaseURL.Parse(page)
		p += fmt.Sprintf(`... <a href="%s">%d</a>`, url, i.Count-1)
	}

	pagination := fmt.Sprintf(`<div class="pagination">%s</div>`, p)
	return template.HTML(pagination)
}

// ArticleModel is the model passed to the article template
type ArticleModel struct {
	GitURL  string
	Article *Article
	BaseURL *url.URL
}
