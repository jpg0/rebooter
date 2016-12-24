package main

import (
	"syscall"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const DEFAULT_DELAY = time.Second * 3

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage %v <port>\n", os.Args[0])
		os.Exit(2)
	}

	if os.Geteuid() != 0 {
		fmt.Print("Must be run as root")
		os.Exit(3)
	}

	r := mux.NewRouter()
	r.Methods("DELETE").Path("/").HandlerFunc(Reboot)
	http.Handle("/", r)
	err := http.ListenAndServe(":" + os.Args[1], nil)

	if err != nil {
		panic(err)
	}
}
func Reboot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Remote request received, attempting reboot...")
	go ScheduleReboot(DEFAULT_DELAY)
	w.WriteHeader(http.StatusAccepted)
}

func ScheduleReboot(d time.Duration) {

	time.Sleep(d)

	err := RebootWithInit()

	if err != nil {
		fmt.Println("Failed to reboot via init. Attempting direct reboot...")
	}

	err = RebootWithSyscall()

	if err != nil {
		fmt.Println("Failed to reboot directly")
	}
}

func RebootWithInit() error {
	return exec.Command("/sbin/reboot").Wait()
}

const LINUX_REBOOT_CMD_RESTART uintptr = 0x1234567

func RebootWithSyscall() error {
	return syscall.Reboot(int(LINUX_REBOOT_CMD_RESTART))
}