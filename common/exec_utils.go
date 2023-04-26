package common

import (
	"fmt"
	"log"
	"os/exec"
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

func GetGroupIdBypId2(pid int, sType string) string {
	cmd := exec.Command("bash", "-c", "netstat -nltp |grep "+strconv.Itoa(pid)+"/"+sType)

	out, err := cmd.Output()
	if err != nil {
		println("error:", err.Error())
	}
	if len(out) == 0 {
		fmt.Printf("out: can't find proc with pid %v\n", pid)
		return ""
	}

	fields := strings.Fields(string(out))
	//println(string(out))

	return fields[3]
}

func GetGroupIdBypId(pid int) string {

	c := cmd.NewCmd("bash", "-c", "netstat -nltp |grep "+strconv.Itoa(pid))
	status := <-c.Start()

	//fmt.Printf("nestat result: %v\n", status.Stdout)
	if len(status.Stdout) == 0 {
		fmt.Printf("status.Stdout: can't find proc with pid %v\n", pid)
		return ""
	}

	fields := strings.Fields(status.Stdout[0])
	if len(fields) == 0 {
		fmt.Printf("Fields: can't find proc with pid %v\n", pid)
		return ""
	}
	addr := fields[3]
	fmt.Println(addr)

	return addr
}

func GetPidByport(port string) string {
	c := cmd.NewCmd("bash", "-c", "netstat -nltp |grep "+port)
	status := <-c.Start()
	fmt.Printf("nestat result: %v", status.Stdout)

	if len(status.Stdout) == 0 {
		fmt.Printf("can't find proc with port %v", port)
		//w.Write([]byte("can't find proc with port " + port))
		return ""
	}

	fields := strings.Fields(status.Stdout[0])
	if len(fields) == 0 {
		fmt.Printf("can't find proc with port %v", port)
		//w.Write([]byte("can't find proc with port " + port))
		return ""
	}
	//取出进程pid
	pid := strings.Split(fields[6], "/")[0]

	return pid
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
