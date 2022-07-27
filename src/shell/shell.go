package shell

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strings"
	"os/exec"
	"runtime"
	"reflect"
)

func ShellLoop() {
	const cmdPrompt = "wyash$ "
	const exit = "exit\n"
	cmdExecStatus := true
	for cmdExecStatus {
		fmt.Print(cmdPrompt)
		cmd := shellReadCmd()
		if cmd == exit {
			os.Exit(1)
		}
		cmdArgs := shellGetArgs(cmd)
		shellRun(cmdArgs)
	}
}

func shellReadCmd() string {
	reader := bufio.NewReader(os.Stdin)
	cmd, err := reader.ReadString('\n')
	if err != nil {
		dieWithLog(traceFunc() + ": error of input reading.")
	}
	cmdLen := len(cmd)
	cmd = cmd[0:cmdLen - 1]
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

var builtins = map[string]func([]string){
	"cd": shellCd,
}

func shellRun(cmdArgs []string) {
	fmt.Println(traceFunc(), "cmdArgs = ", cmdArgs)
	for builtin, builtin_func := range builtins {
		if builtin == cmdArgs[0] {
			fmt.Println(traceFunc(), "cmdArgs[0] = " + cmdArgs[0])
			builtin_func(cmdArgs)
			return
		}
	}
	shellExecCmd(cmdArgs)
}

func shellCd(cmdArgs []string) {
	fmt.Println(traceFunc(), "CD")
	cmdArgsLen := len(cmdArgs)
	if cmdArgsLen < 2 {
		os.Chdir("..")
	} else {
		os.Chdir(cmdArgs[1])
	}
}

func shellExecCmd(cmdArgs []string) {
	fmt.Println(traceFunc(), "cmdArgs = ", cmdArgs)
	executable := cmdArgs[0]
	// fmt.Println("executable = ", executable)
	// fmt.Println("executable type of ", reflect.TypeOf(executable))
	execArgs := cmdArgs[1:]
	tmpArgs := strings.Join(execArgs, " ")
	// execArgsLen := len(execArgs)
	// execArgs = execArgs[0:execArgsLen - 1]
	fmt.Println(traceFunc(), "execArgs = ", tmpArgs)
	fmt.Println(traceFunc(), "execArgs type of ", reflect.TypeOf(tmpArgs))
	cmdOutput, err := exec.Command(executable, tmpArgs).Output()
	fmt.Println(traceFunc(), string(cmdOutput))
	if err != nil {
		fmt.Println(err)
	}
}

func traceFunc() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	funcName := f.Name()
	funcNameSplited := strings.Split(funcName, "/")
	shortFuncName := funcNameSplited[len(funcNameSplited) - 1]
	return shortFuncName
}