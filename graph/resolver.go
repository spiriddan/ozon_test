package graph

import (
	"main/packages/comment"
	"main/packages/post"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	postRepo          post.Repo
	subscriberManager post.Subscriber
	commentRepo       comment.Repo
}

func NewResolver(postRepo post.Repo, subManager post.Subscriber, commentRepo comment.Repo) *Resolver {
	return &Resolver{
		postRepo:          postRepo,
		subscriberManager: subManager,
		commentRepo:       commentRepo,
	}
}
