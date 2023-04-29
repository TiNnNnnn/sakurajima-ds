package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/shirou/gopsutil/process"
)

var LogChan chan string
var pid = -1

func main() {
	//testLab1_Put()
	testLab1_Get()
}

func Test_test(t *testing.T) {
	stm := make(map[string]string)
	k := stm["hah"]
	fmt.Println(k)
}

// test put
func testLab1_Put() {

	LogChan = make(chan string)
	defer KillCommand()

	go startServer("./output/kvserver")
	time.Sleep(1 * time.Second)
	cmd := exec.Command("./output/kvclient", "put", "hello", "world")

	go ListenStdout(cmd)

	err := cmd.Start()
	if err != nil {
		fmt.Println("failed to begin the clientserver!")
	}

	select {
	case v := <-LogChan:
		fmt.Println("client output:" + v)
		if v == "success" {
			fmt.Println("pass lab1 put")
		}
		return
	case <-time.After(6 * time.Second):
		fmt.Println("server reponse time out")
		return
	}
}

// test lab1 all
func testLab1_Get() {
	LogChan = make(chan string)
	defer KillCommand()

	go startServer("./output/kvserver")
	time.Sleep(1 * time.Second)

	cmd1 := exec.Command("./output/kvclient", "put", "hello", "world")
	err1 := cmd1.Run()
	if err1 != nil {
		fmt.Println("failed to begin the clientserver!")
	}

	time.Sleep(1 * time.Second)

	cmd2 := exec.Command("./output/kvclient", "get", "hello")
	go ListenStdout(cmd2)
	err2 := cmd2.Start()
	if err2 != nil {
		fmt.Println("failed to begin the clientserver!")
	}

	select {
	case v := <-LogChan:
		fmt.Println("client output:" + v)
		if v == "world" {
			fmt.Println("pass lab1 get")
		}
		return
	case <-time.After(4 * time.Second):
		fmt.Println("server reponse time out")
		return
	}
}

func startServer(name string, arg ...string) {
	cmd := exec.Command(name)
	err := cmd.Start()
	if err != nil {
		fmt.Println("failed to begin the kvserver!")
	}
	pid = cmd.Process.Pid
	fmt.Printf("get the server pid: %v\n", pid)
	cmd.Wait()
}

func KillCommand() {
	if pid != -1 {
		prc := exec.Command("kill", "-9", strconv.Itoa(pid))
		out, err := prc.Output()
		if err != nil {
			fmt.Printf("kill proc with pid %v failed\n", pid)
			return
		}
		fmt.Printf("kill proc with pid %v success! %v\n", pid, string(out))
	}
}

func ListenStdout(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	go func() {
		for {
			l, err := out.ReadBytes('\n')
			if err != nil && err.Error() != "EOF" {
				log.Print(err)
				continue
			}
			if len(l) == 0 {
				continue
			}
			LogChan <- string(l)
		}
	}()
}

func GetNameBypId(pid int) string {
	pn, err := process.NewProcess(int32(pid))
	if err != nil {
		log.Println("process.NewProcess failed")
		return ""
	}
	pName, err := pn.Name()
	if err != nil {
		log.Print("get pn Name() failed")
		return ""
	}
	return pName
}
