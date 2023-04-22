package tinnraft

import (
	_ "fmt"
	"log"
	_ "runtime"
	_ "time"
)

// Debugging
const Debug = true

func DLog(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
		// 获取当前函数的文件名和行号
	}
	return
}

// func DLog(msg string) {
// 	if Debug {
// 		log.Printf("%s %s \n", time.Now().Format("2023-03-20 22:23:24"), msg)
// 	}
// }

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
