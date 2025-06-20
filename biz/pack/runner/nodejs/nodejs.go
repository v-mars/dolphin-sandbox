package nodejs

import (
	"dolphin-sandbox/biz/pack/runner"
	"dolphin-sandbox/biz/pack/runner/types"
	"dolphin-sandbox/biz/pack/static"
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

type NodeJsRunner struct {
	runner.TempDirRunner
}

//go:embed prescript.js
var nodejsSandboxFs []byte

var (
	RequiredFs = []string{
		path.Join(LibPath, ProjectName, "node_temp"),
		path.Join(LibPath, LibName),
		"/etc/ssl/certs/ca-certificates.crt",
		"/etc/nsswitch.conf",
		"/etc/resolv.conf",
		"/run/systemd/resolve/stub-resolv.conf",
		"/etc/hosts",
	}
)

func (p *NodeJsRunner) Run(
	code string,
	timeout time.Duration,
	stdin []byte,
	preload string,
	options *types.RunnerOptions,
) (chan []byte, chan []byte, chan bool, error) {
	configuration := static.GetSandboxGlobalConfigurations()

	// capture the output
	outputHandler := runner.NewOutputCaptureRunner()
	outputHandler.SetTimeout(timeout)

	err := p.WithTempDir("/", RequiredFs, func(rootPath string) error {
		outputHandler.SetAfterExitHook(func() {
			os.RemoveAll(rootPath)
		})

		// initialize the environment
		scriptPath, err := p.InitializeEnvironment(code, preload, rootPath)
		if err != nil {
			return err
		}

		// create a new process
		cmd := exec.Command(
			static.GetSandboxGlobalConfigurations().NodejsPath,
			scriptPath,
			strconv.Itoa(static.SandboxUserUid),
			strconv.Itoa(static.SandboxGroupId),
			options.Json(),
		)
		cmd.Env = []string{}

		if len(configuration.AllowedSyscalls) > 0 {
			cmd.Env = append(
				cmd.Env,
				fmt.Sprintf("ALLOWED_SYSCALLS=%s", strings.Trim(
					strings.Join(strings.Fields(fmt.Sprint(configuration.AllowedSyscalls)), ","), "[]",
				)),
			)
		}

		// capture the output
		err = outputHandler.CaptureOutput(cmd)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return outputHandler.GetStdout(), outputHandler.GetStderr(), outputHandler.GetDone(), nil
}

func (p *NodeJsRunner) InitializeEnvironment(code string, preload string, root_path string) (string, error) {
	if !checkLibAvaliable() {
		releaseLibBinary()
	}

	nodeSandboxFile := string(nodejsSandboxFs)
	if preload != "" {
		nodeSandboxFile = fmt.Sprintf("%s\n%s", preload, nodeSandboxFile)
	}

	// join nodejs_sandbox_fs and code
	// encode code with base64
	code = base64.StdEncoding.EncodeToString([]byte(code))
	// FIXE: redeclared function causes code injection
	evalCode := fmt.Sprintf("eval(Buffer.from('%s', 'base64').toString('utf-8'))", code)
	code = nodeSandboxFile + evalCode

	// override root_path/tmp/sandbox-nodejs-project/prescript.js
	scriptPath := path.Join(root_path, LibPath, ProjectName, "node_temp/node_temp/test.js")
	err := os.WriteFile(scriptPath, []byte(code), 0755)
	if err != nil {
		return "", err
	}

	return scriptPath, nil
}
