package service

import (
	"dolphin-sandbox/biz/response"
	"time"

	"dolphin-sandbox/biz/pack/runner/nodejs"
	runnertypes "dolphin-sandbox/biz/pack/runner/types"
	"dolphin-sandbox/biz/pack/static"
)

func RunNodeJsCode(code string, preload string, options *runnertypes.RunnerOptions) *response.Response {
	if err := checkOptions(options); err != nil {
		return response.ErrorResponse(-400, err.Error())
	}

	if !static.GetSandboxGlobalConfigurations().EnablePreload {
		preload = ""
	}

	timeout := time.Duration(
		static.GetSandboxGlobalConfigurations().WorkerTimeout * int(time.Second),
	)

	runner := nodejs.NodeJsRunner{}
	stdout, stderr, done, err := runner.Run(code, timeout, nil, preload, options)
	if err != nil {
		return response.ErrorResponse(-500, err.Error())
	}

	stdoutStr := ""
	stderrStr := ""

	defer close(done)
	defer close(stdout)
	defer close(stderr)

	for {
		select {
		case <-done:
			return response.SuccessResponse(&RunCodeResponse{
				Stdout: stdoutStr,
				Stderr: stderrStr,
			})
		case out := <-stdout:
			stdoutStr += string(out)
		case err := <-stderr:
			stderrStr += string(err)
		}
	}
}
