package utils

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func RunCommandWithLog(logPath, filePath string, args ...string) {
	log.Info(filePath, args)
	cmd := exec.Command(filePath, args...)
	f, err := os.Create(logPath)
	if err != nil {
		log.Error(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	//stderr, _ := cmd.StderrPipe()
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Error(err)
	}

	defer stdoutPipe.Close()

	if err = cmd.Start(); err != nil {
		log.Error(err)
	}

	scanner := bufio.NewScanner(stdoutPipe)
	for scanner.Scan() {
		m := scanner.Text()
		//log.Debug(m)
		fmt.Fprintln(w, m)
		w.Flush()
	}

	if err = cmd.Wait(); err != nil {
		log.Error(err)
	}
}
