package blog

import (
	"html/template"
	"io"
	"io/ioutil"
	"sort"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/gogits/git"
	"github.com/russross/blackfriday"
)

type node struct {
	IndexTemplate   *template.Template
	ArticleTemplate *template.Template

	Index    Index
	Articles map[string]*Article
	Tree     *git.Tree
}

// Cache gets and caches file trees and articles
type Cache struct {
	Repo *git.Repository `inject:""`

	lock  sync.RWMutex
	cache map[string]node
}

// Clear clears the cache
func (c *Cache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cache = make(map[string]node)
}

// ClearOne clears the cache for one commit id
func (c *Cache) ClearOne(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.cache == nil {
		return
	}

	delete(c.cache, id)
}

// GetIndex gets an Index from tree and commit ids
func (c *Cache) GetIndex(tid string, id string) (Index, bool) {
	if c.exists(id) {
		return c.getIndex(id)
	}

	if c.Build(tid, id) {
		return c.getIndex(id)
	}

	return nil, false
}

func (c *Cache) getIndex(id string) (Index, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return nil, false
	}

	if n, ok := c.cache[id]; ok {
		return n.Index, true
	}

	return nil, false
}

// GetFile gets a file from tree and commit ids
func (c *Cache) GetFile(tid string, id string, path string) (io.Reader, bool) {
	if c.exists(id) {
		return c.getFile(id, path)
	}

	if c.Build(tid, id) {
		return c.getFile(id, path)
	}

	return nil, false
}

func (c *Cache) getFile(id string, path string) (io.Reader, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return nil, false
	}

	if n, ok := c.cache[id]; ok {
		blob, err := n.Tree.GetBlobByPath(path)
		if err != nil {
			return nil, false
		}

		r, err := blob.Data()
		if err != nil {
			return nil, false
		}

		return r, true
	}

	return nil, false
}

// GetIndexTemplate gets the template from the cache
func (c *Cache) GetIndexTemplate(tid string, id string) *template.Template {
	if c.exists(id) {
		if tpl, ok := c.getIndexTemplate(id); ok {
			return tpl
		}
	} else if c.Build(tid, id) {
		if tpl, ok := c.getIndexTemplate(id); ok {
			return tpl
		}
	}

	return IndexTemplate
}

func (c *Cache) getIndexTemplate(id string) (*template.Template, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return nil, false
	}

	if n, ok := c.cache[id]; ok {
		return n.IndexTemplate, true
	}

	return nil, false
}

// GetArticleTemplate gets the template from the cache
func (c *Cache) GetArticleTemplate(tid string, id string) *template.Template {
	if c.exists(id) {
		if tpl, ok := c.getArticleTemplate(id); ok {
			return tpl
		}
	} else if c.Build(tid, id) {
		if tpl, ok := c.getArticleTemplate(id); ok {
			return tpl
		}
	}

	return ArticleTemplate
}

func (c *Cache) getArticleTemplate(id string) (*template.Template, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return nil, false
	}

	if n, ok := c.cache[id]; ok {
		return n.ArticleTemplate, true
	}

	return nil, false
}

func (c *Cache) buildIndexTemplate(tree *git.Tree) *template.Template {
	blob, err := tree.GetBlobByPath("index.tpl")
	if err != nil {
		logrus.WithError(err).Error("Could not read template blob")
		return IndexTemplate
	}

	reader, err := blob.Data()
	if err != nil {
		logrus.WithError(err).Error("Could not read template blob data")
		return IndexTemplate
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		logrus.WithError(err).Error("Could not read template data")
		return IndexTemplate
	}

	tpl, err := template.New("index").Parse(string(bytes))
	if err != nil {
		logrus.WithError(err).Error("Could parse template")
		return IndexTemplate
	}

	return tpl
}

