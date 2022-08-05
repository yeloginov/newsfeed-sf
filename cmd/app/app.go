// Пакет main
// Проект NewsFeed
// Автор: Егор Логинов (GO-11) по заданию SkillFactory в модуле 36 (Новостной агрегатор)

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"newsfeed/pkg/api"
	"newsfeed/pkg/rss"
	"newsfeed/pkg/storage"
)

// Параметры подключения к БД Postgres.
const (
	DBHost     = "89.223.121.125"
	DBPort     = "5432"
	DBName     = "newsfeed"
	DBUser     = "gn_external"
	DBPassword = "Tdf_p9EXa9n"
)

// Конфигурация приложения.
type config struct {
	Sources []string `json:"rss"`
	Period  int      `json:"period"`
}

func main() {

	//  Инициализация зависимостей приложения.
	db, err := storage.New(fmt.Sprintf("postgres://%s:%s@%s:%s/%s", DBUser, DBPassword, DBHost, DBPort, DBName))
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(db)

	// Чтение и раскодирование файла конфигурации.
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Запуск парсинга новостей в отдельном потоке для каждой ссылки.
	chPosts := make(chan []storage.Post)
	chErrs := make(chan error)
	for _, url := range config.Sources {
		go parseURL(url, db, chPosts, chErrs, config.Period)
	}

	// Запись потока новостей в БД.
	go func() {
		for posts := range chPosts {
			db.StorePosts(posts)
		}
	}()

	// Обработка потока ошибок.
	go func() {
		for err := range chErrs {
			log.Println("ошибка:", err)
		}
	}()

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// Асинхронное чтение потока RSS. Раскодированные
// новости и ошибки пишутся в каналы.
func parseURL(url string, db *storage.DB, posts chan<- []storage.Post, errs chan<- error, period int) {
	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Minute * time.Duration(period))
	}
}
