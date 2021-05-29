package amz

import (
	"errors"
	"github.com/kerbalwzy/SinRiuTool/goApp/utils"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var ReturnDataMatchWord = map[string]bool{
	"Ajuste de tarifa":                true,
	"Ajustement des frais":            true,
	"Återbetalning":                   true,
	"Chargeback Refund":               true,
	"Erstattung":                      true,
	"Erstattung durch Rückbuchung":    true,
	"Fee Adjustment":                  true,
	"Gebührenerstattung":              true,
	"Reembolso":                       true,
	"Reembolso de reversión de cargo": true,
	"Refund":                          true,
	"Refund_Retrocharge":              true,
	"Regolazione dei costi":           true,
	"Remboursement":                   true,
	"Remboursement de contestation de prélèvement": true,
	"Rimborso":      true,
	"Terugbetaling": true,
	"返金":            true,
}

func IsReturn(text string) bool {
	_, ok := ReturnDataMatchWord[text]
	return ok
}

func readFileData(dirPath, filePath string) (string, [][]string, error) {
	fileTip := strings.ReplaceAll(filePath, dirPath, "")
	res := make([][]string, 0)
	titleAdded := false
	csvData, err := utils.ReadCSV(filePath)
	if nil != err {
		return "", nil, err
	}
	for _, row := range csvData {
		if len(row) < 5 || row[6] == "" {
			continue
		}
		if !titleAdded {
			titleAdded = true
			res = append(res, row)
		}
		if IsReturn(row[2]) {
			res = append(res, row)
		}
	}
	return fileTip, res, nil
}

type OrderReturn struct {
	DirPath       string
	filePathS     []string
	Errors        []error
	HandledFiles  chan string
	ExcelSheetS   []utils.ExcelSheet
	ExcelSavePath string
}

func (obj *OrderReturn) Init() {
	filePathS, err := utils.ListDirFiles(obj.DirPath, ".csv")
	if nil != err {
		obj.Errors = append(obj.Errors, err)
		return
	}
	obj.filePathS = filePathS
	obj.HandledFiles = make(chan string, len(obj.filePathS))
	wait := new(sync.WaitGroup)
	wait.Add(len(obj.filePathS))
	for _, item := range obj.filePathS {
		go func(filepath string) {
			isUTF8, err := utils.ValidFileUTF8(filepath, 10)
			if nil != err {
				obj.Errors = append(obj.Errors, err)
			} else if !isUTF8 {
				obj.Errors = append(obj.Errors, errors.New(filepath+"\t不能被UTF8解码!"))
			}
			wait.Done()
		}(item)
	}
	wait.Wait()

}

func (obj *OrderReturn) ReadFilesData() {
	if len(obj.filePathS) == 0 {
		panic(errors.New("请先调用Init方法"))
	}
	if len(obj.Errors) != 0 {
		return
	}
	wait := sync.WaitGroup{}
	lock := sync.Mutex{}
	wait.Add(len(obj.filePathS))
	titleMap := sync.Map{}
	titleExcelSheet := utils.ExcelSheet{Name: "表头汇总"}
	returnExcelSheet := utils.ExcelSheet{Name: "退款数据"}
	// 并发从文件提取数据
	for _, item := range obj.filePathS {
		go func(filepath string) {
			defer func() {
				wait.Done()
				obj.HandledFiles <- filepath
			}()
			fileTip, res, err := readFileData(obj.DirPath, filepath)
			if nil != err {
				obj.Errors = append(obj.Errors, err)
				return
			}
			if len(res) > 1 {
				lock.Lock()
				defer lock.Unlock()
				for index, row := range res {
					copyRow := make([]interface{}, 1, len(row)+1)
					copyRow[0] = fileTip
					for _, item := range row {
						copyRow = append(copyRow, item)
					}
					if index == 0 {
						titleMap.Store(utils.MultiStringMD5Hash(row...), copyRow)
					}
					returnExcelSheet.Content = append(returnExcelSheet.Content, copyRow)
				}
			}
		}(item)
	}
	wait.Wait()
	// 数据处理
	titleMap.Range(func(key, value interface{}) bool {
		titleExcelSheet.Content = append(titleExcelSheet.Content, value.([]interface{}))
		return true
	})
	if len(returnExcelSheet.Content) > 0 {
		sort.Slice(returnExcelSheet.Content, func(i, j int) bool {
			munA, _ := strconv.Atoi(returnExcelSheet.Content[i][0].(string))
			munB, _ := strconv.Atoi(returnExcelSheet.Content[j][0].(string))
			return munA > munB
		})
	}
	obj.ExcelSheetS = append(obj.ExcelSheetS, titleExcelSheet)
	obj.ExcelSheetS = append(obj.ExcelSheetS, returnExcelSheet)
}

func (obj *OrderReturn) SaveExcel() {
	if len(obj.Errors) > 0 {
		return
	}
	fp, err := utils.SafeMakeExcelFp(obj.ExcelSheetS...)
	if nil != err {
		obj.Errors = append(obj.Errors, err)
		return
	}
	obj.ExcelSavePath = obj.DirPath + "/AMZ放款表-数据提取-退款订单.xlsx"
	err = fp.SaveAs(obj.ExcelSavePath)
	if nil != err {
		obj.Errors = append(obj.Errors, err)
	}
}
