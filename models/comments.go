package models

import (
	"github.com/upper/db/v4"
	"time"
)

type Comment struct {
	ID        int       `db:"comment_id,omitempty"`
	CreatedAt time.Time `db:"c_created_at,omitempty"`
	Body      string    `db:"body"`
	PostId    int       `db:"post_id"`
	UserId    int       `db:"user_id"`
	User      `db:",inline"`
}

type CommentsModel struct {
	db db.Session
}

func (m CommentsModel) Table() string {
	return "comments"
}

func (m CommentsModel) GetForPost(postId int) ([]Comment, error) {
	var comments []Comment

	q := m.db.SQL().Select("c.id as comment_id", "c.created_at as c_created_at", "*").
		From(m.Table() + " as c").
		Join("users as u").On("c.user_id = u.id").
		Where(db.Cond{"c.post_id": postId}).
		OrderBy("c.created_at desc")

	err := q.All(&comments)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (m CommentsModel) Insert(body string, postId, userId int) error {
	_, err := m.db.Collection(m.Table()).Insert(map[string]interface{}{
		"body":    body,
		"user_id": userId,
		"post_id": postId,
	})

	if err != nil {
		return err
	}

	return nil
}
