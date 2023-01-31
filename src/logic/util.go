package logic

import "path/filepath"

type FileInfo struct {
	Name     string
	Path     string
	FullPath string
}

// ExtractDirPathAndName Pathからフォルダパスとファイル名を抽出する
func ExtractDirPathAndName(path string) FileInfo {
	dir, filename := filepath.Split(path)

	return FileInfo{Name: getFileNameWithoutExt(filename), Path: dir, FullPath: path}
}

// 拡張子を取り除いたファイル名を返却する
func getFileNameWithoutExt(fileName string) string {
	return filepath.Base(fileName[:len(fileName)-len(filepath.Ext(fileName))])
}
