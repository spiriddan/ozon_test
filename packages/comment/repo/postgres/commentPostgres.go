package postgres

import (
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"main/graph/model"
)

const (
	commentPostFields    = "id, body, parent_post"
	commentCommentFields = "id, body, parent_comment"
	insertPost           = "body, parent_post"
	insertComment        = "body, parent_comment"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(conn *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: conn,
	}
}

func (repo *PostgresRepo) CreateComment(comm model.CreateCommentInput) (model.Comment, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.Insert("comment").RunWith(repo.db)
	var columns string
	if comm.ParentType == model.ParentPost {
		columns = insertPost
	} else {
		columns = insertComment
	}
	query = query.Columns(columns).
		Values(comm.Body, comm.ParentID).
		Suffix("RETURNING id")

	rows, err := query.Query()
	if err != nil {
		return model.Comment{}, err
	}
	defer rows.Close()

	res := model.Comment{}
	rows.Next()
	err = rows.Scan(&res.ID)
	if err != nil {
		return model.Comment{}, err
	}
	res.ParentID = comm.ParentID
	res.Body = comm.Body
	res.ParentType = comm.ParentType
	return res, nil
}

func (repo *PostgresRepo) GetPostComments(id int) ([]*model.Comment, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	query := stBuilder.
		Select(commentPostFields).
		From("comment").
		Where(qb.Eq{"parent_post": id}).
		OrderBy("id").
		RunWith(repo.db)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*model.Comment, 0)
	var scanning *model.Comment
	for rows.Next() {
		scanning = &model.Comment{}
		err = rows.Scan(&scanning.ID, &scanning.Body, &scanning.ParentID)
		if err != nil {
			return nil, err
		}
		scanning.ParentType = model.ParentPost
		res = append(res, scanning)
	}

	return res, nil
}

func (repo *PostgresRepo) GetCommentReplies(id int) ([]*model.Comment, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	query := stBuilder.
		Select(commentCommentFields).
		From("comment").
		Where(qb.Eq{"parent_comment": id}).
		OrderBy("id").
		RunWith(repo.db)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]*model.Comment, 0)
	var scanning *model.Comment
	for rows.Next() {
		scanning = &model.Comment{}
		err = rows.Scan(&scanning.ID, &scanning.Body, &scanning.ParentID)
		if err != nil {
			return nil, err
		}
		scanning.ParentType = model.ParentComment
		res = append(res, scanning)
	}

	return res, nil
}
