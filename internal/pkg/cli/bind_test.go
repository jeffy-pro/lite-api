package cli

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestBindEnv(t *testing.T) {
	viper.Reset()

	t.Run("Default Values", func(t *testing.T) {
		BindEnv()

		require.Equal(t, defaultAppPort, viper.GetString(appPortEnv))
		require.Equal(t, defaultHotelbedsHost, viper.GetString(hotelbedsHostEnv))
		require.Equal(t, defaultAppMode, viper.GetString(appModeEnv))
	})

	t.Run("Custom Environment Variables", func(t *testing.T) {
		os.Setenv(appPortEnv, "9000")
		os.Setenv(hotelbedsHostEnv, "custom.hotelbeds.com")
		os.Setenv(hotelbedsApiKeyEnv, "testApiKey")
		os.Setenv(hotelbedsSecretEnv, "testSecret")

		BindEnv()

		require.Equal(t, "9000", viper.GetString(appPortEnv))
		require.Equal(t, "custom.hotelbeds.com", viper.GetString(hotelbedsHostEnv))
		require.Equal(t, "testApiKey", viper.GetString(hotelbedsApiKeyEnv))
		require.Equal(t, "testSecret", viper.GetString(hotelbedsSecretEnv))

		// Clean up environment variables
		os.Unsetenv(appPortEnv)
		os.Unsetenv(hotelbedsHostEnv)
		os.Unsetenv(hotelbedsApiKeyEnv)
		os.Unsetenv(hotelbedsSecretEnv)
	})

	t.Run("Missing Required Environment Variables", func(t *testing.T) {
		BindEnv()

		require.Empty(t, viper.GetString(hotelbedsApiKeyEnv))
		require.Empty(t, viper.GetString(hotelbedsSecretEnv))
	})
}
