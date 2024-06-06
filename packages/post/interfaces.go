package post

import "main/graph/model"

type (
	Repo interface {
		GetPost(filter model.PostFilter) (*model.PostPayload, error)
		GetPosts() (*model.PostsPayload, error)
		CreatePost(input model.CreatePostInput) (*model.CreatePostPayload, error)
	}

	Subscriber interface {
		AddSubscriber(postId int) (int, chan *model.Comment, error)
		DeleteSubscriber(id int) error
		SendNotification(id int, comment model.Comment) error
	}
)
