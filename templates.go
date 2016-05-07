package blog

import (
	"html/template"
)

// ArticleTemplate is the default article template
var ArticleTemplate, _ = template.New("article").Parse(`<!doctype html><html lang="en"><head> <meta charset="utf-8"> <title>Git based blogging</title> <meta name="description" content="Adam Talbot's code ramblings"> <meta name="author" content="Adam Talbot"> <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.4.1/themes/prism.min.css"> <style>@import url(https://fonts.googleapis.com/css?family=Open+Sans:400,800); html, body{padding: 0; margin: 0; font-family: 'Open Sans', sans-serif;}.header{background: #222; padding: 0.8em 1em; color: #CCC;}.header:after{content:''; display:block; clear:both;}.header__logo{display: inline-block; text-align: center; font-weight: 900; font-family: monospace; font-size: 25px; border: 2px solid #CCCCCC; padding: 2px 5px; margin: 0 0.8em; vertical-align: middle;}.header__title{display: inline-block; vertical-align: middle;}.header__git{display: inline-block; float: right; font-style: italic; font-family: monospace;}.home{display:block; margin: 1em;}.article{border: 2px solid #222; margin: 1em; padding: 1em;}</style></head><body> <header class="header"> <div class="header__logo">B L<br/>O G</div><h1 class="header__title">Git based blogging</h1> <div class="header__git">git clone{{.GitURL}}</div></header> <a class="home" href="{{.BaseURL}}">Home</a> <div class="article"> <p>{{.Article.Full}}</p><i>Posted on{{.Article.Mod}}</i> </div><script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.4.1/prism.min.js"></script> <script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.4.1/components/prism-go.min.js"></script></body></html>`)

// IndexTemplate is the default index template
var IndexTemplate, _ = template.New("index").Parse(`<!doctype html><html lang="en"><head> <meta charset="utf-8"> <title>Git based blogging</title> <meta name="description" content="Adam Talbot's code ramblings"> <meta name="author" content="Adam Talbot"> <style>@import url(https://fonts.googleapis.com/css?family=Open+Sans:400,800); html, body{padding: 0; margin: 0; font-family: 'Open Sans', sans-serif;}.header{background: #222; padding: 0.8em 1em; color: #CCC;}.header:after{content:''; display:block; clear:both;}.header__logo{display: inline-block; text-align: center; font-weight: 900; font-family: monospace; font-size: 25px; border: 2px solid #CCCCCC; padding: 2px 5px; margin: 0 0.8em; vertical-align: middle;}.header__title{display: inline-block; vertical-align: middle;}.header__git{display: inline-block; float: right; font-style: italic; font-family: monospace;}.article{border: 2px solid #222; margin: 1em; padding: 1em;}.pagination{text-align: center;}.pagination a{text-decoration: none;}</style></head><body> <header class="header"> <div class="header__logo">B L<br/>O G</div><h1 class="header__title">Git based blogging</h1> <div class="header__git">git clone{{.GitURL}}</div></header>{{range $article :=.Articles}}<div class="article"> <p>{{$article.Preview $.BaseURL}}</p><i>Posted on{{$article.Mod}}</i> </div>{{end}}<div class="pagination">{{.Pagination}}</div></body></html>`)