func (c *Cache) buildArticleTemplate(tree *git.Tree) *template.Template {
	blob, err := tree.GetBlobByPath("article.tpl")
	if err != nil {
		logrus.WithError(err).Error("Could not read template blob")
		return ArticleTemplate
	}

	reader, err := blob.Data()
	if err != nil {
		logrus.WithError(err).Error("Could not read template blob data")
		return ArticleTemplate
	}

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		logrus.WithError(err).Error("Could not read template data")
		return ArticleTemplate
	}

	tpl, err := template.New("index").Parse(string(bytes))
	if err != nil {
		logrus.WithError(err).Error("Could parse template")
		return ArticleTemplate
	}

	return tpl
}

// GetArticle gets an article from tree and commit ids
func (c *Cache) GetArticle(tid string, id string, article string) (*Article, bool) {
	if c.exists(id) {
		return c.getArticle(id, article)
	}

	if c.Build(tid, id) {
		return c.getArticle(id, article)
	}

	return nil, false
}

func (c *Cache) getArticle(id string, article string) (*Article, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return nil, false
	}

	if n, ok := c.cache[id]; ok {
		if article, ok := n.Articles[article]; ok {
			return article, true
		}
	}

	return nil, false
}

func (c *Cache) exists(id string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.cache == nil {
		return false
	}

	_, ok := c.cache[id]

	return ok
}

// Build gets and caches information on a tree and commit id combo
func (c *Cache) Build(tid string, id string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	logrus.WithField("commit", id).WithField("tree", tid).Info("Building cache")

	if c.cache == nil {
		c.cache = make(map[string]node)
	}

	sha1, err := git.NewIdFromString(tid)
	if err != nil {
		logrus.WithError(err).WithField("tree", tid).Error("Could not build data")
		return false
	}

	tree := git.NewTree(c.Repo, sha1)

	n := node{
		Articles: make(map[string]*Article),
		Tree:     tree,
	}

	scanner, err := tree.Scanner()
	if err != nil {
		logrus.WithError(err).WithField("tree", tid).Error("Could not build data")
		return false
	}

	for scanner.Scan() {
		entry := scanner.TreeEntry()

		name := entry.Name()

		if entry.IsDir() {
			logrus.
				WithField("tree", tid).
				WithField("directory", name).
				Info("Directory ignored")

			continue
		}

		if len(name) <= 3 || name[len(name)-3:] != ".md" {
			logrus.
				WithField("tree", tid).
				WithField("filename", name).
				Info("Non markdown file ignored")

			continue
		}

		reader, err := entry.Blob().Data()
		if err != nil {
			logrus.
				WithError(err).
				WithField("tree", tid).
				WithField("filename", name).
				Warn("File blob could not be generated")

			continue
		}

		markdown, err := ioutil.ReadAll(reader)
		if err != nil {
			logrus.
				WithError(err).
				WithField("tree", tid).
				WithField("filename", name).
				Warn("File blob could not be read")

			continue
		}

		commit, err := c.Repo.GetCommit(id)
		if err != nil {
			logrus.
				WithError(err).
				WithField("commit", id).
				Warn("Could not get commit")

			continue
		}

		fileCommit, err := commit.GetCommitOfRelPath(name)
		if err != nil {
			logrus.
				WithError(err).
				WithField("commit", id).
				WithField("filename", name).
				Warn("Could not get relative commit")

			continue
		}

		if fileCommit.Committer == nil {
			logrus.
				WithError(err).
				WithField("commit", fileCommit.Id.String()).
				Warn("Committer information not set")

			continue
		}

		article := Article{
			Name: name[:len(name)-3],
			Mod:  fileCommit.Committer.When,
			Data: blackfriday.MarkdownCommon(markdown),
		}

		logrus.
			WithField("commit", id).
			WithField("tree", tid).
			WithField("article", article.Name).
			Info("Article cached")

		n.Index = append(n.Index, article)
		n.Articles[article.Name] = &article
	}

	sort.Sort(n.Index)

	n.IndexTemplate = c.buildIndexTemplate(tree)
	n.ArticleTemplate = c.buildArticleTemplate(tree)

	c.cache[id] = n

	logrus.WithField("commit", id).WithField("tree", tid).Info("Cache built")

	return true
}
