package shell

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strings"
	"os/exec"
	"runtime"
	// "reflect"
)

func ShellLoop() {
	cmdPrompt := shellPromptCreate()
	const exit = "exit"
	cmdExecStatus := true
	for cmdExecStatus {
		fmt.Print(cmdPrompt)
		cmd := shellReadCmd()
		if cmd == exit {
			os.Exit(0)
		}
		cmdArgs := shellGetArgs(cmd)
		shellRun(cmdArgs)
	}
}

func shellPromptCreate() string {
	cmdWhoami, _ := exec.Command("whoami").Output()
	cmdHostname, _ := exec.Command("hostname").Output()
	cmdWhoamiStr := rmLastChrUint8(cmdWhoami)
	cmdHostnameStr := rmLastChrUint8(cmdHostname)
	cmdPrompt := cmdWhoamiStr + "|" + cmdHostnameStr + "~ "
	return cmdPrompt
}

func rmLastChrUint8(str []uint8) string {
	return string(str[0:len(str) - 1])
}

func rmLastChrStr(str string) string {
	return str[0:len(str) - 1]
}

func shellReadCmd() string {
	reader := bufio.NewReader(os.Stdin)
	cmd, err := reader.ReadString('\n')
	if err != nil {
		dieWithLog(traceFunc() + ": error of input reading.")
	}
	cmdStr := rmLastChrStr(cmd)
	return cmdStr
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
	// fmt.Println(traceFunc(), "cmdArgs = ", cmdArgs)
	for builtin, builtin_func := range builtins {
		if builtin == cmdArgs[0] {
			// fmt.Println(traceFunc(), "cmdArgs[0] = " + cmdArgs[0])
			builtin_func(cmdArgs)
			return
		}
	}
}

func shellCd(cmdArgs []string) {
	// fmt.Println(traceFunc(), "CD")
	cmdArgsLen := len(cmdArgs)
	if cmdArgsLen < 2 {
		os.Chdir("..")
	} else {
		os.Chdir(cmdArgs[1])
	}
}

func shellExecCmd(cmdArgs []string) {
	// fmt.Println(traceFunc(), "cmdArgs = ", cmdArgs)
	executable := cmdArgs[0]
	// fmt.Println("executable = ", executable)
	// fmt.Println("executable type of ", reflect.TypeOf(executable))
	execArgs := cmdArgs[1:]
	// execArgsLen := len(execArgs)
	// execArgs = execArgs[0:execArgsLen - 1]
	// fmt.Println(traceFunc(), "execArgs = ", execArgs)
	// fmt.Println(traceFunc(), "execArgs type of ", reflect.TypeOf(execArgs))
	cmd := exec.Command(executable, execArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmdOutput, err := cmd.Output()
	fmt.Println(string(cmdOutput))
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