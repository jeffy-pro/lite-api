package cli

import "github.com/spf13/viper"

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

func BindEnv() {
	// Initialize Viper
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault(appPortEnv, defaultAppPort)
	viper.SetDefault(hotelbedsHostEnv, defaultHotelbedsHost)
	viper.SetDefault(appModeEnv, defaultAppMode)

	for _, env := range []string{appPortEnv, hotelbedsHostEnv, hotelbedsApiKeyEnv, hotelbedsSecretEnv} {
		_ = viper.BindEnv(env)
	}
}
