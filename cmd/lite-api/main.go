// main package contains the driver code for running the application
package main

import (
	"context"
	"errors"
	"fmt"
	"lite-api/internal/app"
	"lite-api/internal/client/hotelbeds"
	"lite-api/internal/service/hotel"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.nhat.io/clock"
)

const (
	appPortEnv           = "APP_PORT"
	defaultAppPort       = ":8080"
	hotelbedsHostEnv     = "HOTELBEDS_HOST"
	defaultHotelbedsHost = "https://api.test.hotelbeds.com"
	hotelbedsApiKeyEnv   = "HOTELBEDS_API_KEY"
	hotelbedsSecretEnv   = "HOTELBEDS_SECRET"
	appModeEnv           = "MODE"
	defaultAppMode       = "dev"
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

func start(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {
	realClock := clock.New()
	hotelbedsClient := hotelbeds.NewHotelBeds(hotelbedsHost, hotelbedsApiKey, hotelbedsSecret, realClock)
	hotelsService := hotel.NewHotelService(hotelbedsClient)
	hotelApp := app.NewHotel(hotelsService)

	if strings.ToLower(appMode) == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

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

	serve(ctx, appPort, hotelApp)
}

// main initiates new app from argument receiver over cli args or env and calls serve to start the server
// it also spawns a goroutine to listen to os signals SIGINT or SIGTERM
// once the os signal is received the cancel func of ctx passed to serve is called
// notifying it to initiate a graceful shutdown
func main() {
	var (
		appPort         string
		appMode         string
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
	_ = viper.BindEnv(appPortEnv)
	_ = viper.BindEnv(hotelbedsHostEnv)
	_ = viper.BindEnv(hotelbedsApiKeyEnv)
	_ = viper.BindEnv(hotelbedsSecretEnv)

	// Create start command
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the application",
		Run: func(cmd *cobra.Command, args []string) {
			// Get values from command line flags or environment variables
			appPort = viper.GetString(appPortEnv)
			appMode = viper.GetString(appModeEnv)
			hotelbedsHost = viper.GetString(hotelbedsHostEnv)
			hotelbedsApiKey = viper.GetString(hotelbedsApiKeyEnv)
			hotelbedsSecret = viper.GetString(hotelbedsSecretEnv)

			start(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret)
		},
	}

	// Bind command line flags
	startCmd.Flags().StringVarP(&appPort, "port", "p", defaultAppPort, "Application port")
	startCmd.Flags().StringVarP(&appMode, "mode", "m", defaultAppMode, "Application mode")
	startCmd.Flags().StringVarP(&hotelbedsHost, "host", "o", defaultHotelbedsHost, "Hotelbeds API host")
	startCmd.Flags().StringVarP(&hotelbedsApiKey, "apikey", "k", "", "Hotelbeds API key")
	startCmd.Flags().StringVarP(&hotelbedsSecret, "secret", "s", "", "Hotelbeds API secret")

	// Bind flags with viper
	_ = viper.BindPFlag(appPortEnv, startCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag(appModeEnv, startCmd.Flags().Lookup("mode"))
	_ = viper.BindPFlag(hotelbedsHostEnv, startCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag(hotelbedsApiKeyEnv, startCmd.Flags().Lookup("apikey"))
	_ = viper.BindPFlag(hotelbedsSecretEnv, startCmd.Flags().Lookup("secret"))

	// Create root command
	var rootCmd = &cobra.Command{Use: "lite-api"}
	rootCmd.AddCommand(startCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
