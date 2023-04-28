package test

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

var LogChan = make(chan string)
var Timeout = 5 * time.Second
var pid = -1

// test put
func TestLab1_Put(t *testing.T) {
	go startServer("./output/kvserver")

	time.Sleep(1 * time.Second)

	cmd := exec.Command("./output/kvclient", "put", "hello", "world")

	go ListenStdout(cmd)

	err := cmd.Run()
	if err != nil {
		fmt.Println("failed to begin the clientserver!")
	}

	select {
	case v := <-LogChan:
		fmt.Println("res:" + v)
		KillCommand()
		return
	case <-time.After(3 * time.Second):
		fmt.Println("server reponse time out")
		KillCommand()
	}

}

// test get
func TestLab1_Get(t *testing.T) {
	go startServer("./output/kvserver")

	time.Sleep(1 * time.Second)

	cmd := exec.Command("./output/kvclient", "get", "hello")

	go ListenStdout(cmd)

	err := cmd.Start()
	if err != nil {
		fmt.Println("failed to begin the clientserver!")
	}

	select {
	case v := <-LogChan:
		fmt.Println("res:" + v)
		KillCommand()
		return

	case <-time.After(3 * time.Second):
		fmt.Println("server reponse time out")
		KillCommand()
	}

	cmd.Wait()
}

// test lab1 all
func TestLab1_ALL(t *testing.T) {
	go startServer("./output/kvserver")

	time.Sleep(1 * time.Second)

	cmd := exec.Command("./output/kvclient", "put", "hello", "world")

	err := cmd.Start()
	if err != nil {
		fmt.Println("failed to begin the clientserver!")
	}

	cmd = exec.Command("./output/kvclient", "get", "hello")

	go ListenStdout(cmd)

	err = cmd.Start()
	if err != nil {
		fmt.Println("failed to begin the clientserver!")
	}

	select {
	case v := <-LogChan:
		fmt.Println("res:" + v)
		KillCommand()
		return

	case <-time.After(3 * time.Second):
		fmt.Println("server reponse time out")
		KillCommand()
	}
	cmd.Wait()
}

/*
start kvserver
*/
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

/*
start kvserver
*/
func KillCommand() {
	if pid != -1 {
		//杀死进程
		prc := exec.Command("kill", "-9", strconv.Itoa(pid))
		out, err := prc.Output()
		if err != nil {
			fmt.Printf("kill proc with port %v failed\n", pid)
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
