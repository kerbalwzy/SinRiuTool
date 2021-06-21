package amz

import (
	"fmt"
	"sync"
	"testing"
)

func TestIsReturn(t *testing.T) {
	key := "Remboursement de contestation de prélèvement"
	res := IsReturn(key)
	t.Log(res)
}

func TestOrderReturn(t *testing.T) {
	dirPath := "C:\\Users\\admin\\Desktop\\202001-202103放款数据"
	obj := OrderReturn{DirPath: dirPath}
	obj.Init()
	if len(obj.Errors) == 0 {
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			defer wait.Done()
			obj.ReadFilesData()
		}()
		for i := 0; i < len(obj.filePathS); i++ {
			filePath := <-obj.HandledFiles
			fmt.Println("处理文件完成: ", filePath)
		}
		wait.Wait()
		obj.SaveExcel()
	}
	t.Log(obj.Errors)
}
