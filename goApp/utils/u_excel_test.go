package utils

import "testing"

func TestSheetDataContainer(t *testing.T) {
	sheet := SheetDataContainer{
		SheetName: "测试",
	}
	t.Log(sheet.Len())
	sheet.AddRow([]interface{}{1, 2})
	sheet.AddRow([]interface{}{3, 4})
	sheet.AddRow([]interface{}{"hello", 2})
	sheet.AddRow([]interface{}{1, 2})
	sheet.AddRow([]interface{}{1, 2})
	t.Log(sheet.Data())
	t.Log(sheet.ToSheetData())
}

func TestMakeExcelFp(t *testing.T) {
	sheet := SheetDataContainer{
		SheetName: "测试",
	}
	sheet.AddRow([]interface{}{1, 2})
	sheet.AddRow([]interface{}{3, 4})
	sheet.AddRow([]interface{}{"hello", 2})
	sheet.AddRow([]interface{}{1, 2})
	sheet.AddRow([]interface{}{1, 2})
	fp, err := MakeExcelFp(sheet.ToSheetData()...)
	if nil != err {
		t.Fatal(err)
	} else {
		t.Log(&fp)
		fp.SaveAs("u_excel_test.xlsx")
	}
}
