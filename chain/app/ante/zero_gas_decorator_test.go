package ante_test

import (
	"testing"

	"github.com/cliffordx/opencasedigests/chain/app/ante"
	sdk "github.com/cosmos/cosmos-sdk/types"
	protov2 "google.golang.org/protobuf/proto"
)

// mockMsg is a simple implementation of sdk.Msg for unit testing.
type mockMsg struct {
	typeURL string
}

func (m mockMsg) Reset()                       {}
func (m mockMsg) String() string               { return m.typeURL }
func (m mockMsg) ProtoMessage()                {}
func (m mockMsg) ValidateBasic() error         { return nil }
func (m mockMsg) GetSigners() []sdk.AccAddress { return nil }

// mockFeeTx implements sdk.FeeTx for testing.
type mockFeeTx struct {
	fee  sdk.Coins
	gas  uint64
	msgs []sdk.Msg
}

func (tx mockFeeTx) GetGas() uint64                            { return tx.gas }
func (tx mockFeeTx) GetFee() sdk.Coins                         { return tx.fee }
func (tx mockFeeTx) FeePayer() []byte                          { return nil }
func (tx mockFeeTx) FeeGranter() []byte                        { return nil }
func (tx mockFeeTx) GetMsgs() []sdk.Msg                        { return tx.msgs }
func (tx mockFeeTx) GetMsgsV2() ([]protov2.Message, error)     { return nil, nil }
func (tx mockFeeTx) ValidateBasic() error                      { return nil }

func TestMaxZeroFeeGasConstant(t *testing.T) {
	if ante.MaxZeroFeeGas != 500_000 {
		t.Fatalf("expected MaxZeroFeeGas to be 500000, got %d", ante.MaxZeroFeeGas)
	}
}

func TestZeroGasAgencyDecorator_AllowsNonZeroFeeTx(t *testing.T) {
	decorator := ante.NewZeroGasAgencyDecorator()
	ctx := sdk.Context{}

	// Transaction with a non-zero fee should pass directly to next
	tx := mockFeeTx{
		fee: sdk.NewCoins(sdk.NewInt64Coin("uatom", 1000)),
		gas: 200_000,
		msgs: []sdk.Msg{
			mockMsg{typeURL: "/cosmos.bank.v1beta1.MsgSend"},
		},
	}

	nextCalled := false
	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		nextCalled = true
		return ctx, nil
	}

	_, err := decorator.AnteHandle(ctx, tx, false, next)
	if err != nil {
		t.Fatalf("unexpected error for non-zero fee tx: %v", err)
	}
	if !nextCalled {
		t.Fatal("expected next handler to be called")
	}
}

func TestZeroGasAgencyDecorator_RejectsDisallowedZeroFeeMsg(t *testing.T) {
	decorator := ante.NewZeroGasAgencyDecorator()
	ctx := sdk.Context{}

	// Zero fee transaction with a non-whitelisted message type
	tx := mockFeeTx{
		fee: sdk.NewCoins(),
		gas: 200_000,
		msgs: []sdk.Msg{
			mockMsg{typeURL: "/cosmos.bank.v1beta1.MsgSend"},
		},
	}

	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		return ctx, nil
	}

	_, err := decorator.AnteHandle(ctx, tx, false, next)
	if err == nil {
		t.Fatal("expected error for disallowed zero fee message, got nil")
	}
}

func TestZeroGasAgencyDecorator_RejectsEmptyMsgs(t *testing.T) {
	decorator := ante.NewZeroGasAgencyDecorator()
	ctx := sdk.Context{}

	tx := mockFeeTx{
		fee:  sdk.NewCoins(),
		gas:  200_000,
		msgs: []sdk.Msg{},
	}

	next := func(ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, error) {
		return ctx, nil
	}

	_, err := decorator.AnteHandle(ctx, tx, false, next)
	if err == nil {
		t.Fatal("expected error for empty messages in zero fee tx, got nil")
	}
}
