package static

import (
	"dolphin-sandbox/biz/pack/schema"
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

	python_pip_mirror_url := os.Getenv("PIP_MIRROR_URL")
	if python_pip_mirror_url != "" {
		SandboxGlobalConfigurations.PythonPipMirrorURL = python_pip_mirror_url
	}

	python_deps_update_interval := os.Getenv("PYTHON_DEPS_UPDATE_INTERVAL")
	if python_deps_update_interval != "" {
		SandboxGlobalConfigurations.PythonDepsUpdateInterval = python_deps_update_interval
	}

	// if not set "PythonDepsUpdateInterval", update python dependencies every 30 minutes to keep the sandbox up-to-date
	if SandboxGlobalConfigurations.PythonDepsUpdateInterval == "" {
		SandboxGlobalConfigurations.PythonDepsUpdateInterval = "30m"
	}

	nodejs_path := os.Getenv("NODEJS_PATH")
	if nodejs_path != "" {
		SandboxGlobalConfigurations.NodejsPath = nodejs_path
	}

	if SandboxGlobalConfigurations.NodejsPath == "" {
		SandboxGlobalConfigurations.NodejsPath = "/usr/local/bin/node"
	}

	enable_network := os.Getenv("ENABLE_NETWORK")
	if enable_network != "" {
		SandboxGlobalConfigurations.EnableNetwork, _ = strconv.ParseBool(enable_network)
	}

	enable_preload := os.Getenv("ENABLE_PRELOAD")
	if enable_preload != "" {
		SandboxGlobalConfigurations.EnablePreload, _ = strconv.ParseBool(enable_preload)
	}

	allowed_syscalls := os.Getenv("ALLOWED_SYSCALLS")
	if allowed_syscalls != "" {
		strs := strings.Split(allowed_syscalls, ",")
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
		socks5_proxy := os.Getenv("SOCKS5_PROXY")
		if socks5_proxy != "" {
			SandboxGlobalConfigurations.Proxy.Socks5 = socks5_proxy
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

// avoid global modification, use value copy instead
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
	file, err := os.ReadFile("dependencies/python-requirements.txt")
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}

	runnerDependencies.PythonRequirements = string(file)

	return nil
}
