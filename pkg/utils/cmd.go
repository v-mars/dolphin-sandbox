package utils

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"time"
)

type CmdResp struct {
	ExitCode int    `json:"exitCode"`
	Out      string `json:"out"`
	ErrOut   string `json:"errOut"`
	Stdout   string `json:"stdOut"`
}

// CheckLsExists 预先检查命令是否存在
func CheckLsExists(cmd string) {
	path, err := exec.LookPath(cmd)
	if err != nil {
		fmt.Printf("didn't find 'ls' executable\n")
	} else {
		fmt.Printf("'ls' executable is in '%s'\n", path)
	}
}

func read(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser, stdAndErr, out *string) {
	reader := bufio.NewReader(std)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
			*out = *out + readString
		}
	}
}

// Command 执行命令实时打印
func Command(cmd string, timeout int64, f *os.File) (CmdResp, error) {
	//c := exec.CommandContext(ctx, "cmd", "/C", cmd) // windows
	//fmt.Println("split:", strings.Split(cmd, "\\n"))
	if timeout <= 0 {
		timeout = 10
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancelFunc context.CancelFunc) {
		time.Sleep(time.Duration(timeout) * time.Second)
		cancelFunc()
	}(cancel)
	var result = CmdResp{}
	c := exec.CommandContext(ctx, "bash", "-cxe", cmd) // mac linux
	stdout, err := c.StdoutPipe()
	if err != nil {
		result.ErrOut = err.Error()
		result.ExitCode = c.ProcessState.ExitCode()
		return result, err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		result.ErrOut = err.Error()
		result.ExitCode = c.ProcessState.ExitCode()
		return result, err
	}
	var wg sync.WaitGroup
	// 因为有2个任务, 一个需要读取stderr 另一个需要读取stdout
	wg.Add(2)
	go func() {
		reader := bufio.NewReader(stderr)
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					return
				}
				//fmt.Print(readString)
				result.ErrOut = result.ErrOut + readString
				result.Out = result.Out + readString
				if f != nil {
					_, _ = f.Write([]byte(readString))
				}
			}
		}
	}()

	go func() {
		reader := bufio.NewReader(stdout)
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				readString, err := reader.ReadString('\n')
				if err != nil || err == io.EOF {
					return
				}
				//fmt.Print(readString)
				result.Stdout = result.Stdout + readString
				result.Out = result.Out + readString
				if f != nil {
					_, _ = f.Write([]byte(readString))
				}
			}
		}
	}()
	//go read(ctx, &wg, stderr, &result.ErrOut)
	//go read(ctx, &wg, stdout, &result.Stdout)
	// 这里一定要用start,而不是run 详情请看下面的图
	err = c.Start()
	if err != nil {
		result.ErrOut = err.Error()
		result.ExitCode = c.ProcessState.ExitCode()
		return result, err
	}

	// 等待任务结束
	wg.Wait()
	//fmt.Println("c.ProcessState:", c.ProcessState.String())
	result.ExitCode = c.ProcessState.ExitCode()
	return result, err
}

// CommandOut 执行命令一次性返回结果
func CommandOut(cmd string, timeout int64) (CmdResp, error) {
	//c := exec.CommandContext(ctx, "cmd", "/C", cmd) // windows
	//fmt.Println("split:", strings.Split(cmd, "\\n"))
	if timeout <= 0 {
		timeout = 10
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func(cancelFunc context.CancelFunc) {
		time.Sleep(time.Duration(timeout) * time.Second)
		cancelFunc()
	}(cancel)
	var result = CmdResp{}
	c := exec.CommandContext(ctx, "bash", "-cxe", cmd) // mac linux

	//output, err := c.CombinedOutput()
	//if err != nil {
	//	result.ErrOut = err.Error()
	//	result.ExitCode = c.ProcessState.ExitCode()
	//	return result, err
	//}
	//result.Stdout = string(output)

	c.Stdout = os.Stdout
	//c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		result.ErrOut = err.Error()
		result.ExitCode = c.ProcessState.ExitCode()
		return result, err
	}

	result.ExitCode = c.ProcessState.ExitCode()
	return result, err
}

// RunShell
// run shell code
func RunShell(ctx context.Context, scriptCode string) (*exec.Cmd, string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}
	tmpFile, err := ioutil.TempFile("", "*.sh")
	if err != nil {
		return nil, "", err
	}
	shellCodePath := tmpFile.Name()
	_, err = tmpFile.WriteString(scriptCode)
	if err != nil {
		return nil, "", err
	}

	_ = tmpFile.Sync()
	err = tmpFile.Close()
	if err != nil {
		return nil, "", err
	}
	cmd := exec.CommandContext(ctx, shell, shellCodePath)
	return cmd, shellCodePath, nil
}

// RunShellPipe run shell
func RunShellPipe(ctx context.Context, cmdIn string) (out io.ReadCloser) {
	pr, pw := io.Pipe()
	go func() {
		var (
			exitCode = -1
			err      error
			codePath string
			cmd      *exec.Cmd
		)
		defer func(pw *io.PipeWriter) { _ = pw.Close() }(pw)
		defer func() {
			now := time.Now().Local().Format("2006-01-02 15:04:05")
			//finishTxt := fmt.Sprintf("#################### Bash Shell Finish ####################")
			_, _ = pw.Write([]byte(fmt.Sprintf("\n%s\n%s Shell Run Finished, Return exitCode:%5d", "", now, exitCode))) // write exitCode,total 5 byte
			if codePath != "" {
				_ = os.Remove(codePath)
			}
		}()
		cmd, codePath, err = RunShell(ctx, cmdIn)
		if err != nil {
			_, _ = pw.Write([]byte(err.Error()))
			return
		}
		cmd.Stdout = pw
		cmd.Stderr = pw
		//startTxt := fmt.Sprintf("#################### Bash Shell Start ####################\n\n")
		//_, _ = pw.Write([]byte(startTxt))
		err = cmd.Start()
		if err != nil {
			_, _ = pw.Write([]byte(err.Error()))
			return
		}

		err = cmd.Wait()
		if err != nil {
			// try to get the exit code
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			}
			// deal err
			// if context err,will change err to custom msg
			switch ctx.Err() {
			case context.DeadlineExceeded:
				_, err = pw.Write([]byte("ErrCtxDeadlineExceeded"))
				if err != nil {
					return
				}
			case context.Canceled:
				_, _ = pw.Write([]byte("ErrCtxCanceled"))
			default:
				_, err = pw.Write([]byte(err.Error()))
				if err != nil {
					return
				}
			}
		} else {
			exitCode = 0
		}

	}()
	return pr
}
