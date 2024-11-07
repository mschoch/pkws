package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mschoch/pkws"
)

func main() {

	var (
		httpAddr = flag.String("http.addr", ":8179", "HTTP listen address")
	)
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	log.Info("transport HTTP", "addr", *httpAddr)
	var ln net.Listener
	var err error
	ln, err = net.Listen("tcp", *httpAddr)
	if err != nil {
		log.Info(fmt.Sprintf("error listening"), "add", *httpAddr, "err", err.Error())
		os.Exit(1)
	}

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	env := &fireproofContext{}

	// create the party server
	ps := pkws.NewPartyServer(env, "/parties", log)

	// register the custom server for the fireproof party
	ps.RegisterServer("fireproof", NewFireproof)

	srv := &http.Server{Addr: *httpAddr, Handler: ps}
	go func() {
		errs <- srv.Serve(ln)
	}()

	log.Info("exit", <-errs)
}
