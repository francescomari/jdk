package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: jdk [version] [command...]")
		os.Exit(1)
	}

	javaHome, err := readJavaHome(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid Java version.")
		os.Exit(1)
	}

	if err := runProcess(javaHome, os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to start process.")
		os.Exit(1)
	}

	fmt.Fprintln(os.Stderr, "Invalid process state.")
	os.Exit(1)
}

func normalizeJavaVersion(version string) string {
	switch version {
	case "8":
		return "1.8"
	default:
		return version
	}
}

func readJavaHome(version string) (string, error) {
	cmd := exec.Command("/usr/libexec/java_home", "-v", normalizeJavaVersion(version))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("capture stdout: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("start: %v", err)
	}

	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", fmt.Errorf("read stdout: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("wait: %v", err)
	}

	if code := cmd.ProcessState.ExitCode(); code != 0 {
		return "", fmt.Errorf("invalid exit code: %v", code)
	}

	return strings.TrimSpace(string(data)), nil
}

func runProcess(javaHome string, commandLine []string) error {
	command, err := exec.LookPath(commandName(commandLine))
	if err != nil {
		return err
	}

	env := os.Environ()
	env = append(env, fmt.Sprintf("JAVA_HOME=%s", javaHome))

	return syscall.Exec(command, commandLine, env)
}

func commandName(commandLine []string) string {
	if len(commandLine) < 1 {
		panic("invalid command line")
	}
	return commandLine[0]
}
