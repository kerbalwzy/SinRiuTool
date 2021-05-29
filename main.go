package main

import (
	"github.com/kerbalwzy/SinRiuTool/goApp/server"
	"github.com/kerbalwzy/SinRiuTool/goApp/utils"
	"io"
	"log"
	"os"
)

func init() {
	logger := utils.GetLogger()
	logger.SetLevel(utils.Debug)
	logger.SetOutput(io.MultiWriter(
		os.Stdout,
		utils.NewRotateFileWriter("goApp.log", "./", 3, 1024*1024*100),
	))
	logger.SetPrefix("GoApp: ")
	logger.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func main() {
	server.StartLocalServer()
}
