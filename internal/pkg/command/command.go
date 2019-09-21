package command

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/fatih/color"
)

// ExecuteGroup executes a group of external processes
func ExecuteGroup(cmds []*exec.Cmd) {
	executeWG := &sync.WaitGroup{}
	executeWG.Add(len(cmds))
	defer executeWG.Wait()
	for _, cmd := range cmds {
		go Execute(cmd, executeWG)
	}
}

// Execute executes an external process
func Execute(cmd *exec.Cmd, wg *sync.WaitGroup) {
	defer wg.Done()
	binary, lookErr := exec.LookPath(cmd.Args[0])
	if lookErr != nil {
		panic(lookErr)
	}
	// fmt.Println(cmd.Path)
	// os.Chdir(cmd.Path)
	// fmt.Println(os.Getwd())
	env := os.Environ()
	execErr := syscall.Exec(binary, cmd.Args, env)
	if execErr != nil {
		panic(execErr)
	}
}

// SpawnGroup spawns a group of processes
func SpawnGroup(cmds []*exec.Cmd) {
	spawnWg := &sync.WaitGroup{}
	spawnWg.Add(len(cmds))
	defer spawnWg.Wait()
	for _, cmd := range cmds {
		go Spawn(cmd, spawnWg)
	}
}

// Spawn spawn a child process
func Spawn(cmd *exec.Cmd, wg *sync.WaitGroup) {
	defer wg.Done()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		color.Red(err.Error())
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	errScanner := bufio.NewScanner(stderr)
	for errScanner.Scan() {
		color.Red(scanner.Text())
	}
	cmd.Wait()
}
