package native

import (
	"testing"

	"github.com/hyperledger/burrow/acm"
	"github.com/hyperledger/burrow/acm/acmstate"
	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/execution/engine"
	"github.com/hyperledger/burrow/execution/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestState_CreateAccount(t *testing.T) {
	st := acmstate.NewMemoryState()
	address := AddressFromName("frogs")
	err := CreateAccount(st, address)
	require.NoError(t, err)
	err = CreateAccount(st, address)
	require.Error(t, err)
	require.Equal(t, errors.Codes.DuplicateAddress, errors.GetCode(err))

	st = acmstate.NewMemoryState()
	err = CreateAccount(st, address)
	require.NoError(t, err)
	err = InitEVMCode(st, address, []byte{1, 2, 3})
	require.NoError(t, err)
}

func TestState_Sync(t *testing.T) {
	backend := acmstate.NewCache(acmstate.NewMemoryState())
	st := engine.NewCallFrame(backend)
	address := AddressFromName("frogs")

	err := CreateAccount(st, address)
	require.NoError(t, err)
	amt := uint64(1232)
	addToBalance(t, st, address, amt)
	err = st.Sync()
	require.NoError(t, err)

	acc, err := backend.GetAccount(address)
	require.NoError(t, err)
	assert.Equal(t, acc.Balance, amt)
}

func addToBalance(t testing.TB, st acmstate.ReaderWriter, address crypto.Address, amt uint64) {
	err := UpdateAccount(st, address, func(account *acm.Account) error {
		return account.AddToBalance(amt)
	})
	require.NoError(t, err)
}
