package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var GitExecutable string

func RunCommand(dir string, wg *sync.WaitGroup) {
    defer wg.Done()
    cmd := &exec.Cmd{
        Path: GitExecutable,
        Dir: dir,
        Args: os.Args,
    }
    log.Printf("command:%s dir:%s\n", cmd.String(), dir)

    if output, err := cmd.CombinedOutput(); err != nil {
        log.Printf("command error: %+v\n", err)
    } else {
        log.Printf("Output: %+v\n", output)
    }
}

func main() {
    var err error
    GitExecutable, err = exec.LookPath("git")
    if err != nil {
        log.Fatalf("git executable can't be found in path.")
        return
    }

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Can't get current working directory. error is: %+v", err)
		return
	}

	directories := []string{}
	if _, err := os.Stat(strings.Join([]string{pwd, ".git"}, "/")); os.IsNotExist(err) {
		files, err := ioutil.ReadDir(pwd)
		if err != nil {
			log.Fatalf("Can't read file of current working directory. error is: %+v", err)
		}

		for _, file := range files {
			gitPathArr := []string{pwd, file.Name(), ".git"}
			if _, err := os.Stat(strings.Join(gitPathArr, "/")); os.IsNotExist(err) {

			} else {

                gitPath := strings.Join([]string{pwd, file.Name()}, "/")
				directories = append(directories, gitPath)
			}
		}
	} else {
		log.Println("Will work on current directory: " + pwd)
		directories = append(directories, pwd)
	}

    wg := &sync.WaitGroup{}
    for _, dir := range directories {
        wg.Add(1)
        go RunCommand(dir, wg);
    }
    wg.Wait()
    log.Printf("%+v", directories)
}
