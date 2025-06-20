package nodejs

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
)

const (
	LibPath     = "/var/sandbox/sandbox-nodejs"
	LibName     = "nodejs.so"
	ProjectName = "nodejs-project"
)

//go:embed nodejs.so
var nodejsLib []byte

//go:embed dependens
var nodejsDependens embed.FS // it's a directory

func init() {
	releaseLibBinary()
}

func releaseLibBinary() {
	log.Println("initializing nodejs runner environment...")
	os.RemoveAll(LibPath)

	err := os.MkdirAll(LibPath, 0755)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to create %s", LibPath))
	}
	err = os.WriteFile(path.Join(LibPath, LibName), nodejsLib, 0755)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to write %s", path.Join(LibPath, ProjectName)))
	}

	// copy the nodejs project into /tmp/sandbox-nodejs-project
	err = os.MkdirAll(path.Join(LibPath, ProjectName), 0755)
	if err != nil {
		log.Panic(fmt.Sprintf("failed to create %s", path.Join(LibPath, ProjectName)))
	}

	// copy the nodejs project into /tmp/sandbox-nodejs-project
	var recursively_copy func(src string, dst string) error
	recursively_copy = func(src string, dst string) error {
		entries, err := nodejsDependens.ReadDir(src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			src_path := src + "/" + entry.Name()
			dst_path := dst + "/" + entry.Name()
			if entry.IsDir() {
				err = os.Mkdir(dst_path, 0755)
				if err != nil {
					return err
				}
				err = recursively_copy(src_path, dst_path)
				if err != nil {
					return err
				}
			} else {
				data, err := nodejsDependens.ReadFile(src_path)
				if err != nil {
					return err
				}
				err = os.WriteFile(dst_path, data, 0755)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}

	err = recursively_copy("dependens", path.Join(LibPath, ProjectName))
	if err != nil {
		log.Panic("failed to copy nodejs project")
	}
	log.Println("nodejs runner environment initialized")
}

func checkLibAvaliable() bool {
	if _, err := os.Stat(path.Join(LibPath, LibName)); err != nil {
		return false
	}

	return true
}
