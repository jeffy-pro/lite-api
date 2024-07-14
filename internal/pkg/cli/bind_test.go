package cli

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestBindEnv(t *testing.T) {
	viper.Reset()

	t.Run("Default Values", func(t *testing.T) {
		BindEnv()

		require.Equal(t, DefaultAppPort, viper.GetString(AppPortEnv))
		require.Equal(t, DefaultHotelbedsHost, viper.GetString(HotelbedsHostEnv))
		require.Equal(t, DefaultAppMode, viper.GetString(AppModeEnv))
	})

	t.Run("Custom Environment Variables", func(t *testing.T) {
		os.Setenv(AppPortEnv, "9000")
		os.Setenv(HotelbedsHostEnv, "custom.hotelbeds.com")
		os.Setenv(HotelbedsApiKeyEnv, "testApiKey")
		os.Setenv(HotelbedsSecretEnv, "testSecret")

		BindEnv()

		require.Equal(t, "9000", viper.GetString(AppPortEnv))
		require.Equal(t, "custom.hotelbeds.com", viper.GetString(HotelbedsHostEnv))
		require.Equal(t, "testApiKey", viper.GetString(HotelbedsApiKeyEnv))
		require.Equal(t, "testSecret", viper.GetString(HotelbedsSecretEnv))

		// Clean up environment variables
		os.Unsetenv(AppPortEnv)
		os.Unsetenv(HotelbedsHostEnv)
		os.Unsetenv(HotelbedsApiKeyEnv)
		os.Unsetenv(HotelbedsSecretEnv)
	})

	t.Run("Missing Required Environment Variables", func(t *testing.T) {
		BindEnv()

		require.Empty(t, viper.GetString(HotelbedsApiKeyEnv))
		require.Empty(t, viper.GetString(HotelbedsSecretEnv))
	})
}
