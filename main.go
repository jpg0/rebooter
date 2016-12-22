package main

import (
	"syscall"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage %v <port>", os.Args[0])
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

const LINUX_REBOOT_CMD_RESTART uintptr = 0x1234567

func Reboot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Remote request received, attempting reboot...")
	err := syscall.Reboot(int(LINUX_REBOOT_CMD_RESTART))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusAccepted)
	}
}
