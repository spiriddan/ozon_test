package comment

import "fmt"

var NoPostError = fmt.Errorf("no post with such ID")
var NoCommentError = fmt.Errorf("no comment with such ID")
var CommentNotAllowed = fmt.Errorf("can't comment this post")
