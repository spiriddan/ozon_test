package graph

import "main/packages/post"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	postRepo post.Repo
}

func NewResolver(postRepo post.Repo) *Resolver {
	return &Resolver{
		postRepo: postRepo,
	}
}
