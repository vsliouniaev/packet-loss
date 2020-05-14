package cmd

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/spf13/cobra"
	"github.com/vsliouniaev/packet-loss/core"
	"github.com/vsliouniaev/packet-loss/logging"
	"github.com/vsliouniaev/packet-loss/receiver"
	"github.com/vsliouniaev/packet-loss/sender"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	rootCmd = &cobra.Command{
		Use:     "packet-loss",
		Short:   "Run to figure out if you're dropping packets",
		PreRun:  configureLogging,
		Run:     rootCommand,
		Version: core.Version,
	}

	logCfg       = &logging.Config{}
	receiveAddr  string
	sendAddr     string
	sendInterval time.Duration

	logger log.Logger
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags()
	rootCmd.PersistentFlags().StringVar(&receiveAddr, "receive.host", ":8000", "Specify ip:port to receive on")
	rootCmd.PersistentFlags().StringVar(&sendAddr, "send.host", "localhost:8000", "Specify target to send to")
	rootCmd.PersistentFlags().DurationVar(&sendInterval, "send.interval", time.Millisecond*10, "Specify target to send to")
	logging.AddFlags(rootCmd, logCfg)
}

func configureLogging(_ *cobra.Command, _ []string) {
	logger = logging.New(logCfg)
}

func rootCommand(cmd *cobra.Command, _ []string) {
	var (
		hup  = make(chan os.Signal, 1)
		term = make(chan os.Signal, 1)
	)
	signal.Notify(hup, syscall.SIGHUP)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	rec, err := receiver.New(log.With(logger, "component", "receiver"), receiveAddr)
	if err != nil {
		level.Error(logger).Log("msg", "cannot start receiver", "err", err)
		os.Exit(1)
	}
	sen, err := sender.New(log.With(logger, "component", "sender"), sendAddr, sendInterval)
	if err != nil {
		level.Error(logger).Log("msg", "cannot start sender", "err", err)
		os.Exit(1)
	}

	go rec.Receive()
	go sen.Send()

	for {
		select {
		case <-term:
			level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			sen.Stop()
			rec.Stop()
			return
		}
	}
}
