package cli

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCreateStartCmdHandler(t *testing.T) {
	t.Run("Successful command creation", func(t *testing.T) {
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {}
		cmd, err := CreateStartCmdHandler(mockStart)

		assert.NoError(t, err)
		assert.NotNil(t, cmd)
		assert.Equal(t, "start", cmd.Use)
		assert.Equal(t, "Start lite-api application", cmd.Short)
	})

	t.Run("Flag bindings", func(t *testing.T) {
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {}
		cmd, err := CreateStartCmdHandler(mockStart)

		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		assert.NotNil(t, cmd.Flags().Lookup("port"))
		assert.NotNil(t, cmd.Flags().Lookup("mode"))
		assert.NotNil(t, cmd.Flags().Lookup("host"))
		assert.NotNil(t, cmd.Flags().Lookup("apikey"))
		assert.NotNil(t, cmd.Flags().Lookup("secret"))
	})

	t.Run("Viper bindings", func(t *testing.T) {
		viper.Reset()
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {}
		cmd, err := CreateStartCmdHandler(mockStart)

		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		assert.NotNil(t, viper.GetString(appPortEnv))
		assert.NotNil(t, viper.GetString(appModeEnv))
		assert.NotNil(t, viper.GetString(hotelbedsHostEnv))
		assert.NotNil(t, viper.GetString(hotelbedsApiKeyEnv))
		assert.NotNil(t, viper.GetString(hotelbedsSecretEnv))
	})

	t.Run("Command execution", func(t *testing.T) {
		var executedStart bool
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {
			executedStart = true
			assert.Equal(t, defaultAppPort, appPort)
			assert.Equal(t, defaultAppMode, appMode)
			assert.Equal(t, defaultHotelbedsHost, hotelbedsHost)
			assert.Empty(t, hotelbedsApiKey)
			assert.Empty(t, hotelbedsSecret)
		}

		cmd, err := CreateStartCmdHandler(mockStart)
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		assert.NoError(t, cmd.Execute())

		assert.True(t, executedStart)
	})

	t.Run("Command execution with flags", func(t *testing.T) {
		var executedStart bool
		mockStart := func(appPort, appMode, hotelbedsHost, hotelbedsApiKey, hotelbedsSecret string) {
			executedStart = true
			assert.Equal(t, "8080", appPort)
			assert.Equal(t, "test", appMode)
			assert.Equal(t, "testhost", hotelbedsHost)
			assert.Equal(t, "testkey", hotelbedsApiKey)
			assert.Equal(t, "testsecret", hotelbedsSecret)
		}

		cmd, err := CreateStartCmdHandler(mockStart)
		assert.NoError(t, err)
		assert.NotNil(t, cmd)

		require.NoError(t, cmd.Flags().Set("port", "8080"))
		require.NoError(t, cmd.Flags().Set("mode", "test"))
		require.NoError(t, cmd.Flags().Set("host", "testhost"))
		require.NoError(t, cmd.Flags().Set("apikey", "testkey"))
		require.NoError(t, cmd.Flags().Set("secret", "testsecret"))

		require.NoError(t, cmd.Execute())

		assert.True(t, executedStart)
	})
}
