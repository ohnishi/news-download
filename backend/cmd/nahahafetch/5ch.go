package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ohnishi/nahaha/backend/cmd"
	"github.com/ohnishi/nahaha/backend/common/core"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// 5ch板一覧URL
const boardListURL = "https://menu.5ch.net/bbstable.html"

// 板以外のリンクURL
var excludeLink = map[string]struct{}{
	"https://www.5ch.net/":             {}, // 5chの入り口
	"https://www.5ch.net/kakolog.html": {}, // 過去ログ倉庫
	"https://newsnavi.5ch.net/":        {}, // 2NN+
	"https://info.5ch.net/":            {}, // 5ch総合案内
	"https://search.5ch.net/":          {}, // 検索[ベータ版]
	"https://dig.5ch.net/":             {}, // 超スレタイ検索
	"https://stat.5ch.net/":            {}, // 5ch投稿数
	"https://o.5ch.net/":               {}, // お絵描き観測所
	"https://i.5ch.net/":               {}, // スマホメニュー
	"https://be.5ch.net/":              {}, // be.5ch.net
	"https://premium.5ch.net/":         {}, // 5chプレミアム浪人
	"https://info.5ch.net/wiki/":       {}, // 5chプロジェクト
	"https://matsuri.5ch.net/maru/":    {}, // ●
	"https://info.5ch.net/?curid=2078": {}, // 書き込む前に
	"mailto:admin@5ch.net":             {}, // メール
	"https://www.bbspink.com/":         {}, // Pinkちゃんねる
	"https://ronin.bbspink.com/":       {}, // 浪人
	"https://info.5ch.net/rank/":       {}, // いろいろランク
}

type link struct {
	text string
	href string
}

// fetchNews5ch  ニュースソースとなる5chの subject.txt を保存する
func fetchNews5ch(dest string, maxRetry uint) (err error) {
	var links []link
	retry := uint(0)
	for {
		links, err = func() ([]link, error) {
			res, err := http.Get(boardListURL)
			if err != nil {
				return nil, errors.Wrapf(err, "failed request url : %s", boardListURL)
			}
			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return nil, errors.Errorf("status code expected 200 but was %d : url=%s", res.StatusCode, boardListURL)
			}

			links, err := getLinks(res.Body)
			if err != nil {
				return nil, err
			}
			return links, nil
		}()
		retry++
		if err == nil || retry > maxRetry {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return err
	}

	fetchDir := filepath.Join(dest, time.Now().Format("20060102"))
	boards := toBoards(links, fetchDir, maxRetry)

	return saveBoards(fetchDir, boards)
}

// スクレイピングで取得した板URLの一覧から有効な板の情報を取得して返す
func toBoards(links []link, out string, maxRetry uint) []cmd.Board {
	linkSet := core.StringSet{}
	var boards []cmd.Board
	for _, l := range links {
		if _, ok := excludeLink[l.href]; ok {
			// 板以外へのリンクURLならスキップする
			continue
		}
		if linkSet.Include(l.href) {
			//たまに重複したリンクURLがあるので重複チェックする
			continue
		}
		linkSet.Add(l.href)
		boardID, err := getSubject(l.href, out, maxRetry)
		if err != nil {
			// subject.txtの取得に失敗しても処理は止めずに warnnig log を出力する
			fmt.Println("failed to fetch subject.txt.", zap.String("url", l.href), zap.Error(err))
			continue
		}

		name, _, err := transform.String(japanese.ShiftJIS.NewDecoder(), l.text)
		if err != nil {
			fmt.Println("failed to decode board name", zap.String("text", l.text), zap.Error(err))
			continue
		}

		boards = append(boards, cmd.Board{ID: boardID, Name: name, URL: l.href})
	}
	return boards
}

// 板一覧HTMLをパースして板URLと板名を取得する
func getLinks(r io.Reader) ([]link, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read response body")
		}
		return nil, errors.Wrapf(err, "failed parse response body : %s", string(b))
	}

	var links []link
	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			l := link{
				text: s.Text(),
				href: href,
			}
			links = append(links, l)
		}
	})
	return links, nil
}

// スクレイピングで取得した板URLからsubject.txtを取得して保存する
func getSubject(href, out string, maxRetry uint) (string, error) {
	boardURL, err := url.Parse(href)
	if err != nil {
		return "", errors.Wrapf(err, "failed parse url : %s", href)
	}

	// 先頭、末尾のスラッシュを除去
	boardID := strings.Replace(boardURL.Path, "/", "", 2)
	if boardID == "" {
		return boardID, errors.Errorf("illegal 5ch board URL : %s", boardURL.String())
	}
	boardURL.Path = path.Join(boardURL.Path, "subject.txt")

	var res *http.Response
	retry := uint(0)
	for {
		res, err = http.Get(boardURL.String())
		retry++
		if err == nil || retry > maxRetry {
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return boardID, errors.Wrapf(err, "failed request url : %s", boardURL.String())
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return boardID, errors.Errorf("status code expected 200 but was %d : url=%s", res.StatusCode, boardURL.String())
	}

	return boardID, save(res, filepath.Join(out, boardID))
}

// 板URLと板名を保存する
func saveBoards(out string, boards []cmd.Board) error {
	if len(boards) == 0 {
		return nil
	}

	fi := cmd.FetchInfo{
		Date:   time.Now().Format(time.RFC3339),
		Boards: boards,
	}

	path := filepath.Join(out, "fetch_info.json")
	f, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create file: %s", out)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(fi)
	if err != nil {
		return errors.Wrapf(err, "failed to write json : path=%s, value=%v", path, fi)
	}
	return nil
}
