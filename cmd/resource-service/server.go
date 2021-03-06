package main

import (
	"net/http"
	"os"
	"time"

	"os/signal"

	"context"

	"fmt"
	"text/tabwriter"

	"git.containerum.net/ch/resource-service/pkg/router"
	m "git.containerum.net/ch/resource-service/pkg/router/middleware"
	"git.containerum.net/ch/resource-service/pkg/util/validation"
	"github.com/containerum/kube-client/pkg/model"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func initServer(c *cli.Context) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent|tabwriter.Debug)
	for _, f := range c.GlobalFlagNames() {
		fmt.Fprintf(w, "Flag: %s\t Value: %s\n", f, c.String(f))
	}
	w.Flush()

	setupLogs(c)

	translate := setupTranslator()
	validate := validation.StandardResourceValidator(translate)

	tv := &m.TranslateValidate{UniversalTranslator: translate, Validate: validate}

	mongo, err := setupMongo(c)
	exitOnError(err)
	defer mongo.Close()

	err = mongo.Init()
	exitOnError(err)

	kube, err := setupKube(c)
	exitOnError(err)

	permissions := setupPermissions(c)

	status := model.ServiceStatus{
		Name:     c.App.Name,
		Version:  c.App.Version,
		StatusOK: true,
	}

	app := router.CreateRouter(mongo, permissions, kube, &status, tv, c.Bool("cors"), c.String("ingress_suffix"), c.Uint("min_port"), c.Uint("max_port"))

	srv := &http.Server{
		Addr:    ":" + c.String("port"),
		Handler: app,
	}

	// serve connections
	go exitOnError(srv.ListenAndServe())

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt) // subscribe on interrupt event
	<-quit                            // wait for event
	logrus.Infoln("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}

func exitOnError(err error) {
	if err != nil {
		logrus.WithError(err).Fatalf("can`t setup resource-service")
		os.Exit(1)
	}
}
