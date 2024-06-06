package postgres

import (
	"database/sql"
	qb "github.com/Masterminds/squirrel"
	"main/graph/model"
)

const (
	postFields       = "id, title, body, canComment"
	postInsertFields = "title, body, canComment"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(conn *sql.DB) *PostgresRepo {
	return &PostgresRepo{
		db: conn,
	}
}

func (repo *PostgresRepo) GetPost(filter model.PostFilter) (*model.PostPayload, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	res := &model.PostPayload{}

	// тут можно добавить setMap, когда в фильтре больше полей появится
	query := stBuilder.
		Select(postFields).
		From("post").
		Where(qb.Eq{"id": filter.IDIn}).
		RunWith(repo.db)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tmp := model.Post{}
		err = rows.Scan(&tmp.ID, &tmp.Title, &tmp.Body, &tmp.CanComment)
		if err != nil {
			return nil, err
		}
		res.Posts = append(res.Posts, &tmp)
	}
	return res, nil
}

func (repo *PostgresRepo) GetPosts() (*model.PostsPayload, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)
	res := &model.PostsPayload{}

	query := stBuilder.
		Select(postFields).
		From("post").
		RunWith(repo.db)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tmp := model.Post{}
		err = rows.Scan(&tmp.ID, &tmp.Title, &tmp.Body, &tmp.CanComment)
		if err != nil {
			return nil, err
		}
		res.Posts = append(res.Posts, &tmp)
	}
	return res, nil
}

func (repo *PostgresRepo) CreatePost(input model.CreatePostInput) (*model.CreatePostPayload, error) {
	stBuilder := qb.StatementBuilder.PlaceholderFormat(qb.Dollar)

	query := stBuilder.
		Insert("post").
		Columns(postInsertFields).
		Values(input.Title, input.Body, input.CanComment).
		Suffix("RETURNING id").
		RunWith(repo.db)

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := &model.CreatePostPayload{Post: &model.Post{}}
	rows.Next()
	err = rows.Scan(&res.Post.ID)
	if err != nil {
		return nil, err
	}
	res.Post.Title = input.Title
	res.Post.Body = input.Body
	res.Post.CanComment = input.CanComment
	return res, nil
}
