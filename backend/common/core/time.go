package core

import (
	"time"

	"github.com/pkg/errors"
)

// NewLocalDate は time.Date(year, month, date, 0, 0, 0, 0, time.Local) のショートハンド
func NewLocalDate(year int, month time.Month, date int) time.Time {
	return time.Date(year, month, date, 0, 0, 0, 0, time.Local)
}

// ParseLocal はvalueをローカルの時間とみなしてパースした結果を返す。
func ParseLocal(layout string, value string) (time.Time, error) {
	t, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "cannot parse as %q", time.Local)
	}
	return t, nil
}

// TruncateDayLocal はローカルの時間のもとで日より小さい時間をTruncateした結果を返す。
func TruncateDayLocal(t time.Time) time.Time {
	return t.Truncate(time.Hour).Add(-time.Duration(t.Hour()) * time.Hour)
}
