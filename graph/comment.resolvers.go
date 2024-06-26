package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.47

import (
	"context"
	"fmt"
	"main/graph/model"
	"main/packages/comment"
	"main/packages/post"
)

// Replies is the resolver for the replies field.
func (r *commentResolver) Replies(ctx context.Context, obj *model.Comment) ([]*model.Comment, error) {
	return r.commentRepo.GetCommentReplies(obj.ID)
}

// CreateComment is the resolver for the createComment field.
func (r *mutationResolver) CreateComment(ctx context.Context, input model.CreateCommentInput) (*model.CreateCommentPayload, error) {
	if len(input.Body) > 2000 {
		return nil, fmt.Errorf("too long comment (max length 2000 symbols)")
	}
	if input.ParentType == model.ParentPost {
		posts, err := r.postRepo.GetPost(model.PostFilter{IDIn: input.ParentID})
		if err != nil {
			return nil, err
		}
		if len(posts.Posts) == 0 {
			return nil, post.NoPostError
		}
		if !posts.Posts[0].CanComment {
			return nil, comment.CommentNotAllowed
		}
	}
	res, err := r.commentRepo.CreateComment(input)
	if err != nil {
		return nil, err
	}
	return &model.CreateCommentPayload{Comment: &res}, nil
}

// CommentsSubscription is the resolver for the CommentsSubscription field.
func (r *subscriptionResolver) CommentsSubscription(ctx context.Context, input model.SubsInput) (<-chan *model.Comment, error) {
	id, ch, err := r.subscriberManager.AddSubscriber(input.ID)
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		err = r.subscriberManager.DeleteSubscriber(id)
	}()

	return ch, nil
}

// Comment returns CommentResolver implementation.
func (r *Resolver) Comment() CommentResolver { return &commentResolver{r} }

type commentResolver struct{ *Resolver }
