package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateStartCmdHandler(t *testing.T) {
	t.Run("Successful command creation", func(t *testing.T) {
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {}
		cmd, err := CreateStartCmdHandler(mockStart)

		require.NoError(t, err)
		require.NotNil(t, cmd)
		require.Equal(t, "start", cmd.Use)
		require.Equal(t, "Start lite-api application", cmd.Short)
	})

	t.Run("Flag bindings", func(t *testing.T) {
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {}
		cmd, err := CreateStartCmdHandler(mockStart)

		require.NoError(t, err)
		require.NotNil(t, cmd)

		require.NotNil(t, cmd.Flags().Lookup("port"))
		require.NotNil(t, cmd.Flags().Lookup("mode"))
		require.NotNil(t, cmd.Flags().Lookup("host"))
		require.NotNil(t, cmd.Flags().Lookup("apikey"))
		require.NotNil(t, cmd.Flags().Lookup("secret"))
	})

	t.Run("Viper bindings", func(t *testing.T) {
		viper.Reset()
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {}
		cmd, err := CreateStartCmdHandler(mockStart)

		require.NoError(t, err)
		require.NotNil(t, cmd)

		require.NotNil(t, viper.GetString(appPortEnv))
		require.NotNil(t, viper.GetString(appModeEnv))
		require.NotNil(t, viper.GetString(hotelbedsHostEnv))
		require.NotNil(t, viper.GetString(hotelbedsApiKeyEnv))
		require.NotNil(t, viper.GetString(hotelbedsSecretEnv))
	})

	t.Run("Command execution", func(t *testing.T) {
		var executedStart bool
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {
			executedStart = true
			require.Equal(t, defaultAppPort, appPort)
			require.Equal(t, defaultAppMode, appMode)
			require.Equal(t, defaultHotelbedsHost, hotelbedsHost)
			require.Empty(t, hotelbedsApiKey)
			require.Empty(t, hotelbedsSecret)
		}

		cmd, err := CreateStartCmdHandler(mockStart)
		require.NoError(t, err)
		require.NotNil(t, cmd)

		require.NoError(t, cmd.Execute())

		require.True(t, executedStart)
	})

	t.Run("Command execution with flags", func(t *testing.T) {
		var executedStart bool
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {
			executedStart = true
			require.Equal(t, "8080", appPort)
			require.Equal(t, "test", appMode)
			require.Equal(t, "testhost", hotelbedsHost)
			require.Equal(t, "testkey", hotelbedsApiKey)
			require.Equal(t, "testsecret", hotelbedsSecret)
		}

		cmd, err := CreateStartCmdHandler(mockStart)
		require.NoError(t, err)
		require.NotNil(t, cmd)

		require.NoError(t, cmd.Flags().Set("port", "8080"))
		require.NoError(t, cmd.Flags().Set("mode", "test"))
		require.NoError(t, cmd.Flags().Set("host", "testhost"))
		require.NoError(t, cmd.Flags().Set("apikey", "testkey"))
		require.NoError(t, cmd.Flags().Set("secret", "testsecret"))

		require.NoError(t, cmd.Execute())

		assert.True(t, executedStart)
	})
}
