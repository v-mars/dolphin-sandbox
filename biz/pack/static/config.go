package static

import (
	"dolphin-sandbox/biz/pack/schema"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strconv"
	"strings"
)

var SandboxGlobalConfigurations schema.SandboxGlobalConfigurations

func InitConfig(path string) error {
	SandboxGlobalConfigurations = schema.SandboxGlobalConfigurations{}

	// read config file
	configFile, err := os.Open(path)
	if err != nil {
		return err
	}

	defer configFile.Close()

	// parse config file
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&SandboxGlobalConfigurations)
	if err != nil {
		return err
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err == nil {
		SandboxGlobalConfigurations.App.Debug = debug
	}

	maxWorkers := os.Getenv("MAX_WORKERS")
	if maxWorkers != "" {
		SandboxGlobalConfigurations.MaxWorkers, _ = strconv.Atoi(maxWorkers)
	}

	maxRequests := os.Getenv("MAX_REQUESTS")
	if maxRequests != "" {
		SandboxGlobalConfigurations.MaxRequests, _ = strconv.Atoi(maxRequests)
	}

	port := os.Getenv("SANDBOX_PORT")
	if port != "" {
		SandboxGlobalConfigurations.App.Port, _ = strconv.Atoi(port)
	}

	timeout := os.Getenv("WORKER_TIMEOUT")
	if timeout != "" {
		SandboxGlobalConfigurations.WorkerTimeout, _ = strconv.Atoi(timeout)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey != "" {
		SandboxGlobalConfigurations.App.Key = apiKey
	}

	pythonPath := os.Getenv("PYTHON_PATH")
	if pythonPath != "" {
		SandboxGlobalConfigurations.PythonPath = pythonPath
	}

	if SandboxGlobalConfigurations.PythonPath == "" {
		SandboxGlobalConfigurations.PythonPath = "/usr/local/bin/python3"
	}

	pythonLibPath := os.Getenv("PYTHON_LIB_PATH")
	if pythonLibPath != "" {
		SandboxGlobalConfigurations.PythonLibPaths = strings.Split(pythonLibPath, ",")
	}

	if len(SandboxGlobalConfigurations.PythonLibPaths) == 0 {
		SandboxGlobalConfigurations.PythonLibPaths = DEFAULT_PYTHON_LIB_REQUIREMENTS
	}

	pythonPipMirrorUrl := os.Getenv("PIP_MIRROR_URL")
	if pythonPipMirrorUrl != "" {
		SandboxGlobalConfigurations.PythonPipMirrorURL = pythonPipMirrorUrl
	}

	pythonDepsUpdateInterval := os.Getenv("PYTHON_DEPS_UPDATE_INTERVAL")
	if pythonDepsUpdateInterval != "" {
		SandboxGlobalConfigurations.PythonDepsUpdateInterval = pythonDepsUpdateInterval
	}

	// if not set "PythonDepsUpdateInterval", update python dependencies every 30 minutes to keep the sandbox up-to-date
	if SandboxGlobalConfigurations.PythonDepsUpdateInterval == "" {
		SandboxGlobalConfigurations.PythonDepsUpdateInterval = "30m"
	}

	nodejsPath := os.Getenv("NODEJS_PATH")
	if nodejsPath != "" {
		SandboxGlobalConfigurations.NodejsPath = nodejsPath
	}

	if SandboxGlobalConfigurations.NodejsPath == "" {
		SandboxGlobalConfigurations.NodejsPath = "/usr/local/bin/node"
	}

	enableNetwork := os.Getenv("ENABLE_NETWORK")
	if enableNetwork != "" {
		SandboxGlobalConfigurations.EnableNetwork, _ = strconv.ParseBool(enableNetwork)
	}

	enablePreload := os.Getenv("ENABLE_PRELOAD")
	if enablePreload != "" {
		SandboxGlobalConfigurations.EnablePreload, _ = strconv.ParseBool(enablePreload)
	}

	allowedSyscalls := os.Getenv("ALLOWED_SYSCALLS")
	if allowedSyscalls != "" {
		strs := strings.Split(allowedSyscalls, ",")
		ary := make([]int, len(strs))
		for i := range ary {
			ary[i], err = strconv.Atoi(strs[i])
			if err != nil {
				return err
			}
		}
		SandboxGlobalConfigurations.AllowedSyscalls = ary
	}

	if SandboxGlobalConfigurations.EnableNetwork {
		log.Println("network has been enabled")
		socks5Proxy := os.Getenv("SOCKS5_PROXY")
		if socks5Proxy != "" {
			SandboxGlobalConfigurations.Proxy.Socks5 = socks5Proxy
		}

		if SandboxGlobalConfigurations.Proxy.Socks5 != "" {
			log.Println("using socks5 proxy: %s", SandboxGlobalConfigurations.Proxy.Socks5)
		}

		https_proxy := os.Getenv("HTTPS_PROXY")
		if https_proxy != "" {
			SandboxGlobalConfigurations.Proxy.Https = https_proxy
		}

		if SandboxGlobalConfigurations.Proxy.Https != "" {
			log.Println("using https proxy: %s", SandboxGlobalConfigurations.Proxy.Https)
		}

		http_proxy := os.Getenv("HTTP_PROXY")
		if http_proxy != "" {
			SandboxGlobalConfigurations.Proxy.Http = http_proxy
		}

		if SandboxGlobalConfigurations.Proxy.Http != "" {
			log.Println("using http proxy: %s", SandboxGlobalConfigurations.Proxy.Http)
		}
	}
	return nil
}

// GetSandboxGlobalConfigurations avoid global modification, use value copy instead
func GetSandboxGlobalConfigurations() schema.SandboxGlobalConfigurations {
	return SandboxGlobalConfigurations
}

type RunnerDependencies struct {
	PythonRequirements string
}

var runnerDependencies RunnerDependencies

func GetRunnerDependencies() RunnerDependencies {
	return runnerDependencies
}

func SetupRunnerDependencies() error {
	file, err := os.ReadFile(fmt.Sprintf("dependencies/%s", "python-requirements.txt"))
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}

	runnerDependencies.PythonRequirements = string(file)

	return nil
}
