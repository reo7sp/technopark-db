package dbutil

import (
	"os/exec"
)

func KillPostgres() {
	cmd := "/etc/init.d/postgresql stop -m immediate"
	exec.Command("bash", "-c", cmd).Output()

	cmd = "sleep 10 && kill -9 $(pgrep postgres)"
	exec.Command("bash", "-c", cmd).Output()
}
