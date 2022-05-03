package signalhandler

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func sendSignals(sugar *zap.SugaredLogger, sig os.Signal, subscribers []chan os.Signal, done chan string) {
	for _, sub := range subscribers {
		sub <- sig
	}
	sugar.Warnw("Sent shutdown signals to components")
	chanCount := len(subscribers)
	start := time.Now()
	delay := time.Millisecond * 100
	timeout := time.Millisecond * 2000

	for {
		select {
		case _ = <-done:
			chanCount--
			if chanCount == 0 {
				sugar.Warnw("All components shut down successfully")
				return
			}
		default:
			time.Sleep(delay)
		}
		if time.Since(start) > timeout {
			sugar.Errorw("Timouted while waiting for proper shutdown",
				"componentsStillUp", strconv.Itoa(chanCount),
				"timeout", timeout.String(),
			)
			return
		}
	}
}

// Run catches and handles signals
func Run(sugar *zap.SugaredLogger, subscribers []chan os.Signal, done chan string, server *http.Server) {
	sugar = sugar.With("module", "signalhandler")
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	var sig os.Signal
	for {
		sig := <-signals
		sugarSig := sugar.With("signal", sig.String())
		sugarSig.Infow("Got signal")
		if sig == syscall.SIGTERM {
			// Shutdown daemon on SIGTERM
			break
		}
		sugarSig.Warnw("Ignoring signal. Send SIGTERM to trigger shutdown.")
	}

	sugar.Infow("Sending shutdown signals to components ...")
	sendSignals(sugar, sig, subscribers, done)

	sugar.Infow("Shutting down the server ...")
	if err := server.Shutdown(context.Background()); err != nil {
		sugar.Errorw("Error while shuting down HTTP server",
			"error", err,
		)
	}
}
