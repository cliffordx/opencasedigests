// Package ante provides custom ante decorators for the Case Digestor appchain.
//
// The ZeroGasAgencyDecorator allows authorized curator and counsel transactions
// to execute with zero fees, up to a hard gas ceiling of 500,000 gas units.
// This enables a "tokenless by default" user experience where law students,
// paralegals, and attorneys never need to hold cryptocurrency.
//
// Authorized transactions may be submitted directly or wrapped inside a
// cosmos.authz.v1beta1.MsgExec envelope for delegated execution via x/authz.
package ante

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/gogoproto/proto"
)

const (
	// MaxZeroFeeGas is the hard ceiling for gas consumption on zero-fee
	// transactions. Any transaction exceeding this limit will be rejected
	// even if it contains only authorized message types.
	MaxZeroFeeGas uint64 = 500_000
)

// zeroFeeAllowedMessages is the set of message type URLs that are permitted
// to execute with zero fees. These correspond to the curator and counsel
// transaction types in the attestations and jurisprudence modules.
var zeroFeeAllowedMessages = map[string]struct{}{
	"/casedigestor.attestations.v1.MsgSubmitDigestDraft":    {},
	"/casedigestor.attestations.v1.MsgSubmitCuratorRevision": {},
	"/casedigestor.attestations.v1.MsgAttestDigest":          {},
	"/casedigestor.attestations.v1.MsgLinkSubsequentRuling":  {},
	"/casedigestor.jurisprudence.v1.MsgAnchorDecision":       {},
}

// ZeroGasAgencyDecorator is an AnteDecorator that allows zero-fee execution
// for authorized message types. It enforces a hard gas ceiling to prevent
// abuse and supports unwrapping cosmos.authz.v1beta1.MsgExec envelopes to
// validate the inner messages.
//
// This decorator should be placed BEFORE the standard DeductFeeDecorator
// in the ante handler chain. When it detects a qualifying zero-fee transaction,
// it sets a finite gas limit and skips fee deduction.
type ZeroGasAgencyDecorator struct{}

// NewZeroGasAgencyDecorator creates a new ZeroGasAgencyDecorator.
func NewZeroGasAgencyDecorator() ZeroGasAgencyDecorator {
	return ZeroGasAgencyDecorator{}
}

// AnteHandle implements the sdk.AnteDecorator interface.
//
// Logic:
//  1. If the transaction carries non-zero fees, skip (let the standard
//     DeductFeeDecorator handle it normally).
//  2. Extract all messages, recursively unwrapping any MsgExec envelopes.
//  3. If every unwrapped message is in the allowed set, cap the gas meter
//     at MaxZeroFeeGas and continue the chain.
//  4. Otherwise, reject the transaction with an insufficient fees error.
func (d ZeroGasAgencyDecorator) AnteHandle(
	ctx sdk.Context,
	tx sdk.Tx,
	simulate bool,
	next sdk.AnteHandler,
) (sdk.Context, error) {
	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errors.Wrap(sdkerrors.ErrTxDecode, "transaction does not implement sdk.FeeTx")
	}

	// If fees are provided, this decorator is a no-op — let the standard
	// fee decorator handle validation and deduction.
	feeCoins := feeTx.GetFee()
	if feeCoins != nil && !feeCoins.IsZero() {
		return next(ctx, tx, simulate)
	}

	// Extract and validate all messages (unwrapping MsgExec recursively).
	msgs := tx.GetMsgs()
	innerMsgs, err := unwrapAllMessages(msgs)
	if err != nil {
		return ctx, err
	}

	if len(innerMsgs) == 0 {
		return ctx, errors.Wrap(sdkerrors.ErrInvalidRequest, "transaction contains no messages")
	}

	// Every inner message must be in the zero-fee allowed set.
	for _, msg := range innerMsgs {
		typeURL := proto.MessageName(msg)
		if typeURL == "" {
			typeURL = sdk.MsgTypeURL(msg)
		} else {
			typeURL = "/" + typeURL
		}
		if _, allowed := zeroFeeAllowedMessages[typeURL]; !allowed {
			return ctx, errors.Wrapf(
				sdkerrors.ErrInsufficientFee,
				"message type %s is not eligible for zero-fee execution",
				typeURL,
			)
		}
	}

	// Cap gas at MaxZeroFeeGas to prevent abuse of gasless transactions.
	gasLimit := feeTx.GetGas()
	if gasLimit == 0 || gasLimit > MaxZeroFeeGas {
		gasLimit = MaxZeroFeeGas
	}

	ctx = ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))

	return next(ctx, tx, simulate)
}

// unwrapAllMessages recursively extracts the inner messages from MsgExec
// envelopes. Direct (non-MsgExec) messages are returned as-is.
// This prevents circumvention of the allowlist by wrapping disallowed
// messages inside an authz MsgExec.
func unwrapAllMessages(msgs []sdk.Msg) ([]sdk.Msg, error) {
	var result []sdk.Msg

	for _, msg := range msgs {
		execMsg, ok := msg.(*authztypes.MsgExec)
		if !ok {
			// Direct message — not wrapped in MsgExec.
			result = append(result, msg)
			continue
		}

		// Unwrap the inner messages from the MsgExec envelope.
		innerMsgs, err := execMsg.GetMessages()
		if err != nil {
			return nil, errors.Wrapf(
				sdkerrors.ErrTxDecode,
				"failed to unpack MsgExec inner messages: %v",
				err,
			)
		}

		// Recursively unwrap in case of nested MsgExec (defensive).
		unwrapped, err := unwrapAllMessages(innerMsgs)
		if err != nil {
			return nil, err
		}

		result = append(result, unwrapped...)
	}

	if len(result) > 20 {
		return nil, fmt.Errorf("zero-fee transaction contains %d messages, exceeding the safety limit of 20", len(result))
	}

	return result, nil
}
