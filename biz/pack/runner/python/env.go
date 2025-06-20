package python

import (
	"dolphin-sandbox/biz/pack/runner"
	"dolphin-sandbox/biz/pack/static"
	_ "embed"
	"log"
	"os"
	"os/exec"
	"path"
)

//go:embed env.sh
var env_script string

func PreparePythonDependenciesEnv() error {
	config := static.GetSandboxGlobalConfigurations()

	runner := runner.TempDirRunner{}
	err := runner.WithTempDir("/", []string{}, func(root_path string) error {
		err := os.WriteFile(path.Join(root_path, "env.sh"), []byte(env_script), 0755)
		if err != nil {
			return err
		}

		for _, lib_path := range config.PythonLibPaths {
			// check if the lib path is available
			if _, err := os.Stat(lib_path); err != nil {
				log.Println("python lib path %s is not available", lib_path)
				continue
			}
			exec_cmd := exec.Command(
				"bash",
				path.Join(root_path, "env.sh"),
				lib_path,
				LIB_PATH,
			)
			exec_cmd.Stderr = os.Stderr

			if err := exec_cmd.Run(); err != nil {
				return err
			}
		}

		os.RemoveAll(root_path)
		return nil
	})

	return err
}
