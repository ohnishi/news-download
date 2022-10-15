package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Content struct {
	FormatDate string        `json:"format_date"`
	Date       string        `json:"date"`
	Items      []ContentItem `json:"items"`
}

type ContentItem struct {
	Word     string    `json:"word"`
	Count    int       `json:"count"`
	Articles []Article `json:"articles"`
}

type Article struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type YahooRSSFeed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

// FetchInfo fetchした板の名前一覧情報
type FetchInfo struct {
	Date   string  `json:"date"`
	Boards []Board `json:"boards"`
}

type Board struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type NewsArticleJSON struct {
	Date     string `json:"date"`
	URL      string `json:"url"`
	Name     string `json:"name"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

func ReadYahooRSSFeed(path string) ([]YahooRSSFeed, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	var articles []YahooRSSFeed
	d := json.NewDecoder(f)
	for d.More() {
		var article YahooRSSFeed
		if err := d.Decode(&article); err != nil {
			return nil, errors.Wrapf(err, "could not unmarshal: %v", article)
		}
		articles = append(articles, article)
	}
	return articles, nil
}

// CreateOutFile データ書き込み用のファイルを生成する
func CreateOutFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create output directory: %s", dir)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	return f, nil
}

// AppendOutFile はファイルにJSONを追記する
func AppendOutFile(f *os.File, v interface{}) error {
	jsonl, err := toJSON(v)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(jsonl)); err != nil {
		return errors.Wrapf(err, "failed to write line: %s", jsonl)
	}
	return nil
}

//JSON文字列に変換する
func toJSON(r interface{}) (string, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return "", errors.Wrapf(err, "could not marshal: %v", r)
	}
	return fmt.Sprintf("%s\n", jsonStr), nil
}

// ReadFileJSON はJSONを読み込んでoutputに入れる
func ReadFileJSON(file string, output interface{}) error {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, output)
}
