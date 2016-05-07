package blog

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

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

	if err != nil {
		return nil, err
	}

	err = graph.Provide(
		&inject.Object{Value: &app},
		&inject.Object{Value: config},
		&inject.Object{Value: repo},
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
	log := GetLog(ctx)

	log.
		WithError(err).
		WithField("status", status).
		Error("Handler encountered an error")

	w.WriteHeader(status)
	w.Write([]byte("Error!"))
}

// NotFound is a not found handler
func (a *App) NotFound(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("Page not found: %s", r.URL.Path)
	errors.GetErrorHandler(ctx, 404).ServeErrorPage(ctx, w, r, 404, err)
}

// TrailingSlashMiddleware forces urls without extensions to have trailing slashes
func (a *App) TrailingSlashMiddleware(next scaffold.Handler) scaffold.Handler {
	return scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		if path.Ext(r.URL.Path) == "" && !strings.HasSuffix(r.URL.Path, "/") {
			log := GetLog(ctx)
			log.Info("Redirecting to URL with slash appended")

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
			log := GetLog(ctx)
			log.Info("Passed request to WebDav")

			a.git.ServeHTTP(w, r)
		} else {
			next.CtxServeHTTP(ctx, w, r)
		}
	})
}

// LogMetricsMiddleware logs metrics for a request
func (a *App) LogMetricsMiddleware(next scaffold.Handler) scaffold.Handler {
	return scaffold.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log := GetLog(ctx)

		start := time.Now()

		rid := RequestID(r)

		log = log.WithField("start", start)
		log = log.WithField("request_id", rid)
		log = log.WithField("request_headers", r.Header)
		log = log.WithField("request_method", r.Method)
		log = log.WithField("request_cookies", r.Cookies())
		log = log.WithField("url", r.URL.String())

		ctx = StoreLog(ctx, log)

		log.Info("Serving request")

		next.CtxServeHTTP(ctx, w, r)

		duration := time.Since(start)
		log.
			WithField("duration", duration).
			Info("Request served")
	})
}

// RedirectHome redirects home
func (a *App) RedirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 302)
}

// Routes implements scaffold.Platform.Router
func (a *App) Routes(router *scaffold.Router) {
	// WebDav
	server := webdav.NewServer(a.Repo.Path, "/blog.git/", true)
	server.ReadOnly = true
	a.git = server
	router.Handle("blog.git", a.RedirectHome).Use(a.GitMiddleware)

	// Error handlers
	errorHandlerMiddleware := errors.SetErrorHandler(errors.AllStatusCodes, errors.ErrorHandlerFunc(a.Error))
	router.Use(errorHandlerMiddleware)
	router.NotFound(a.NotFound)

	// Trailing slash middleware
	router.Use(a.TrailingSlashMiddleware)

	// Metric logging middleware
	router.Use(a.LogMetricsMiddleware)

	// App routes
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
