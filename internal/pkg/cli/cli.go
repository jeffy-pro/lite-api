package cli

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StartFunc func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string, logger *slog.Logger)

func CreateStartCmdHandler(start StartFunc, logger *slog.Logger) (*cobra.Command, error) {
	var (
		appPort         string
		appMode         string
		hotelbedsHost   string
		hotelbedsApiKey string
		hotelbedsSecret string
	)

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start lite-api application",
		Run: func(cmd *cobra.Command, args []string) {
			// Get values from command line flags or environment variables
			appPort = viper.GetString(AppPortEnv)
			appMode = viper.GetString(AppModeEnv)
			hotelbedsHost = viper.GetString(HotelbedsHostEnv)
			hotelbedsApiKey = viper.GetString(HotelbedsApiKeyEnv)
			hotelbedsSecret = viper.GetString(HotelbedsSecretEnv)

			start(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret, logger)
		},
	}

	// Bind command line flags
	startCmd.Flags().StringVarP(&appPort, "port", "p", DefaultAppPort, "Application port")
	startCmd.Flags().StringVarP(&appMode, "mode", "m", DefaultAppMode, "Application mode")
	startCmd.Flags().StringVarP(&hotelbedsHost, "host", "o", DefaultHotelbedsHost, "Hotelbeds API host")
	startCmd.Flags().StringVarP(&hotelbedsApiKey, "apikey", "k", "", "Hotelbeds API key")
	startCmd.Flags().StringVarP(&hotelbedsSecret, "secret", "s", "", "Hotelbeds API secret")

	// Bind flags with viper
	if err := viper.BindPFlag(AppPortEnv, startCmd.Flags().Lookup("port")); err != nil {
		return nil, err
	}

	if err := viper.BindPFlag(AppModeEnv, startCmd.Flags().Lookup("mode")); err != nil {
		return nil, err
	}

	if err := viper.BindPFlag(HotelbedsHostEnv, startCmd.Flags().Lookup("host")); err != nil {
		return nil, err
	}

	if err := viper.BindPFlag(HotelbedsApiKeyEnv, startCmd.Flags().Lookup("apikey")); err != nil {
		return nil, err
	}

	if err := viper.BindPFlag(HotelbedsSecretEnv, startCmd.Flags().Lookup("secret")); err != nil {
		return nil, err
	}

	return startCmd, nil
}
