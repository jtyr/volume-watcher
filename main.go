package main

import (
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"net/http"
)

var (
	build string
	logger log.Logger
)

type volumeWatcher struct {
	Dir      string
	Endpoint string
}

func (vw *volumeWatcher) callEndpoint() {
	level.Info(logger).Log(
		"msg", "calling endpoint",
		"endpoint", vw.Endpoint)

	resp, err := http.Get(vw.Endpoint)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to call endpoint",
			"endpoint", vw.Endpoint,
			"err", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		level.Error(logger).Log(
			"msg", "wrong response status code",
			"endpoint", vw.Endpoint,
			"statuscode", resp.StatusCode)
	}
}

func (vw *volumeWatcher) addWatcher() {
	level.Info(logger).Log(
		"msg", "adding watcher",
		"dir", vw.Dir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to create new watcher",
			"dir", vw.Dir,
			"err", err)
		os.Exit(1)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					level.Error(logger).Log(
						"msg", "failed to get watcher events",
						"dir", vw.Dir)
					return
				}

				if event.Op&fsnotify.Create == fsnotify.Create && strings.HasSuffix(event.Name, "..data") {
					vw.callEndpoint()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					level.Error(logger).Log(
						"msg", "failed to get watcher errors",
						"dir", vw.Dir,
						"err", err)
					return
				}
			}
		}
	}()

	err = watcher.Add(vw.Dir)
	if err != nil {
		level.Error(logger).Log(
			"msg", "failed to add watcher",
			"dir", vw.Dir,
			"err", err)
		os.Exit(1)
	}

	<-done
}

func main() {
	// Setting up logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(
		logger,
		"t", log.DefaultTimestampUTC,
		"app", "volume-watcher",
		"build", build)

	dir := os.Getenv("VOLUMEWATCHER_DIR")
	endpoint := os.Getenv("VOLUMEWATCHER_ENDPOINT")

	// Fail if no dir defined
	if dir == "" {
		level.Error(logger).Log(
			"msg", "no VOLUMEWATCHER_DIR defined")
		os.Exit(1)
	}

	// Fail if no endpoint defined
	if endpoint == "" {
		level.Error(logger).Log(
			"msg", "no VOLUMEWATCHER_ENDPOINT defined")
		os.Exit(1)
	}

	// Define volume watcher params
	vw := volumeWatcher{
		Dir:      dir,
		Endpoint: endpoint,
	}

	// Add watcher
	vw.addWatcher()
}
