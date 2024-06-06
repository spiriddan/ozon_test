package inMemory

import (
	"main/graph/model"
	"sync"
	"sync/atomic"
)

type repliesArray = []*model.Comment

type MemoryRepo struct {
	RepliesToPost    *sync.Map // map[int][]*comment.Comment
	RepliesToComment *sync.Map // map[int][]*comment.Comment
	commentsAmount   int64
}

func NewCommentMemoryRepo() *MemoryRepo {
	res := &MemoryRepo{}
	res.RepliesToPost = &sync.Map{}
	res.RepliesToComment = &sync.Map{}
	return res
}

func (repo *MemoryRepo) CreateComment(commen model.CreateCommentInput) (model.Comment, error) {
	comm := model.Comment{}
	atomic.AddInt64(&repo.commentsAmount, 1)
	comm.ID = int(repo.commentsAmount)
	comm.Body = commen.Body
	comm.ParentType = commen.ParentType
	comm.ParentID = commen.ParentID
	if comm.ParentType == model.ParentPost {
		arr, _ := repo.RepliesToPost.Load(comm.ParentID)
		if arr == nil {
			arr = []*model.Comment{&comm}
		} else {
			arr = append(arr.(repliesArray), &comm)
		}
		repo.RepliesToPost.Store(comm.ParentID, arr)
	} else {
		arr, _ := repo.RepliesToComment.Load(comm.ParentID)
		if arr == nil {
			arr = []*model.Comment{&comm}
		} else {
			arr = append(arr.(repliesArray), &comm)
		}
		repo.RepliesToComment.Store(comm.ParentID, arr)
	}

	return comm, nil
}

func (repo *MemoryRepo) GetPostComments(id int) ([]*model.Comment, error) {
	res, _ := repo.RepliesToPost.Load(id)
	if res == nil {
		return nil, nil
	}
	return res.(repliesArray), nil
}

func (repo *MemoryRepo) GetCommentReplies(id int) ([]*model.Comment, error) {
	res, _ := repo.RepliesToComment.Load(id)
	if res == nil {
		return nil, nil
	}
	return res.(repliesArray), nil
}
