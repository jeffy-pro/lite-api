// main package contains the driver code for running the application
package main

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.nhat.io/clock"
	"lite-api/internal/app"
	"lite-api/internal/client/hotelbeds"
	"lite-api/internal/pkg/server"
	"lite-api/internal/service/hotel"
	"log"
	"os"
	"os/signal"
	"syscall"
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
