package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	// 	golang的时间格式 默认采用的是RFC333 RFC3339     = "2006-01-02T15:04:05Z07:00"
	// 这个格式还不能随便写,要不然生成的文件时间不对 正确的固定的格式:// 这是个奇葩,必须是这个时间点, 据说是go诞生之日, 记忆方法:6-1-2-3-4-5
	//    return []byte(fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))), nil
	//TimeFormat  = "20200601"
	TimeFormat = "20060102"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)
	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filename string) *os.File {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			mkDir()
		} else if os.IsPermission(err) {
			log.Fatalf("Permission : %v", err)
		}
	}

	handler, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("open file failed: %v", err)
	}
	return handler
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
