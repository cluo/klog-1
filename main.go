package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var (
	name  string // app 名称
	podNO int    // 第几个pod

	logOption string // 查看日志选项
)

func init() {
	flag.StringVar(&name, "name", "", "app name")
	flag.IntVar(&podNO, "podNO", 1, "spec witch pod, starting from 1")
	flag.StringVar(&logOption, "logOption", "-f --tail=100", "kubectl logs options")
	flag.Parse()

	if name == "" {
		fmt.Println("error: name is empty.")
		fmt.Println("usage: klog --name=<name> [--podNO=podNO]\n")
		fmt.Println("example：")
		fmt.Println("the first pod log:\n    klog --name=account")
		fmt.Println("spec pod log:\n    klog --name=account --podNO=3\n")

		fmt.Println(`view some modules log, copy and execute:
	klog --name=wq-account
	klog --name=wq-email
	klog --name=wq-exchange-order
	klog --name=wq-file
	klog --name=wq-history-quote
	klog --name=wq-order-prepare
	klog --name=wq-pcc
	klog --name=wq-quote
	klog --name=wq-risk`, "\n")

		os.Exit(1)
	}

	if podNO == 0 {
		podNO = 1
	}

	if logOption == "" {
		logOption = "--tail=100 -f"
	}
}

func main() {
	// 获取pod name
	arg := fmt.Sprintf("kubectl get pod | grep %s | sed -n \"%d, 1p\" | awk '{print $1}'", name, podNO)
	out, err := execShellForNoneBlock(arg)
	if err != nil {
		fmt.Printf("error: %s, arg = %s\n", err.Error(), arg)
		return
	}
	podName := getData(out)

	// 获取pod name失败处理
	if podName == "" {
		arg = fmt.Sprintf("kubectl get pod | grep %s | wc -l", name)
		out, err = execShellForNoneBlock(arg)
		if err != nil {
			fmt.Printf("error: %s, arg = %s\n", err.Error(), arg)
			return
		}
		fmt.Printf("error: can not find the %dth pod name, total %s pods.\n\n", podNO, getData(out))

		arg = fmt.Sprintf("kubectl get pod | grep %s", name)
		out, err := execShellForNoneBlock(arg)
		if err != nil {
			fmt.Printf("error: %s, arg = %s\n", err.Error(), arg)
			return
		}
		fmt.Println(getData(out), "\n")
		return
	}

	// 查看日志
	arg = fmt.Sprintf("kubectl logs %s %s", logOption, podName)
	execShellForBlock(arg)
}

// 执行非阻塞shell命令
func execShellForNoneBlock(arg string) ([]byte, error) {
	cmd := exec.Command("/bin/sh", "-c", arg)
	return cmd.Output()
}

// 去掉换行符
func getData(data []byte) string {
	l := len(data)
	if l > 0 {
		return string(data[:l-1])
	}

	return ""
}

// 执行阻塞shell命令
func execShellForBlock(arg string) {
	cmd := exec.Command("/bin/sh", "-c", arg)

	//显示运行的命令
	fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		return
	}

	cmd.Start()
	reader := bufio.NewReader(stdout)

	// 读取内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}

		fmt.Printf(line)
	}

	cmd.Wait()
}
