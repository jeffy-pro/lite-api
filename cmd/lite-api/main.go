// main package contains the driver code for running the application
package main

import (
	"context"
	"lite-api/internal/app"
	"lite-api/internal/client/hotelbeds"
	"lite-api/internal/pkg/cli"
	"lite-api/internal/pkg/server"
	"lite-api/internal/service/hotel"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.nhat.io/clock"
)

func start(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {
	realClock := clock.New()
	hotelbedsClient := hotelbeds.NewHotelBeds(hotelbedsHost, hotelbedsApiKey, hotelbedsSecret, realClock)
	hotelsService := hotel.NewHotelService(hotelbedsClient)
	hotelApp := app.NewHotel(hotelsService, appMode)

	defer func() {
		if err := recover(); err != nil {
			log.Printf("recovering from panic %s", err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		log.Printf("system call received")
		cancel()
	}()

	handler := hotelApp.RegisterRoutes()
	server.ServeHTTP(ctx, appPort, handler)
}

// main initiates new app from argument receiver over cli args or env and calls serve to start the server
// it also spawns a goroutine to listen to os signals SIGINT or SIGTERM
// once the os signal is received the cancel func of ctx passed to serve is called
// notifying it to initiate a graceful shutdown
func main() {
	cli.BindEnv()

	startCmdHandler, err := cli.CreateStartCmdHandler(start)
	if err != nil {
		log.Fatal("error booting lite-api application", err)
	}

	var rootCmd = &cobra.Command{Use: "lite-api"}
	rootCmd.AddCommand(startCmdHandler)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("error starting lite-api application", err)
	}
}
