package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"syscall"

	"github.com/gen2brain/beeep"
	"github.com/gorilla/handlers"
)

const (
	appName = "mwan3-notify-fcgi"
)

var (
	quiet         = flag.Bool("q", false, "do not produce any output")
	sockName      = flag.String("l", "/run/"+appName+"/fcgi.sock", "unix socket path for listening")
	unlinkSock    = flag.Bool("u", false, "unlink unix socket path")
	allowedSecret = flag.String("s", "", "only notifications with matching request secret are displayed")
	appIcon       = flag.String("i", "", "app icon for displayed notification")
)

func fcgiHandle(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	secret := r.Form.Get("secret")
	if secret != *allowedSecret {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	title := fmt.Sprintf("mwan3: %s", r.Form.Get("hostname"))
	message := fmt.Sprintf("'%s' (%s) '%s'", r.Form.Get("device"), r.Form.Get("interface"), r.Form.Get("action"))

	if err := beeep.Alert(title, message, *appIcon); err != nil {
		log.Printf("[alert] unable to send %s <%s>: %s", title, message, err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func main() {
	flag.Parse()
	if *quiet {
		log.SetOutput(ioutil.Discard)
	}

	addr := &net.UnixAddr{Name: *sockName, Net: "unix"}
	log.Printf("[startup] %s: will listen `%s'", appName, addr.Name)
	if *unlinkSock {
		_ = syscall.Unlink(addr.Name)
	}

	mask := syscall.Umask(0)
	ln, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Fatalf("[startup] %s: unable to listen: %s", appName, err)
	}
	_ = syscall.Umask(mask)

	hf := http.HandlerFunc(fcgiHandle)
	if err := fcgi.Serve(ln, handlers.LoggingHandler(os.Stderr, hf)); err != nil {
		log.Fatalf("[serve] %s: fcgi: %s", appName, err)
	}
}
