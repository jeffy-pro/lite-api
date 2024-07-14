package cli

import "github.com/spf13/viper"

const (
	AppPortEnv           = "APP_PORT"
	DefaultAppPort       = ":8080"
	HotelbedsHostEnv     = "HOTELBEDS_HOST"
	DefaultHotelbedsHost = "https://api.test.hotelbeds.com"
	HotelbedsApiKeyEnv   = "HOTELBEDS_API_KEY"
	HotelbedsSecretEnv   = "HOTELBEDS_SECRET"
	AppModeEnv           = "MODE"
	DefaultAppMode       = "dev"
	LogLevel             = "LOG_LEVEL"
)

func BindEnv() {
	// Initialize Viper
	viper.AutomaticEnv()

	// Set default values
	viper.SetDefault(AppPortEnv, DefaultAppPort)
	viper.SetDefault(HotelbedsHostEnv, DefaultHotelbedsHost)
	viper.SetDefault(AppModeEnv, DefaultAppMode)

	for _, env := range []string{AppPortEnv, HotelbedsHostEnv, HotelbedsApiKeyEnv, HotelbedsSecretEnv} {
		_ = viper.BindEnv(env)
	}
}
