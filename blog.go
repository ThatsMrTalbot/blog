package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/ThatsMrTalbot/scaffold"
	"github.com/ThatsMrTalbot/scaffold/errors"
	"github.com/gogits/git"
	"golang.org/x/net/context"
)

// Blog is the blog platfomr
type Blog struct {
	Repo      *git.Repository    `inject:""`
	Config    *Config            `inject:""`
	Cache     *Cache             `inject:""`
	Templates *template.Template `inject:""`
}

// IndexModel creates an IndexModel for use in the index template
func (b *Blog) IndexModel(ctx context.Context, r *http.Request, index Index) *IndexModel {
	page, _ := scaffold.GetParam(ctx, "page").Int()

	return &IndexModel{
		Logo:     template.HTML(b.Config.Logo),
		Title:    template.HTML(b.Config.Title),
		Page:     page,
		Count:    index.Pages(20),
		Articles: index.Page(page, 20),
		BaseURL:  b.BaseURL(ctx),
		GitURL:   b.GitURL(r),
	}
}

// ArticleModel creates an ArticleModel for use in the article tempalte
func (b *Blog) ArticleModel(ctx context.Context, r *http.Request, article *Article) *ArticleModel {
	return &ArticleModel{
		Logo:    template.HTML(b.Config.Logo),
		Title:   template.HTML(b.Config.Title),
		Article: article,
		BaseURL: b.BaseURL(ctx),
		GitURL:  b.GitURL(r),
	}
}

// Index is the index handler
func (b *Blog) Index(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	logrus.
		WithField("http_request", r).
		Info("Index handler called")

	tid, id, err := getID(ctx, b.Repo)
	if err != nil {
		return errors.ConvertErrorStatus(500, err)
	}

	index, ok := b.Cache.GetIndex(tid, id)
	if !ok {
		return errors.NewErrorStatus(404, "Index not found")
	}

	logrus.
		WithField("http_request", r).
		Info("Loaded index from cache")

	model := b.IndexModel(ctx, r, index)

	var buffer bytes.Buffer
	err = b.Templates.ExecuteTemplate(&buffer, "index.tpl", model)
	if err != nil {
		return errors.ConvertErrorStatus(500, err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(buffer.Bytes())

	return nil
}

// Article is the article handler
func (b *Blog) Article(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	logrus.
		WithField("http_request", r).
		Info("Article handler called")

	tid, id, err := getID(ctx, b.Repo)
	if err != nil {
		return errors.ConvertErrorStatus(500, err)
	}

	name, _ := scaffold.GetParam(ctx, "article").String()

	article, ok := b.Cache.GetArticle(tid, id, name)
	if !ok {
		return errors.NewErrorStatus(404, "Article not found")
	}

	logrus.
		WithField("http_request", r).
		Info("Loaded article from cache")

	model := b.ArticleModel(ctx, r, article)

	var buffer bytes.Buffer
	err = b.Templates.ExecuteTemplate(&buffer, "article.tpl", model)
	if err != nil {
		return errors.ConvertErrorStatus(500, err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(buffer.Bytes())

	return nil
}

// BaseURL calculates the base url, for example / or /branch/master/
func (b *Blog) BaseURL(ctx context.Context) *url.URL {
	base := "/"
	if branch := scaffold.GetParam(ctx, "branch"); branch != "" {
		base = fmt.Sprintf("/branch/%s/", branch)
	}
	if commit := scaffold.GetParam(ctx, "commit"); commit != "" {
		base = fmt.Sprintf("/commit/%s/", commit)
	}
	baseURL, _ := url.Parse(base)
	return baseURL
}

// GitURL calculates the git url based on the request host
func (b *Blog) GitURL(r *http.Request) string {
	return (&url.URL{
		Scheme: "http",
		Host:   r.Host,
		Path:   "/blog.git",
	}).String()
}

// FileLoaderMiddleware is used to load files from the tree
// Middleware is used as it does not require a specific path
func (b *Blog) FileLoaderMiddleware(next scaffold.Handler) scaffold.Handler {
	return scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		article, _ := scaffold.GetParam(ctx, "article").String()

		baseURL := b.BaseURL(ctx)
		basePath := path.Join(baseURL.Path, "article", article) + "/"
		path := strings.TrimPrefix(r.URL.Path, basePath)
		tid, id, err := getID(ctx, b.Repo)

		if err == nil {
			if reader, ok := b.Cache.GetFile(tid, id, path); ok {
				logrus.
					WithField("http_request", r).
					WithField("filepath", path).
					Info("Serving file")

				if ctype := mime.TypeByExtension(filepath.Ext(r.URL.Path)); ctype != "" {
					w.Header().Set("Content-Type", ctype)
				}
				io.Copy(w, reader)
			} else {
				next.CtxServeHTTP(ctx, w, r)
			}
		}
	})
}

// Routes implements scaffold.Platform.Routes
func (b *Blog) Routes(router *scaffold.Router) {
	router.AddHandlerBuilder(errors.HandlerBuilder)

	router.Get("/", b.Index)
	router.Get("page/:page", b.Index)
	router.Get("article/:article", b.Article)
	router.Get("article/:article/:file").Use(b.FileLoaderMiddleware)
}

// Get tree and commit id based off current request
func getID(ctx context.Context, repo *git.Repository) (string, string, error) {
	if branch, err := scaffold.GetParam(ctx, "branch").String(); branch != "" && err == nil {
		c, err := repo.GetCommitOfBranch(branch)
		if err != nil {
			return "", "", err
		}

		return c.TreeId().String(), c.Id.String(), nil
	}
	if commit, err := scaffold.GetParam(ctx, "commit").String(); commit != "" && err == nil {
		c, err := repo.GetCommit(commit)
		if err != nil {
			return "", "", err
		}

		return c.TreeId().String(), c.Id.String(), nil
	}

	c, err := repo.GetCommitOfBranch("master")
	if err != nil {
		return "", "", err
	}

	return c.TreeId().String(), c.Id.String(), nil
}
