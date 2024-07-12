// main package contains the driver code for running the application
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.nhat.io/clock"
	"lite-api/internal/app"
	"lite-api/internal/client/hotelbeds"
	"lite-api/internal/service/hotel"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	appPortEnv           = "APP_PORT"
	defaultAppPort       = ":8080"
	hotelbedsHostEnv     = "HOTELBEDS_HOST"
	defaultHotelbedsHost = "https://api.test.hotelbeds.com"
	hotelbedsApiKeyEnv   = "HOTELBEDS_API_KEY"
	hotelbedsSecretEnv   = "HOTELBEDS_SECRET"
)

// serve handles the logic of running  server in a goroutine and waiting for signal to gracefully stop the server
// on ctx.Done signal a request to shut down the server is sent, so that no new requests will be served.
func serve(ctx context.Context, appPort string, hotel *app.Hotel) {
	if appPort[0] != ':' {
		appPort = ":" + appPort
	}
	router := hotel.RegisterRoutes()

	srv := &http.Server{Addr: appPort, Handler: router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen:%s\n", err)
		}
	}()

	log.Printf("server started on port %s", appPort)

	<-ctx.Done()

	log.Printf("graceful shutdown request received")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%s", err.Error())
	}

	log.Println("application stopped accepting requests")
}

func start(appPort, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {
	realClock := clock.New()
	hotelbedsClient := hotelbeds.NewHotelBeds(hotelbedsHost, hotelbedsApiKey, hotelbedsSecret, realClock)
	hotelbedsService := hotel.NewHotelService(hotelbedsClient)
	app := app.NewHotel(hotelbedsService)

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

	serve(ctx, appPort, app)
}

// main initiates new app from argument receiver over cli args or env and calls serve to start the server
// it also spawns a goroutine to listen to os signals SIGINT or SIGTERM
// once the os signal is received the cancel func of ctx passed to serve is called
// notifying it to initiate a graceful shutdown
func main() {
	var (
		appPort         string
		hotelbedsHost   string
		hotelbedsApiKey string
		hotelbedsSecret string
	)

	// Initialize Viper
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault(appPortEnv, defaultAppPort)
	viper.SetDefault(hotelbedsHostEnv, defaultHotelbedsHost)
	// Load from environment variables
	viper.BindEnv(appPortEnv)
	viper.BindEnv(hotelbedsHostEnv)
	viper.BindEnv(hotelbedsApiKeyEnv)
	viper.BindEnv(hotelbedsSecretEnv)

	// Create start command
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the application",
		Run: func(cmd *cobra.Command, args []string) {
			// Get values from command line flags or environment variables
			appPort = viper.GetString(appPortEnv)
			hotelbedsHost = viper.GetString(hotelbedsHostEnv)
			hotelbedsApiKey = viper.GetString(hotelbedsApiKeyEnv)
			hotelbedsSecret = viper.GetString(hotelbedsSecretEnv)

			start(appPort, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret)
		},
	}

	// Bind command line flags
	startCmd.Flags().StringVarP(&appPort, "port", "p", defaultAppPort, "Application port")
	startCmd.Flags().StringVarP(&hotelbedsHost, "host", "o", defaultHotelbedsHost, "Hotelbeds API host")
	startCmd.Flags().StringVarP(&hotelbedsApiKey, "apikey", "k", "", "Hotelbeds API key")
	startCmd.Flags().StringVarP(&hotelbedsSecret, "secret", "s", "", "Hotelbeds API secret")

	// Bind flags with viper
	viper.BindPFlag(appPortEnv, startCmd.Flags().Lookup("port"))
	viper.BindPFlag(hotelbedsHostEnv, startCmd.Flags().Lookup("host"))
	viper.BindPFlag(hotelbedsApiKeyEnv, startCmd.Flags().Lookup("apikey"))
	viper.BindPFlag(hotelbedsSecretEnv, startCmd.Flags().Lookup("secret"))

	// Create root command
	var rootCmd = &cobra.Command{Use: "lite-api"}
	rootCmd.AddCommand(startCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
