// Пакет storage
// Реализует обработку RSS-потоков
// (в основе код из примера Дмитрия Титова)
//
// Проект NewsFeed
// Автор: Егор Логинов (GO-11) по заданию SkillFactory в модуле 36 (Новостной агрегатор)

package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"time"

	"newsfeed/pkg/storage"

	strip "github.com/grokify/html-strip-tags-go"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Chanel  Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

// Parse читает rss-поток из источника url и возвращет слайс раскодированных новостей.
func Parse(url string) ([]storage.Post, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var f Feed
	err = xml.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}

	var data []storage.Post
	for _, item := range f.Chanel.Items {
		var p storage.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = strip.StripTags(p.Content)
		p.URL = item.Link
		// Парсим время поста и сохраняем в int64
		t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", item.PubDate)
		}
		if err == nil {
			p.PostedAt = t.Unix()
		}
		data = append(data, p)
	}

	return data, nil
}
