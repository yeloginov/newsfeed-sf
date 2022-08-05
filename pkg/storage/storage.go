// Пакет storage
// CRUD-действия с базой данных Postgres
//
// Проект NewsFeed
// Автор: Егор Логинов (GO-11) по заданию SkillFactory в модуле 36 (Новостной агрегатор)

package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// База данных.
type DB struct {
	pool *pgxpool.Pool
}

// Публикация, получаемая из RSS.
type Post struct {
	ID       int    // идентификатор
	Title    string // заголовок
	Content  string // содержание
	URL      string // ссылка на источник
	PostedAt int64  // время публикации
}

// New(p string) конструктор объекта БД Postgress с параметрами p.
func New(p string) (*DB, error) {

	c, err := pgxpool.Connect(context.Background(), p)
	if err != nil {
		return nil, err
	}

	s := DB{
		pool: c,
	}

	return &s, nil
}

// StoreNews сохраняет посты из слайса pp в базу данных db.
func (db *DB) StorePosts(pp []Post) error {

	for _, post := range pp {
		_, err := db.pool.Exec(context.Background(), `
		INSERT INTO news(title, content, url, posted_at)
		VALUES ($1, $2, $3, $4)`,
			post.Title, post.Content, post.URL, post.PostedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// News возвращает n последних постов из базы данных db.
func (db *DB) News(n int) ([]Post, error) {

	if n < 1 {
		n = 10
	}

	rows, err := db.pool.Query(context.Background(), `
	SELECT id, title, content, url, posted_at FROM news
	ORDER BY posted_at DESC
	LIMIT $1
	`,
		n,
	)
	if err != nil {
		return nil, err
	}

	var nn []Post
	for rows.Next() {
		var p Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.URL, &p.PostedAt)
		if err != nil {
			return nil, err
		}
		nn = append(nn, p)
	}

	return nn, rows.Err()
}
