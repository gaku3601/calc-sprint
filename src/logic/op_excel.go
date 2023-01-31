package logic

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type OperationExcel struct {
	file   *excelize.File
	sheets []string
}

// NewOperationExcel constructor
func NewOperationExcel(path string) (*OperationExcel, error) {
	o := new(OperationExcel)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	o.file = f
	// Your Jira Issuesのシートのみ抽出
	for _, sheet := range o.file.GetSheetMap() {
		if sheet == "Your Jira Issues" {
			o.sheets = append(o.sheets, sheet)
		}
	}
	return o, nil
}

func (o OperationExcel) Execute() error {
	// シート毎に処理する
	for _, sheet := range o.sheets {
		rows := o.file.GetRows(sheet)
		var headers []string
		var values [][]string
		for i, row := range rows {
			if i == 0 {
				headers = analyzeHeaders(row)
			} else {
				values = append(values, analyzeValues(row, len(headers)))
			}
		}
		if err := calcTimeToComplete(headers, values); err != nil {
			return err
		}
	}
	return nil
}

// 消費時間の位置を取得する
func getConsumptionTimePosition(headers []string) (int, error) {
	for i, col := range headers {
		if col == "Σ 消費時間" {
			return i, nil
		}
	}
	return 0, errors.New("Σ 消費時間がヘッダーに含まれていません")
}

// 解決状況の位置を取得する
func getResolutionPosition(headers []string) (int, error) {
	for i, col := range headers {
		if col == "解決状況" {
			return i, nil
		}
	}
	return 0, errors.New("解決状況がヘッダーに含まれていません")
}

// Story point estimateの位置を取得する
func getStoryPointEstimatePosition(headers []string) (int, error) {
	for i, col := range headers {
		if col == "Story point estimate" {
			return i, nil
		}
	}
	return 0, errors.New("Story point estimateがヘッダーに含まれていません")
}

// 要約の位置を取得する
func getSummaryPosition(headers []string) (int, error) {
	for i, col := range headers {
		if col == "要約" {
			return i, nil
		}
	}
	return 0, errors.New("要約がヘッダーに含まれていません")
}

/**
完了タスクのすべての時間を合算した値を計算し算出する
*/
func calcTimeToComplete(headers []string, values [][]string) error {
	ctp, err := getConsumptionTimePosition(headers)
	if err != nil {
		return err
	}
	rp, err := getResolutionPosition(headers)
	if err != nil {
		return err
	}
	spe, err := getStoryPointEstimatePosition(headers)
	if err != nil {
		return err
	}
	sp, err := getSummaryPosition(headers)
	if err != nil {
		return err
	}

	var consumptionTimes int
	var storyPoints int
	var bufferTime int
	for _, row := range values {
		// バッファ以外のタスクでステータスが完了で時間が格納されているものだけ抽出して計算に必要な値を算出
		if row[rp] == "完了" && row[ctp] != "" && row[sp] != "バッファ" {
			consumptionTime, err := strconv.Atoi(row[ctp])
			if err != nil {
				return err
			}
			consumptionTimes += consumptionTime

			if row[spe] == "" {
				return fmt.Errorf("ストーリーポイントが入っていないストーリーが存在します: %s", row[sp])
			}
			storyPoint, err := strconv.Atoi(row[spe])
			if err != nil {
				return err
			}
			storyPoints += storyPoint
		}
		// バッファタスクの消費時間を算出
		if row[sp] == "バッファ" {
			bufferTime, err = strconv.Atoi(row[ctp])
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf("バッファ以外の完了タスクにおけるトータル消費時間: %.1fh\n", float64(consumptionTimes)/60/60)
	fmt.Printf("完了したストーリーポイント: %d\n", storyPoints)
	fmt.Printf("バッファで消費した時間: %.1fh\n", float64(bufferTime)/60/60)
	return nil
}

// 値が入っているもののみheaderとして認識して返却する
func analyzeHeaders(row []string) []string {
	var cols []string
	for _, col := range row {
		if len(col) > 0 {
			cols = append(cols, col)
		}
	}
	return cols
}

// headerの列の数を実際の値であると認識して返却する
func analyzeValues(row []string, headerCount int) []string {
	var cols []string
	for i, col := range row {
		if i < headerCount {
			cols = append(cols, deleteNewLineExcelCode(col))
		}
	}
	return cols
}

// 不要な改行コードを削除する
func deleteNewLineExcelCode(col string) string {
	return strings.Replace(col, "_x000D_", "", -1)
}

// シートを並べ替える
func (o OperationExcel) organizeSheets() {
	sort.SliceStable(o.sheets, func(i, j int) bool { return o.sheets[i] < o.sheets[j] })
}

// シートからテーブル名を抽出する
func extractTableName(sheet string) string {
	rep := regexp.MustCompile(`\d*\.`)
	return rep.ReplaceAllString(sheet, "")
}
