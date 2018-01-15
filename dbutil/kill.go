package dbutil

import (
	"log"
	"os/exec"
)

func KillPostgres() {
	cmd := "kill -9 $(pgrep postgres)"
	out, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Println("error: kill postgres: cannot execute kill command:", err)
		return
	}
	log.Println("info: kill postgres: out:", out)
}
