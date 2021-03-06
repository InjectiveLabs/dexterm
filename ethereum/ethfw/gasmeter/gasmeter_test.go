package gasmeter

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestGasMeter(t *testing.T) {
	require := require.New(t)
	gs, err := NewGasStation("https://ethgasstation.info/json/ethgasAPI.json", time.Minute)
	require.NoError(err)

	gas, dur := gs.Estimate(GasPrioritySafeLow)
	logrus.Printf("Safe Low: %s Gwei %s", gas.StringGwei(), dur)
	gas, dur = gs.Estimate(GasPriorityFast)
	logrus.Printf("Fast: %s Gwei %s", gas.StringGwei(), dur)
	gas, dur = gs.Estimate(GasPriorityFastest)
	logrus.Printf("Fastest: %s Gwei %s", gas.StringGwei(), dur)
}
