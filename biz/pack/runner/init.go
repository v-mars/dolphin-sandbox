package runner

import (
	"dolphin-sandbox/biz/pack/static"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func init() {
	// create sandbox user
	user := static.SandboxUser
	uid := static.SandboxUserUid

	// check if user exists
	_, err := exec.Command("id", user).Output()
	if err != nil {
		// create user
		output, err := exec.Command("bash", "-c", "useradd -u "+strconv.Itoa(uid)+" "+user).Output()
		if err != nil {
			log.Panicf("failed to create user: %v, %v", err, string(output))
		}
	}

	// get gid of sandbox user and setgid
	gid, err := exec.Command("id", "-g", static.SandboxUser).Output()
	if err != nil {
		log.Panic("failed to get gid of user: %v", err)
	}

	static.SandboxGroupId, err = strconv.Atoi(strings.TrimSpace(string(gid)))
	if err != nil {
		log.Panic("failed to convert gid: %v", err)
	}
}
