package gasmeter

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGasMeter(t *testing.T) {
	require := require.New(t)
	gs, err := NewGasStation("https://ethgasstation.info/json/ethgasAPI.json", time.Minute)
	require.NoError(err)

	gas, dur := gs.Estimate(GasPrioritySafeLow)
	log.Printf("Safe Low: %s Gwei %s", gas.StringGwei(), dur)
	gas, dur = gs.Estimate(GasPriorityFast)
	log.Printf("Fast: %s Gwei %s", gas.StringGwei(), dur)
	gas, dur = gs.Estimate(GasPriorityFastest)
	log.Printf("Fastest: %s Gwei %s", gas.StringGwei(), dur)
}
