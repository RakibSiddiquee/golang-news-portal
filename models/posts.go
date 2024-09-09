package models

import (
	"database/sql"
	"errors"
	"github.com/golang-module/carbon/v2"
	"github.com/upper/db/v4"
	"net/url"
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

func (m PostsModel) GetAll(f Filter) ([]Post, Metadata, error) {
	var posts []Post
	var rows *sql.Rows
	var err error
	meta := Metadata{}

	q := f.applyTemplate(queryTemplate)

	if len(f.Query) > 0 {
		rows, err = m.db.SQL().Query(q, "%"+strings.ToLower(f.Query)+"%", f.limit(), f.offset())
	} else {
		rows, err = m.db.SQL().Query(q, f.limit(), f.offset())
	}

	if err != nil {
		return nil, meta, err
	}

	iter := m.db.SQL().NewIterator(rows)
	err = iter.All(&posts)

	if err != nil {
		return nil, meta, err
	}

	if len(posts) == 0 {
		return nil, meta, errors.New("no record found")
	}

	first := posts[0]
	return posts, calculateMetadata(first.TotalRecords, f.Page, f.PageSize), nil
}

func (m PostsModel) Vote(postId, userId int) error {
	col := m.db.Collection("votes")

	_, err := col.Insert(map[string]int{
		"post_id": postId,
		"user_id": userId,
	})

	if err != nil {
		if errHasDuplicate(err, "votes_pkey") {
			return ErrDuplicateVotes
		}
		return err
	}
	return nil
}

func (p *Post) DateHuman() string {
	return carbon.CreateFromStdTime(p.CreateAt).DiffForHumans()
}

func (p *Post) Host() string {
	ur, err := url.Parse(p.Url)
	if err != nil {
		return ""
	}
	return ur.Host
}
