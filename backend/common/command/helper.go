package command

import (
	"fmt"
	"os"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ohnishi/nahaha/backend/common/core"
)

// DatesFlagFormat は`--date`フラグで用いる日付のフォーマットを表す。
const DatesFlagFormat = "20060102"

// StringSliceVarSetter はStringSliceフラグをセットするインタフェースを表す
type StringSliceVarSetter interface {
	StringSliceVar(p *[]string, name string, value []string, usage string)
}

// SetRangeFlag は`--<引数で渡された名前>`フラグをセットアップする。
func SetRangeFlag(f StringSliceVarSetter, p *[]string, name string, purpose string) {
	const (
		format = "%s in 'YYYYmmdd' or period in 'YYYYmmdd,YYYYmmdd' " +
			"(e.g. date: '20180101', period: '20180101,20180131')"
	)
	f.StringSliceVar(p, name, []string{}, fmt.Sprintf(format, purpose))
}

// SetDatesFlag は`--date`フラグをセットアップする。
func SetDatesFlag(f StringSliceVarSetter, p *[]string, purpose string) {
	SetRangeFlag(f, p, "date", purpose)
}

// SetMonthsFlag は`--month`フラグをセットアップする。
func SetMonthsFlag(f StringSliceVarSetter, p *[]string, purpose string) {
	SetRangeFlag(f, p, "month", purpose)
}

// StringVarSetter はStringフラグをセットするインタフェースを表す
type StringVarSetter interface {
	StringVar(p *string, name string, value string, usage string)
}

// SetDataDirFlag は`--data-dir`フラグをセットアップする。
func SetDataDirFlag(f StringVarSetter, p *string, v string, purpose string) {
	const (
		name   = "data-dir"
		format = "data directory path %s"
	)
	f.StringVar(p, name, v, fmt.Sprintf(format, purpose))
}

// PrintErrorAndExit はエラーを詳細に出力してステータスコード1で終了する。
func PrintErrorAndExit(err error) {
	fmt.Printf("ERROR: %+v\n", err)
	os.Exit(1)
}

// ParseDayRange は日付範囲を含む文字列のスライスからtime.Timeで開始日と終了日を返す。
func ParseDayRange(date []string) (time.Time, time.Time, error) {
	var s, e string
	switch len(date) {
	case 1:
		s = date[0]
		e = date[0]
	case 2:
		s = date[0]
		e = date[1]
	default:
		return time.Time{}, time.Time{}, errors.Errorf("invalid date: %v", date)
	}

	start, err := core.ParseLocal("20060102", s)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Errorf("failed to parse start date: %s", s)
	}
	end, err := core.ParseLocal("20060102", e)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Errorf("failed to parse end date: %s", e)
	}
	return start, end, nil
}

// EachDate はdateで指定された範囲の期間における日付ごとに、引数fnを実行する。
func EachDate(date []string, fn func(time.Time) error) error {
	step, err := getStepFunc("daily")
	if err != nil {
		return err
	}
	return EachByStep(date, step, fn)
}

// EachWeek はweekで指定された範囲の期間における週ごとに引数fnを実行する
func EachWeek(week []string, fn func(time.Time) error) error {
	step, err := getStepFunc("weekly")
	if err != nil {
		return err
	}
	return EachByStep(week, step, fn)
}

// EachMonth はmonthで指定された範囲の期間における月ごとに、引数fnを実行する。
func EachMonth(month []string, fn func(time.Time) error) error {
	for i := 0; i < len(month); i++ {
		var parsed time.Time
		parsed, err := parseMonthLocal(month[i])
		if err == nil {
			month[i] = parsed.Format(DatesFlagFormat)
		}
	}
	step, err := getStepFunc("monthly")
	if err != nil {
		return err
	}
	return EachByStep(month, step, fn)
}

// EachQuarter はweekで指定された範囲の期間における四半期ごとに引数fnを実行する
func EachQuarter(quarter []string, fn func(time.Time) error) error {
	for i := 0; i < len(quarter); i++ {
		var parsed time.Time
		parsed, err := parseMonthLocal(quarter[i])
		if err == nil {
			quarter[i] = parsed.Format(DatesFlagFormat)
		}
	}
	step, err := getStepFunc("quarterly")
	if err != nil {
		return err
	}
	return EachByStep(quarter, step, fn)
}

// EachYear はyearで指定された範囲の期間における年ごとに引数fnを実行する
func EachYear(year []string, fn func(time.Time) error) error {
	for i := 0; i < len(year); i++ {
		var parsed time.Time
		parsed, err := parseYearLocal(year[i])
		if err == nil {
			year[i] = parsed.Format(DatesFlagFormat)
		}
	}
	step, err := getStepFunc("yearly")
	if err != nil {
		return err
	}
	return EachByStep(year, step, fn)
}

// EachByStep はdateで指定された範囲の期間におけるstepごとに、引数fnを実行する。
func EachByStep(date []string, step func(time.Time) time.Time, fn func(time.Time) error) error {
	switch len(date) {
	case 0:
		return errors.New("one or two date values must be specified")
	case 1:
		d, err := core.ParseLocal(DatesFlagFormat, date[0])
		if err != nil {
			return err
		}
		return fn(d)
	case 2:
		since, err := core.ParseLocal(DatesFlagFormat, date[0])
		if err != nil {
			return err
		}
		until, err := core.ParseLocal(DatesFlagFormat, date[1])
		if err != nil {
			return err
		}
		if since.After(until) {
			since, until = until, since
		}
		var errs error
		for d := since; !d.After(until); d = step(d) {
			err = fn(d)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		}
		return errs
	default:
		return errors.New("more than 2 values cannot be specified for date")
	}
}

// getStepFunc は period に応じた対象日付の次の期間の始めの日付を取得する関数を取得する
func getStepFunc(period string) (func(time.Time) time.Time, error) {
	switch period {
	case "daily":
		return func(d time.Time) time.Time {
			return d.AddDate(0, 0, 1)
		}, nil
	case "weekly":
		return func(d time.Time) time.Time {
			return d.AddDate(0, 0, 7)
		}, nil
	case "monthly":
		return func(d time.Time) time.Time {
			return d.AddDate(0, 1, 0)
		}, nil
	case "quarterly":
		return func(d time.Time) time.Time {
			return d.AddDate(0, 3, 0)
		}, nil
	case "yearly":
		return func(d time.Time) time.Time {
			return d.AddDate(1, 0, 0)
		}, nil
	default:
		return nil, errors.Errorf("invalid period: %s", period)
	}
}

const monthFormat = "200601"

func parseMonthLocal(month string) (time.Time, error) {
	t, err := time.ParseInLocation(monthFormat, month, time.Local)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "timeParseInLocation failed: month=%s", month)
	}
	return t, nil
}

const yearFormat = "2006"

func parseYearLocal(year string) (time.Time, error) {
	t, err := time.ParseInLocation(yearFormat, year, time.Local)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "timeParseInLocation failed: year=%s", year)
	}
	return t, nil
}
