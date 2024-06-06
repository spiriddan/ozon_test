package comment

import "main/graph/model"

type (
	Repo interface {
		CreateComment(comm model.CreateCommentInput) (model.Comment, error)
		GetPostComments(id int) ([]*model.Comment, error)
		GetCommentReplies(id int) ([]*model.Comment, error)
	}
)
