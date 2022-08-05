package storage

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// Параметры подключения к БД Postgres.
const (
	DBHost     = "89.223.121.125"
	DBPort     = "5432"
	DBName     = "newsfeed"
	DBUser     = "gn_external"
	DBPassword = "Tdf_p9EXa9n"
)

func TestNew(t *testing.T) {
	_, err := New(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDB_News(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	posts := []Post{
		{
			Title: "Test Post",
			URL:   strconv.Itoa(rand.Intn(1_000_000_000)),
		},
	}
	db, err := New(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		t.Fatal(err)
	}
	err = db.StorePosts(posts)
	if err != nil {
		t.Fatal(err)
	}
	news, err := db.News(2)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", news)
}
