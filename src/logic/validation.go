package logic

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// IsExistFile 対象ファイルが存在するかチェックします
func IsExistFile(path string) error {
	fInfo, err := os.Stat(path)
	if err != nil {
		return errors.New("指定したパスは存在しません")
	}

	if fInfo.IsDir() {
		return errors.New("指定したパスはファイルではありません")
	}
	return nil
}

// CheckExtension 指定の拡張子かどうかチェックする
func CheckExtension(path string, exts []string) error {
	target := filepath.Ext(path)
	for _, ext := range exts {
		if target == ext {
			return nil
		}
	}
	return fmt.Errorf("拡張子は%sを指定してください", exts)
}
