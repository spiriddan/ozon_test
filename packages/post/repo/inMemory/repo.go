package inMemory

import (
	"main/graph/model"
	"main/packages/post"
	"sync"
)

type MemoryRepo struct {
	mx    sync.RWMutex
	Posts []*model.Post
}

func NewPostMemoryRepo() *MemoryRepo {
	res := &MemoryRepo{}
	res.Posts = make([]*model.Post, 0)
	return res
}

func (repo *MemoryRepo) GetPost(filter model.PostFilter) (*model.PostPayload, error) {
	repo.mx.RLock()
	defer repo.mx.RUnlock()

	for _, p := range repo.Posts { // TODO
		if p.ID == filter.IDIn {
			return &model.PostPayload{Title: p.Title, Body: p.Body}, nil
		}
	}

	return nil, post.NoPostErr
}

func (repo *MemoryRepo) CreatePost(input model.CreatePostInput) (*model.CreatePostPayload, error) {
	repo.mx.Lock()
	defer repo.mx.Unlock()
	n := len(repo.Posts)
	res := &model.Post{
		ID:         n + 1,
		Title:      input.Title,
		Body:       input.Body,
		CanComment: true,
	}
	repo.Posts = append(repo.Posts, res)
	return &model.CreatePostPayload{Post: res}, nil
}

func (repo *MemoryRepo) GetPosts() (*model.PostsPayload, error) {
	res := &model.PostsPayload{}
	repo.mx.RLock()
	defer repo.mx.RUnlock()
	res.Posts = append(res.Posts, repo.Posts...)
	return res, nil
}
