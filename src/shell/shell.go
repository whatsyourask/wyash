package shell

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strings"
	"os/exec"
	"runtime"
)

func ShellLoop() {
	const cmdPrompt = "wyash$ "
	cmdExecStatus := true
	for cmdExecStatus {
		fmt.Print(cmdPrompt)
		cmd := shellReadCmd()
		cmdArgs := shellGetArgs(cmd)
		cmdExecStatus = shellRun(cmdArgs)
	}
}

func shellReadCmd() string {
	reader := bufio.NewReader(os.Stdin)
	cmd, err := reader.ReadString('\n')
	if err != nil {
		dieWithLog(traceFunc() + ": error of input reading.")
	}
	return cmd
}

func dieWithLog(err string) {
	log.Fatal(err)
	os.Exit(1)
}

func shellGetArgs(cmd string) []string {
	const cmdArgsSize = 64
	cmdArgs := strings.Split(cmd, " ")
	return cmdArgs
}

var builtins = map[string]func([]string)bool{
	"cd": shellCd,
}

func shellRun(cmdArgs []string) bool {
	for builtin, builtin_func := range builtins {
		if builtin == cmdArgs[0] {
			return builtin_func(cmdArgs)
		}
	}
	return shellExecCmd(cmdArgs)
}

func shellCd(cmdArgs []string) bool {
	fmt.Println("CD")
	cmdArgsLen := len(cmdArgs)
	if cmdArgsLen < 2 {
		dieWithLog(traceFunc() + ": argument required.")
	} else {
		os.Chdir(cmdArgs[1])
	}
	return true
}

func shellExecCmd(cmdArgs []string) bool {
	executable := cmdArgs[0]
	fmt.Println(executable)
	execArgs := cmdArgs[0:]
	fmt.Println(execArgs)
	cmd := exec.Command(executable, execArgs...)
	err := cmd.Run()
	if err != nil {
		dieWithLog(traceFunc() + ": error of command execution.")
	}
	return true
}

func traceFunc() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}