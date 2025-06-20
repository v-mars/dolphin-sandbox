package service

import (
	"dolphin-sandbox/biz/response"
	"time"

	"dolphin-sandbox/biz/pack/runner/python"
	runnertypes "dolphin-sandbox/biz/pack/runner/types"
	"dolphin-sandbox/biz/pack/static"
)

type RunCodeResponse struct {
	Stderr   string `json:"error"`
	Stdout   string `json:"stdout"`
	ExitCode int    `json:"exit_code"`
}

func RunPython3Code(code string, preload string, options *runnertypes.RunnerOptions) *response.Response {
	if err := checkOptions(options); err != nil {
		return response.ErrorResponse(-400, err.Error())
	}

	if !static.GetSandboxGlobalConfigurations().EnablePreload {
		preload = ""
	}

	timeout := time.Duration(
		static.GetSandboxGlobalConfigurations().WorkerTimeout * int(time.Second),
	)

	runner := python.PythonRunner{}
	stdout, stderr, done, err := runner.Run(
		code, timeout, nil, preload, options,
	)
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

type ListDependenciesResponse struct {
	Dependencies []runnertypes.Dependency `json:"dependencies"`
}

func ListPython3Dependencies() *response.Response {
	return response.SuccessResponse(&ListDependenciesResponse{
		Dependencies: python.ListDependencies(),
	})
}

type RefreshDependenciesResponse struct {
	Dependencies []runnertypes.Dependency `json:"dependencies"`
}

func RefreshPython3Dependencies() *response.Response {
	return response.SuccessResponse(&RefreshDependenciesResponse{
		Dependencies: python.RefreshDependencies(),
	})
}

type UpdateDependenciesResponse struct{}

func UpdateDependencies() *response.Response {
	err := python.PreparePythonDependenciesEnv()
	if err != nil {
		return response.ErrorResponse(-500, err.Error())
	}

	return response.SuccessResponse(&UpdateDependenciesResponse{})
}
