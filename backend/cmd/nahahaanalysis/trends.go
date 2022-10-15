package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ohnishi/nahaha/backend/cmd"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	mecab "github.com/shogo82148/go-mecab"
)

const ipadic = "/usr/local/lib/mecab/dic/mecab-ipadic-neologd"

var newsArticleNames = []string{"rss.jsonl", "5ch.jsonl"}

func transformTrends(src, dest string, date time.Time) error {
	dateStr := date.Format("20060102")
	var articles []cmd.NewsArticleJSON
	for _, fileName := range newsArticleNames {
		path := filepath.Join(src, dateStr, fileName)
		a, err := readArticles(path)
		if err != nil {
			fmt.Println("failed to open JSONL file.", zap.String("path", path), zap.Error(err))
			continue
		}
		articles = append(articles, a...)
	}

	contentItems := toContents(articles)
	if len(contentItems) >= 30 {
		contentItems = contentItems[:30]
	}

	content := cmd.Content{
		FormatDate: date.Format("2006/01/02"),
		Date:       date.Format(time.RFC3339),
		Items:      contentItems,
	}

	if err := writeContent(dest, dateStr+".json", content); err != nil {
		return err
	}
	return nil
}

func writeContent(dest, fileName string, c cmd.Content) error {
	f, err := cmd.CreateOutFile(filepath.Join(dest, fileName))
	if err != nil {
		return err
	}
	defer f.Close()

	err = cmd.AppendOutFile(f, c)
	if err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}
	return nil
}

// ニュース記事情報となるJSONLファイルをreadして返す
func readArticles(path string) ([]cmd.NewsArticleJSON, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	var articles []cmd.NewsArticleJSON
	d := json.NewDecoder(f)
	for d.More() {
		var article cmd.NewsArticleJSON
		if err := d.Decode(&article); err != nil {
			return nil, errors.Wrapf(err, "could not unmarshal: %v", article)
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func toContents(articles []cmd.NewsArticleJSON) []cmd.ContentItem {
	mecab, err := mecab.New(map[string]string{"dicdir": ipadic})
	if err != nil {
		panic(err)
	}
	defer mecab.Destroy()

	m := make(map[string]cmd.ContentItem)
	for _, article := range articles {
		title := strings.TrimSpace(strings.ToLower(article.Title))
		i := strings.LastIndex(title, "(")
		if i >= 0 {
			title = title[:i]
		}
		i = strings.LastIndex(title, "（")
		if i >= 0 {
			title = title[:i]
		}
		i = strings.LastIndex(title, "[")
		if i >= 0 {
			title = title[:i]
		}
		i = strings.LastIndex(title, "〈")
		if i >= 0 {
			title = title[:i]
		}
		i = strings.LastIndex(title, "【")
		if i >= 0 {
			title = title[:i]
		}
		i = strings.Index(title, "]")
		if i >= 0 {
			title = title[i:]
		}
		title = strings.ReplaceAll(title, ":", "")
		title = strings.ReplaceAll(title, "にも", "")
		node, err := mecab.ParseToNode(title)
		if err != nil {
			panic(err)
		}

		for ; !node.IsZero(); node = node.Next() {
			features := strings.Split(node.Feature(), ",")
			if features[0] == "名詞" && features[1] == "固有名詞" && features[2] == "人名" && features[3] == "一般" {
				// fmt.Println(node.String())
				word := node.Surface()
				if _, ok := excludeWord[word]; ok {
					continue
				}
				contentItem, ok := m[word]
				if !ok {
					contentItem = cmd.ContentItem{
						Word:  word,
						Count: 0,
					}
					m[word] = contentItem
				}
				a := cmd.Article{
					Title: article.Title,
					URL:   article.URL,
				}
				contentItem.Articles = append(contentItem.Articles, a)
				contentItem.Count = len(contentItem.Articles)
				m[word] = contentItem
			}
		}
	}
	var ret []cmd.ContentItem
	for _, val := range m {
		// fmt.Println(fmt.Sprintf("\"%s\":            {},", key))
		ret = append(ret, val)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Count > ret[j].Count })
	return ret[:100]
}

var excludeWord = map[string]struct{}{
	"web":     {},
	"no":      {},
	"お姉さん":    {},
	"ニート":     {},
	"ドラ":      {},
	"新劇場版":    {},
	"at":      {},
	"風俗嬢":     {},
	"加藤純一":    {},
	"な！":      {},
	"jk":      {},
	"alt":     {},
	"life":    {},
	"rtx":     {},
	"クラスター":   {},
	"body":    {},
	"mark":    {},
	"ceo":     {},
	"king":    {},
	"id":      {},
	"ユニ":      {},
	"d2":      {},
	"shadows": {},
	"v2":      {},
	"ko":      {},
	"ai":      {},
}
