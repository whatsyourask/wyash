package shell

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strings"
	"os/exec"
	"os/user"
	"runtime"
	// "reflect"
)

var prompt = shellPromptCreate()
var prevPwd string = shellGetPwd()

func ShellLoop() {
	const exit = "exit"
	const quit = "quit"
	execStatus := true
	for execStatus {
		fmt.Print(prompt)
		cmd := shellReadCmd()
		cmdArgs := shellGetArgs(cmd)
		shellRun(cmdArgs)
	}
}

func shellPromptCreate() string {
	user := shellGetUser()
	hostname := shellGetHostname()
	pwd := shellGetPwd()
	homedir := shellGetHomeDir()
	pwd = strings.Replace(pwd, homedir, "~", 1)
	prompt := "# " + user + " @ " + hostname + " # - [" + pwd + "]\n~ $ "
	return prompt
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
	"exit": shellExit,
	"quit": shellExit,
}

func shellRun(cmdArgs []string) {
	for builtin, builtin_func := range builtins {
		if builtin == cmdArgs[0] {
			builtin_func(cmdArgs)
			return
		}
	}
	shellExecCmd(cmdArgs)
}

func shellGetHomeDir() string {
	homedir, _ := os.UserHomeDir()
	homedirStr := string(homedir)
	return homedirStr
}

func shellGetPwd() string {
	pwd, _ := os.Getwd()
	pwdStr := string(pwd)
	return pwdStr
}

func shellGetUser() string {
	user, _ := user.Current()
	userStr := user.Username
	return userStr
}

func shellGetHostname() string {
	hostname, _ := os.Hostname()
	hostnameStr := string(hostname)
	return hostnameStr
}

func shellCd(cmdArgs []string) {
	cmdArgsLen := len(cmdArgs)
	if cmdArgsLen < 2 {
		os.Chdir("..")
	} else {
		switch cmdArgs[1] {
		case "~":
			cmdHome, _ := os.UserHomeDir()
			os.Chdir(cmdHome)
		case "-":
			fmt.Println(prevPwd)
			os.Chdir(prevPwd)
		default:
			os.Chdir(cmdArgs[1])
		}
	}
	prompt = shellPromptCreate()
}

func shellExecCmd(cmdArgs []string) {
	executable := cmdArgs[0]
	execArgs := cmdArgs[1:]
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

func shellExit(cmdArgs []string) {
	os.Exit(0)
}