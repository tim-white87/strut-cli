package command

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/fatih/color"
)

var cwd, _ = os.Getwd()

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
	defer cmd.Wait()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err := cmd.Start()
	if err != nil {
		color.Red(err.Error())
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
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

// SpawnMapGroup Spawns a group of sync processes
func SpawnMapGroup(mapCmds map[string][]*exec.Cmd) {
	wg := &sync.WaitGroup{}
	wg.Add(len(mapCmds))
	defer wg.Wait()
	for _, cmds := range mapCmds {
		go func(cmds []*exec.Cmd, wg *sync.WaitGroup) {
			defer wg.Done()
			for _, cmd := range cmds {
				swg := &sync.WaitGroup{}
				swg.Add(1)
				defer swg.Wait()
				go Spawn(cmd, swg)
			}
		}(cmds, wg)
	}
}

// Spawn spawn a child process
func Spawn(cmd *exec.Cmd, wg *sync.WaitGroup) {
	defer wg.Done()
	defer cmd.Wait()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
	}
}
