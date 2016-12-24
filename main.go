package main

import (
	"syscall"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"os"
	"os/exec"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage %v <port>\n", os.Args[0])
		os.Exit(2)
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
	err := RebootWithInit()

	if err != nil {
		fmt.Println("Failed to reboot via init. Attempting direct reboot...")
	}

	err = RebootWithSyscall()

	if err != nil {
		fmt.Printf("Failed: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}

func RebootWithInit() error {
	return exec.Command("/sbin/reboot").Wait()
}

const LINUX_REBOOT_CMD_RESTART uintptr = 0x1234567

func RebootWithSyscall() error {
	return syscall.Reboot(int(LINUX_REBOOT_CMD_RESTART))
}