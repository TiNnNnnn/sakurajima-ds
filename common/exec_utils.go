package common

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-cmd/cmd"
	"github.com/shirou/gopsutil/process"
)

func GetNameBypId(pid int) string {
	pn, err := process.NewProcess(int32(pid))
	if err != nil {
		log.Println("process.NewProcess failed")
		return ""
	}
	pName, err := pn.Name()
	if err != nil {
		log.Print("pn.Name() failed")
		return ""
	}
	return pName
}

// 查找指定进程
func IsProcessExist(port string) bool {

	c := cmd.NewCmd("bash", "-c", "netstat -nltp |grep 10055")
	<-c.Start()
	fmt.Println(c.Status().Stdout)

	//out := c.Status().Stdout
	if len(c.Status().Stdout) == 0 {
		fmt.Printf("can't find proc with port %v", port)
		return false
	}

	//fmt.Println(fields[6])
	fields := strings.Fields(c.Status().Stdout[0])
	if len(fields) == 0 {
		fmt.Printf("can't find proc with port %v", port)
		return false
	}

	fmt.Println(fields[6])

	//strings.Split(fields[6], "/")
	fmt.Println(strings.Split(fields[6], "/")[0])
	// for _, v := range fields {
	// 	if v == "312144/api_server" {
	// 		fmt.Println(v)
	// 		return true
	// 	}
	// }
	return false
}

func GetPid(addr int64) (pid string, err error) {

	port := fmt.Sprintf("%04X", addr)
	SocketId, err := GetSocketId(port)
	if err != nil {

		return
	}

	SocketInfo := fmt.Sprintf("socket:[%s]", SocketId)

	procDirList, err := ioutil.ReadDir("/proc")
	if err != nil {
		return
	}

	for _, procDir := range procDirList {
		_, err := strconv.Atoi(procDir.Name())
		if err != nil {
			continue
		}
		fdDir := fmt.Sprintf("/proc/%s/fd", procDir.Name())
		fdSonDirList, err := ioutil.ReadDir(fdDir)
		for _, socketFile := range fdSonDirList {

			socket := fmt.Sprintf("/proc/%s/fd/%s", procDir.Name(), socketFile.Name())
			data, err := os.Readlink(socket)
			if err != nil {
				continue
			}
			if SocketInfo == data {
				return procDir.Name(), nil
			}
		}
	}
	return "", errors.New("get pid fail")
}

func GetSocketId(port string) (SocketId string, err error) {
	fi, err := os.Open("/proc/net/tcp")
	if err != nil {
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		Info := strings.Fields(string(a))
		// Info不同位置代表不同，想从哪查就获取啥
		remPort := strings.Split(Info[1], ":")
		if len(remPort) == 2 {
			if remPort[1] == port {
				SocketId = Info[9]
				return
			}
		}
	}
	return "", errors.New("get SocketId fail")

}
