package predeploys

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
)

func TestGethAddresses(t *testing.T) {
	// We test if the addresses in geth match those in bindings, to avoid an import-cycle:
	// we import geth in the monorepo, and do not want to import bindings into geth.
	require.Equal(t, L1BlockAddr, types.L1BlockAddr)
}

// TestL1BlockSlots ensures that the storage layout of the L1Block
// contract matches the hardcoded values in `kroma-geth`.
func TestL1BlockSlots(t *testing.T) {
	layout, err := bindings.GetStorageLayout("L1Block")
	require.NoError(t, err)

	var l1BaseFeeSlot, overHeadSlot, scalarSlot common.Hash
	for _, entry := range layout.Storage {
		switch entry.Label {
		case "l1FeeOverhead":
			overHeadSlot = common.BigToHash(big.NewInt(int64(entry.Slot)))
		case "l1FeeScalar":
			scalarSlot = common.BigToHash(big.NewInt(int64(entry.Slot)))
		case "basefee":
			l1BaseFeeSlot = common.BigToHash(big.NewInt(int64(entry.Slot)))
		}
	}

	require.Equal(t, types.OverheadSlot, overHeadSlot)
	require.Equal(t, types.ScalarSlot, scalarSlot)
	require.Equal(t, types.L1BaseFeeSlot, l1BaseFeeSlot)
}
