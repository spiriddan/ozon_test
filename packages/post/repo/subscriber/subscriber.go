package subscriber

import (
	"main/graph/model"
	"main/packages/post"
	"sync"
	"sync/atomic"
)

type CommentSubscriberManager struct {
	subscribers *sync.Map // [int][]Subscriber
	counter     int64
}

type Subscriber struct {
	id int64
	ch chan *model.Comment
}

func NewCommentSubscriberManager() *CommentSubscriberManager {
	return &CommentSubscriberManager{
		subscribers: &sync.Map{},
	}
}

func (sub *CommentSubscriberManager) AddSubscriber(postId int) (int, chan *model.Comment, error) {
	resChan := make(chan *model.Comment)
	atomic.AddInt64(&sub.counter, 1)
	id := sub.counter
	subs, _ := sub.subscribers.Load(postId)
	if subs != nil {
		subs = append(subs.([]Subscriber), Subscriber{id: id, ch: resChan})
	} else {
		subs = []Subscriber{{id: id, ch: resChan}}
	}
	sub.subscribers.Store(id, subs)
	return int(id), resChan, nil
}

func (sub *CommentSubscriberManager) DeleteSubscriber(id int) error {
	sub.subscribers.Range(func(key, value interface{}) bool {
		for i, k := range value.([]Subscriber) {
			if int(k.id) == id {
				subs, _ := sub.subscribers.Load(key)
				subs = append(subs.([]Subscriber)[:i], subs.([]Subscriber)[:i+1]...)
				return true
			}
		}
		return true
	})
	return nil
}

func (sub *CommentSubscriberManager) SendNotification(id int, comment model.Comment) error {
	chans, ok := sub.subscribers.Load(id)
	if !ok {
		return post.NoPostError
	}
	chs := chans.([]Subscriber)

	for _, ch := range chs {
		ch.ch <- &comment
	}
	return nil
}
