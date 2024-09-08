package models

import (
	"errors"
	"github.com/upper/db/v4"
	"strings"
	"time"
)

var (
	ErrDuplicateTitle = errors.New("title already exist in database")
	ErrDuplicateVotes = errors.New("you already voted")

	queryTemplate = `
	SELECT COUNT(*) OVER() AS total_records, pg.*, u.name as uname FROM (
	    SELECT p.id, p.title, p.url, p.created_at, p.user_id as uid, COUNT(c.post_id) as comment_count, COUNT(v.post_id) as votes
	    FROM posts p 
	    LEFT JOIN comments c ON p.id = c.post_id
	    LEFT JOIN votes v ON p.id = v.post_id
	    #where#
	    GROUP BY p.id
	    #orderby#
	    ) AS pq
	LEFT JOIN  users u ON u.id = uid
	#limit#
	`
)

type Post struct {
	ID           int       `json:"id,omitempty"`
	Title        string    `json:"title"`
	Url          string    `json:"url"`
	CreateAt     time.Time `json:"create_at"`
	UserID       int       `json:"user_id"`
	Votes        int       `json:"votes,omitempty"`
	UserName     string    `json:"user_name,omitempty"`
	CommentCount int       `json:"comment_count,omitempty"`
	TotalRecords int       `json:"total_records,omitempty"`
}

type PostsModel struct {
	db db.Session
}

func (m PostsModel) Table() string {
	return "posts"
}

func (m PostsModel) Get(id int) (*Post, error) {
	var post Post

	q := strings.Replace(queryTemplate, "#where#", "WHERE p.id = $1", 1)
	q = strings.Replace(q, "#orderby#", "", 1)
	q = strings.Replace(q, "#limit#", "", 1)

	row, err := m.db.SQL().Query(q, id)
	if err != nil {
		return nil, err
	}

	iter := m.db.SQL().NewIterator(row)
	err = iter.One(&post)
	if err != nil {
		return nil, err
	}

	return &post, nil
}
