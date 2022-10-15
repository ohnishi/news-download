package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ohnishi/nahaha/backend/cmd"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// 5ch スレッドタイトルから除外するワード
var replaceThreadTitleWords = []string{
	"&#169;2ch.net",
	"&copy;2ch.net",
	"&#169;bbspink.com",
	"&copy;bbspink.com",
	"[無断転載禁止]",
	"[転載禁止]",
}

// transform5ch fetchしたsubject.txtからターゲット日に更新された5chスレッドを抽出する
func transform5ch(src, dest string, date time.Time) error {
	dateStr := date.Format("20060102")
	threadMap, err := toThreadMap(src, dateStr)
	if err != nil {
		return err
	}

	return writeArticleJSOL(dest, dateStr, "5ch.jsonl", threadMap)
}

// ターゲット日に作成されたスレッド情報を返す
func toThreadMap(src string, dateStr string) (map[string]newsArticleJSON, error) {
	m := make(map[string]newsArticleJSON)
	fileDir := filepath.Join(src, dateStr)

	fetchInfo, err := readFetchInfo(fileDir)
	if err != nil {
		return nil, err
	}

	for _, b := range fetchInfo.Boards {
		subjectTextPath := filepath.Join(fileDir, b.ID)

		file, err := os.Open(subjectTextPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read file: %s", subjectTextPath)
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			s := scanner.Text()
			records := strings.Split(s, "<>")
			if len(records) != 2 {
				fmt.Println("unexpected string", zap.String("string", s))
				continue
			}
			threadKey := toThreadKey(records[0])

			threadSec, err := strconv.Atoi(threadKey)
			if err != nil {
				fmt.Println("unexpected thread key", zap.String("threadKey", threadKey), zap.Error(err))
				continue
			}

			threadDate := time.Unix(int64(threadSec), 0)
			if dateStr != threadDate.Format("20060102") {
				continue
			}

			url, err := toThreadURL(b.URL, b.ID, threadKey)
			if err != nil {
				fmt.Println("failed to generate thread URL", zap.String("url", b.URL), zap.String("id", b.ID), zap.String("threadKey", threadKey), zap.Error(err))
				continue
			}

			if _, ok := m[url]; ok {
				continue
			}

			threadTitle, err := toThreadTitle(records[1])
			if err != nil {
				fmt.Println("failed to generate thread title", zap.String("string", records[1]), zap.Error(err))
				continue
			}
			if len(threadTitle) > 512 {
				//512文字以上のタイトルならDBに挿入不可能かつ、画面表示も難しいためスキップ
				continue
			}

			json := newsArticleJSON{
				Date:  threadDate.Format(time.RFC3339),
				URL:   url,
				Name:  b.Name,
				Title: threadTitle,
			}
			m[url] = json
		}
		if err := scanner.Err(); err != nil {
			return nil, errors.Wrapf(err, "failed to read file: %s", subjectTextPath)
		}

		err = file.Close()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to close file: %s", subjectTextPath)
		}
	}
	return m, nil
}

// fetchした板の名前一覧情報を返す
func readFetchInfo(dir string) (cmd.FetchInfo, error) {
	var fetchInfo cmd.FetchInfo

	jsonPath := filepath.Join(dir, "fetch_info.json")
	content, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return fetchInfo, errors.Wrapf(err, "failed to read file: %s", jsonPath)
	}

	err = json.Unmarshal(content, &fetchInfo)
	if err != nil {
		return fetchInfo, errors.Wrapf(err, "failed to unmarshal json: %s", jsonPath)
	}
	return fetchInfo, nil
}

// subject.txtの行文字列からスレッドキーを抽出して返す
func toThreadKey(s string) string {
	lastIndex := strings.LastIndex(s, ".dat")
	if lastIndex >= 0 {
		return s[:lastIndex]
	}
	return s
}

// subject.txtの行文字列からスレッドタイトルを抽出して返す
func toThreadTitle(s string) (string, error) {
	threadTitle, _, err := transform.String(japanese.ShiftJIS.NewDecoder(), s)
	if err != nil {
		return "", errors.Wrapf(err, "failed to encode thread title : %v", threadTitle)
	}
	li := strings.LastIndex(threadTitle, "(")
	if li >= 0 {
		threadTitle = threadTitle[:li]
	}
	for _, w := range replaceThreadTitleWords {
		threadTitle = strings.Replace(threadTitle, w, "", 1)
	}
	threadTitle = strings.TrimSpace(threadTitle)

	return threadTitle, nil
}

// 5ch スレッドURLを生成して返す
func toThreadURL(threadURL, threadID, threadKey string) (string, error) {
	lastIndex := strings.LastIndex(threadURL, threadID)
	if lastIndex < 0 {
		return "", errors.Errorf("unexpected thread URL : %s", threadURL)
	}
	return fmt.Sprintf("%stest/read.cgi/%s/%s/", threadURL[:lastIndex], threadID, threadKey), nil
}
