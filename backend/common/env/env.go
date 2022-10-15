package env

import (
	"go/build"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// AppEnvType はアプリケーションの環境を示す文字列。
type AppEnvType int

const (
	// Development はdevelopmentのアプリケーションの環境を示す
	Development AppEnvType = iota
	// Production はproductionのアプリケーションの環境を示す
	Production
	// Test はgo test時に使用されるのアプリケーションの環境を示す
	Test
	// Unknown はparseできなかった場合などのerror値を示す
	Unknown
)

func (e AppEnvType) String() string {
	switch e {
	case Development:
		return "development"
	case Production:
		return "production"
	case Test:
		return "test"
	default:
		return "unknown"
	}
}

// ParseAppEnv はstringをparseしAppEnvTypeにconvertする
func ParseAppEnv(s string) (AppEnvType, error) {
	switch s {
	case Development.String():
		return Development, nil
	case Production.String():
		return Production, nil
	case Test.String():
		return Test, nil
	default:
		return Unknown, errors.Errorf("app env %s is not supported", s)
	}
}

// DataDir はデータディレクトリのパスを返す。
func DataDir(subpaths ...string) string {
	return filepath.Join(AppRootDir(), "data", filepath.Join(subpaths...))
}

// ConfigDir は設定ディレクトリのパスを返す。
func ConfigDir(subpaths ...string) string {
	return filepath.Join(AppRootDir(), "config", filepath.Join(subpaths...))
}

// LogDir はログディレクトリのパスを返す。
func LogDir(subpaths ...string) string {
	return DataDir("log", filepath.Join(subpaths...))
}

// LogLevel はロガーのログレベルを返します。
// 指定可能なlog levelはzapcore.Levelに準拠します。
// 明示的に指定されていないか不正な場合はfalseが返ります
func LogLevel() (*zapcore.Level, bool) {
	logLevelText := getEnv("LOG_LEVEL", "")
	return logLevelFromText(logLevelText)
}

func logLevelFromText(logLevelText string) (*zapcore.Level, bool) {
	if logLevelText == "" {
		return nil, false
	}

	logLevel := new(zapcore.Level)
	err := logLevel.UnmarshalText([]byte(logLevelText))
	if err != nil {
		return nil, false
	}

	return logLevel, true
}

// AppEnv はアプリケーションの環境(development|production|test)を返す。
// APP_ENVがセットされていないか不正な値であった場合はdevelopmentを返す。
func AppEnv() AppEnvType {
	env := getEnv("APP_ENV", Development.String())
	appEnvType, err := ParseAppEnv(env)
	if err != nil {
		return Development
	}
	return appEnvType
}

// AppRootDir はconfig, dataディレクトリを置くためのルートディレクトリ。
// デフォルトは$GOPATH。
// quaと同じGOPATH下でconfig等を切り替えたい場合に指定すること。
func AppRootDir() string {
	return getEnv("APP_ROOT_DIR", build.Default.GOPATH)
}

func getEnv(name string, alt string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return alt
}
