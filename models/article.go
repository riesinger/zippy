package models

import (
	"time"
)

type Article struct {
	// Title is the Articles plain title, e.g. "Hello World!"
	Title string `bson:"title" json:"title"`
	// Path is the path relative to the domain name that points to the article.
	// If the article would be at http://example.org/blog/posts/technology/awesome-tech.html,
	// Path would be /blog/posts/technology/, Note the leading and trailing slash
	Path string `bson:"path" json:"path"`
	// Slug is the article's title slug (with counter). So for the title "Hello World!", the slug
	// would be "hello-world" (or "hello-world-1"), if a post with that slug alreay exists.
	Slug string `bson:"slug" json:slug`
	// Posts marked as draft cannot be accessed from the web
	IsDraft bool `bson:"isDraft" json:isDraft`
	// MarkdownBody is the actual markdown, that the article is made of
	MarkdownBody string `bson:"mdBody" json:"markdownBody"`
	// HtmlBody is the compiled html for the given markdown
	HtmlBody string `bson:"htmlBody" json:"htmlBody"`
	// Template denotes, which template the article should be rendered in, e.g. "default" or "post".
	// The available templates are given by the currently used theme, but there should always be the
	// "default" template.
	Template string `bson:"tmpl" json:"template"`
	// AuthorID is the ID of the article's writer. Note, that this does not follow the database ID
	// for portability.
	AuthorID  string    `bson:"authorID" json:"authorID"`
	CreatedAt time.Time `bson:"created" json:"createdAt"`
	UpdatedAt time.Time `bson:"updated" json:"updatedAt"`
}
