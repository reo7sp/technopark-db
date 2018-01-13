package dbutil

import (
	"os"
	"os/signal"
	"fmt"
	"os/exec"
)

func SubscribeKillPostgresOnInterupt() {
	if os.Getenv("KILL_POSTGRES") != "1" {
		return
	}
	fmt.Println("warning: will kill postgres on interrupt")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for range c {
			cmd := "kill -9 $(pgrep postgres)"
			exec.Command("sh","-c", cmd).Output()
		}
	}()
}
