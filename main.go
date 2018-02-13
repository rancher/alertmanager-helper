package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/alertmanager-helper/watcher"
	"github.com/urfave/cli"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "alertmanager-helper"
	app.Version = VERSION
	app.Usage = "alertmanager helper is used to watch the alertmanager config file change, and make api call to alertmanager to reload its config"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "watched-file-list",
			Usage: "config file path list for what file need to be watched",
		},
	}
	app.Action = run

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	wg := sync.WaitGroup{}
	jobs := make(chan int, 5)

	filePathList := c.StringSlice("watched-file-list")
	for _, v := range filePathList {
		wg.Add(1)
		go func(file string, jobs <-chan int) {
			defer wg.Done()
			logrus.Info("watched file:", file)
			watcher.Watcherfile(file, jobs)
		}(v, jobs)
	}

	waitForSignal()
	close(jobs)
	wg.Wait()

	fmt.Println("exiting")
	return nil
}

func waitForSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	fmt.Println("receive exiting signal")

}
