package blog

import (
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/ThatsMrTalbot/scaffold"
	"github.com/ThatsMrTalbot/scaffold/errors"
	"github.com/facebookgo/inject"
	"github.com/gogits/git"
	"github.com/gogits/webdav"
	"golang.org/x/net/context"
)

// App is the base platform
type App struct {
	Config *Config         `inject:""`
	Repo   *git.Repository `inject:""`
	Blog   *Blog           `inject:""`

	git http.Handler
}

// NewApp creates a new app and injects dependencies into graph
func NewApp(config *Config) (*App, error) {
	var graph inject.Graph
	var app App

	repo, err := git.OpenRepository(config.Path)
	if err != nil {
		return nil, err
	}

	templates, err := template.ParseFiles(
		filepath.Join(config.TemplatePath, "article.tpl"),
		filepath.Join(config.TemplatePath, "index.tpl"),
	)

	if err != nil {
		return nil, err
	}

	err = graph.Provide(
		&inject.Object{Value: &app},
		&inject.Object{Value: config},
		&inject.Object{Value: repo},
		&inject.Object{Value: templates},
	)

	if err != nil {
		return nil, err
	}

	if err := graph.Populate(); err != nil {
		return nil, err
	}

	return &app, nil
}

// Error is a scaffold error handler
func (a *App) Error(ctx context.Context, w http.ResponseWriter, r *http.Request, status int, err error) {
	logrus.WithError(err).WithField("status", status).Error("Handler encountered an error")

	w.WriteHeader(status)
	w.Write([]byte("Error!"))
}

// TrailingSlashMiddleware forces urls without extensions to have trailing slashes
func (a *App) TrailingSlashMiddleware(next scaffold.Handler) scaffold.Handler {
	return scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if path.Ext(r.URL.Path) == "" && !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", 302)
		} else {
			next.CtxServeHTTP(ctx, w, r)
		}
	})
}

// GitMiddleware handles webdav request for the git repo
func (a *App) GitMiddleware(next scaffold.Handler) scaffold.Handler {
	return scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/blog.git/") {
			a.git.ServeHTTP(w, r)
		} else {
			next.CtxServeHTTP(ctx, w, r)
		}
	})
}

// RedirectHome redirects home
func (a *App) RedirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 302)
}

// Routes implements scaffold.Platform.Router
func (a *App) Routes(router *scaffold.Router) {
	m := errors.SetErrorHandler(errors.AllStatusCodes, errors.ErrorHandlerFunc(a.Error))

	server := webdav.NewServer(a.Repo.Path, "/blog.git/", true)
	server.ReadOnly = true

	a.git = server

	router.Use(m)
	router.Use(a.TrailingSlashMiddleware)
	router.Handle("blog.git", a.RedirectHome).Use(a.GitMiddleware)
	router.Platform("", a.Blog)
	router.Platform("branch/:branch", a.Blog)
	router.Platform("commit/:commit", a.Blog)
}

// Serve the application
func (a *App) Serve() {
	dispatcher := scaffold.DefaultDispatcher()
	scaffold.Scaffold(dispatcher, a)
	listen := ":80"
	if a.Config.Listen != "" {
		listen = a.Config.Listen
	}

	http.ListenAndServe(listen, dispatcher)
}
