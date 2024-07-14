package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StartFunc func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string)

func CreateStartCmdHandler(start StartFunc) (*cobra.Command, error) {
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
	if err := viper.BindPFlag(appPortEnv, startCmd.Flags().Lookup("port")); err != nil {
		return nil, err
	}

	if err := viper.BindPFlag(appModeEnv, startCmd.Flags().Lookup("mode")); err != nil {
		return nil, err
	}

	if err := viper.BindPFlag(hotelbedsHostEnv, startCmd.Flags().Lookup("host")); err != nil {
		return nil, err
	}

	if err := viper.BindPFlag(hotelbedsApiKeyEnv, startCmd.Flags().Lookup("apikey")); err != nil {
		return nil, err
	}

	if err := viper.BindPFlag(hotelbedsSecretEnv, startCmd.Flags().Lookup("secret")); err != nil {
		return nil, err
	}

	return startCmd, nil
}
