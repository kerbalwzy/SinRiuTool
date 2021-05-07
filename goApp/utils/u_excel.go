package utils

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"regexp"
	"strconv"
	"sync"
)

var ExcelIllegalCharactersRe = regexp.MustCompile(`[\000-\010]|[\013-\014]|[\016-\037]`)

//Excel单个sheet最多只能由1048576行,超出的行数据将保存到复制了名称的sheet
const ExcelMaxRowCount = 1048576

//const ExcelMaxRowCount = 2

type SheetData struct {
	Name    string
	Content [][]interface{}
}

type SheetDataContainer struct {
	SheetName  string
	dataSheets [][][]interface{} // Excel单个sheet最多只能由1048576行,超出的行数据将保存到复制了名称的sheet
	rawsLen    int               // 数据行数
	mutex      sync.Mutex        // 互斥锁
}

func (obj *SheetDataContainer) Len() int {
	n := 0
	for _, tempSheet := range obj.dataSheets {
		n += len(tempSheet)
	}
	return n
}

func (obj *SheetDataContainer) AddRow(row []interface{}) {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	if obj.rawsLen == 0 || obj.rawsLen+1 > ExcelMaxRowCount {
		tempSheet := make([][]interface{}, 0)
		obj.dataSheets = append(obj.dataSheets, tempSheet)
		obj.rawsLen = 0
	}
	tempSheet := &obj.dataSheets[len(obj.dataSheets)-1]
	*tempSheet = append(*tempSheet, row)
	obj.rawsLen += 1
}

func (obj *SheetDataContainer) Data() [][]interface{} {
	res := make([][]interface{}, 0)
	for _, tempSheet := range obj.dataSheets {
		res = append(res, tempSheet...)
	}
	return res
}

func (obj *SheetDataContainer) ToSheetData() []SheetData {
	res := make([]SheetData, 0)
	for index, tempSheet := range obj.dataSheets {
		if index == 0 {
			res = append(res, SheetData{Name: obj.SheetName, Content: tempSheet})
		} else {
			res = append(res, SheetData{Name: obj.SheetName + strconv.Itoa(index), Content: tempSheet})
		}
	}
	return res
}

func MakeExcelFp(sheetData ...SheetData) (*excelize.File, error) {
	fp := excelize.NewFile()
	for index, item := range sheetData {
		if index == 0 {
			fp.SetSheetName("Sheet1", item.Name)
		} else {
			fp.NewSheet(item.Name)
		}
		streamWriter, err := fp.NewStreamWriter(item.Name)
		if nil != err {
			fmt.Println(err)
			return nil, err
		}
		for rowIndex, row := range item.Content {
			cell, _ := excelize.CoordinatesToCellName(1, rowIndex+1)
			err = streamWriter.SetRow(cell, row)
			if nil != err {
				return nil, err
			}
		}
		err = streamWriter.Flush()
		if nil != err {
			return nil, err
		}
	}
	return fp, nil
}
