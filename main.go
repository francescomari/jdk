package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: jdk [version] [command...]\n")
		os.Exit(1)
	}

	javaHome, err := readJavaHome(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid Java version.\n")
		os.Exit(1)
	}

	exitCode, err := runProcess(javaHome, os.Args[2:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start process.\n")
		os.Exit(1)
	}

	os.Exit(exitCode)
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

func runProcess(javaHome string, commandLine []string) (int, error) {
	cmd := exec.Command(processName(commandLine), processArgs(commandLine)...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("JAVA_HOME=%s", javaHome))

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return 0, err
	}

	return cmd.ProcessState.ExitCode(), nil
}

func processName(commandLine []string) string {
	if len(commandLine) < 1 {
		panic("invalid command line")
	}
	return commandLine[0]
}

func processArgs(commandLine []string) []string {
	if len(commandLine) < 2 {
		return nil
	}
	return commandLine[1:]
}
