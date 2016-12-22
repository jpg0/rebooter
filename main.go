package main

import (
	"syscall"
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"os"
	"log"
	"strconv"
)

func main() {

	if len(os.Args) != 1 {
		log.Fatalf("Usage %v <port>", os.Args[0])
	}

	port, err := strconv.Atoi(os.Args[1])

	if err != nil {
		log.Fatalf("Failed to parse port number %v [%v]", port, err)
	}

	r := mux.NewRouter()
	r.Methods("DELETE").Path("/").HandlerFunc(Reboot)
	http.Handle("/", r)
	err = http.ListenAndServe(":8080", nil)

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
