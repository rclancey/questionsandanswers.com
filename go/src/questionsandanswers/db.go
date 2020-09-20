package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/uuid"
)

const (
	SortAskDate = "questionDate"
	SortAskDateDesc = "-questionDate"
	SortAnswerDate = "answerDate"
	SortAnswerDateDesc = "-answerDate"
)

type DB struct {
	db *sqlx.DB
}

type Question struct {
	ID *string `json:"id"`
	Question *string `json:"question"`
	QuestionTimeMS *int64 `json:"questionDate" db:"question_date"`
	Answer *string `json:"answer,omitempty"`
	AnswerTimeMS *int64 `json:"answerDate,omitempty" db:"answer_date"`
}

func NewQuestion(question string) (*Question, error) {
	q := question
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	idStr := id.String()
	t := time.Now()
	tMs := t.Unix() * 1000 + int64(t.Nanosecond() / 1e6)
	return &Question{
		ID: &idStr,
		Question: &q,
		QuestionTimeMS: &tMs,
		Answer: nil,
		AnswerTimeMS: nil,
	}, nil
}

func (q *Question) SetAnswer(answer string) error {
	a := answer
	t := time.Now()
	tMs := t.Unix() * 1000 + int64(t.Nanosecond() / 1e6)
	q.Answer = &a
	q.AnswerTimeMS = &tMs
	return nil
}

func NewDB(filename string) (*DB, error) {
	needsInit := false
	st, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			needsInit = true
		} else {
			return nil, err
		}
	} else if st.IsDir() {
		return nil, errors.New("database file is a directory")
	}
	conn, err := sqlx.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	db := &DB{conn}
	if needsInit {
		err = db.initDB()
		if err != nil {
			conn.Close()
			return nil, err
		}
	}
	return db, nil
}

func (db *DB) initDB() error {
	qs := `CREATE TABLE question (id VARCHAR(36) NOT NULL PRIMARY KEY, question TEXT NOT NULL, question_date BIGINT NOT NULL, answer TEXT, answer_date BIGINT)`
	_, err := db.db.Exec(qs)
	return err
}

func (db *DB) Close() error {
	if db.db == nil {
		return nil
	}
	err := db.db.Close()
	if err != nil {
		return err
	}
	db.db = nil
	return nil
}

func (db *DB) CreateQuestion(question string) (*Question, error) {
	q, err := NewQuestion(question)
	if err != nil {
		return nil, err
	}
	qs := `INSERT INTO question (id, question, question_date, answer, answer_date) VALUES(?, ?, ?, ?, ?)`
	_, err = db.db.Exec(qs, q.ID, q.Question, q.QuestionTimeMS, q.Answer, q.AnswerTimeMS)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (db *DB) ReadQuestion(questionId string) (*Question, error) {
	qs := `SELECT * FROM question WHERE id = ?`
	row := db.db.QueryRowx(qs, questionId)
	q := &Question{}
	err := row.StructScan(q)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (db *DB) UpdateQuestion(q *Question) error {
	qs := `UPDATE question SET question = ?, question_date = ?, answer = ?, answer_date = ? WHERE id = ?`
	res, err := db.db.Exec(qs, q.Question, q.QuestionTimeMS, q.Answer, q.AnswerTimeMS, q.ID)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return ErrNoSuchQuestion
	}
	return nil
}

func (db *DB) DeleteQuestion(questionId string) error {
	qs := `DELETE FROM question WHERE id = ?`
	res, err := db.db.Exec(qs, questionId)
	if err != nil {
		return err
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return ErrNoSuchQuestion
	}
	return nil
}

func (db *DB) CountQuestions() (int, error) {
	qs := `SELECT COUNT(*) FROM question`
	row := db.db.QueryRow(qs)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (db *DB) ListQuestions(sortOrder string, pageSize, pageNum int) ([]*Question, error) {
	qs := `SELECT * FROM question`
	switch sortOrder {
	case SortAskDate:
		qs += " ORDER BY question_date"
	case SortAskDateDesc:
		qs += " ORDER BY question_date DESC"
	case SortAnswerDate:
		qs += " ORDER BY answer_date"
	case SortAnswerDateDesc:
		qs += " ORDER BY answer_date DESC"
	default:
		qs += " ORDER BY id"
	}
	qs += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, pageSize * pageNum)
	rows, err := db.db.Queryx(qs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	questions := make([]*Question, pageSize)
	count := 0
	for rows.Next() {
		q := &Question{}
		err := rows.StructScan(q)
		if err != nil {
			return questions[:count], err
		}
		questions[count] = q
		count += 1
	}
	return questions[:count], nil
}
