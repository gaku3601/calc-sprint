package logic

import (
	"os"
)

type OperationFile struct {
	fp *os.File
}

func NewFile(outputPath string) (*OperationFile, error) {
	os.Remove(outputPath)

	f := new(OperationFile)
	fp, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	f.fp = fp
	return f, nil
}

func (f OperationFile) Close() error {
	return f.fp.Close()
}
