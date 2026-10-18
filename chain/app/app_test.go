package app_test

import (
	"context"
	"encoding/json"
	"testing"

	"cosmossdk.io/log"
	"github.com/cliffordx/opencasedigests/chain/app"
	abcitypes "github.com/cometbft/cometbft/abci/types"
)

func TestCaseDigestorApp_InitChain(t *testing.T) {
	logger := log.NewNopLogger()
	casedApp := app.NewCaseDigestorApp(logger)

	req := &abcitypes.RequestInitChain{
		ChainId: "casedigestor-testnet-1",
	}

	resp, err := casedApp.InitChain(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected InitChain error: %v", err)
	}
	if len(resp.AppHash) == 0 {
		t.Fatal("expected non-empty app hash")
	}
}

func TestCaseDigestorApp_CheckTx_GasCeiling(t *testing.T) {
	logger := log.NewNopLogger()
	casedApp := app.NewCaseDigestorApp(logger)

	// Valid message type with excess gas (> 500,000)
	tx := app.TxEnvelope{
		TypeURL: "/casedigestor.jurisprudence.v1.MsgAnchorDecision",
		Sender:  "casedigestor1...",
		Gas:     600_000,
		Payload: json.RawMessage(`{}`),
	}
	txBytes, _ := json.Marshal(tx)

	resp, err := casedApp.CheckTx(context.Background(), &abcitypes.RequestCheckTx{Tx: txBytes})
	if err != nil {
		t.Fatalf("unexpected CheckTx error: %v", err)
	}
	if resp.Code != 2 {
		t.Fatalf("expected code 2 (gas ceiling exceeded), got %d", resp.Code)
	}
}

func TestCaseDigestorApp_CheckTx_AllowedMessage(t *testing.T) {
	logger := log.NewNopLogger()
	casedApp := app.NewCaseDigestorApp(logger)

	tx := app.TxEnvelope{
		TypeURL: "/casedigestor.jurisprudence.v1.MsgAnchorDecision",
		Sender:  "casedigestor1...",
		Gas:     100_000,
		Payload: json.RawMessage(`{}`),
	}
	txBytes, _ := json.Marshal(tx)

	resp, err := casedApp.CheckTx(context.Background(), &abcitypes.RequestCheckTx{Tx: txBytes})
	if err != nil {
		t.Fatalf("unexpected CheckTx error: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("expected code 0, got %d (log: %s)", resp.Code, resp.Log)
	}
}

func TestCaseDigestorApp_FinalizeBlock_And_Query(t *testing.T) {
	logger := log.NewNopLogger()
	casedApp := app.NewCaseDigestorApp(logger)

	// Init chain
	_, err := casedApp.InitChain(context.Background(), &abcitypes.RequestInitChain{ChainId: "casedigestor-testnet-1"})
	if err != nil {
		t.Fatalf("InitChain failed: %v", err)
	}

	// Create AnchorDecision transaction
	decPayload := app.DecisionRecord{
		GRNumber:         "G.R. No. 212345",
		Title:            "People v. Dela Cruz",
		PromulgationDate: "2024-01-15T00:00:00Z",
		DivisionOrEnBanc: "EN_BANC",
		Ponente:          "Leonen, J.",
		IPFSCID:          "bafybeicmockcid123456789",
		SHA256Checksum:   "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2",
	}
	payloadBytes, _ := json.Marshal(decPayload)

	tx := app.TxEnvelope{
		TypeURL: "/casedigestor.jurisprudence.v1.MsgAnchorDecision",
		Sender:  "casedigestor1authority...",
		Gas:     100_000,
		Payload: payloadBytes,
	}
	txBytes, _ := json.Marshal(tx)

	// Finalize block at height 1
	finResp, err := casedApp.FinalizeBlock(context.Background(), &abcitypes.RequestFinalizeBlock{
		Height: 1,
		Txs:    [][]byte{txBytes},
	})
	if err != nil {
		t.Fatalf("FinalizeBlock failed: %v", err)
	}
	if len(finResp.TxResults) != 1 || finResp.TxResults[0].Code != 0 {
		t.Fatalf("tx failed: %+v", finResp.TxResults[0])
	}
	if len(finResp.TxResults[0].Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(finResp.TxResults[0].Events))
	}

	// Query decision by G.R. number
	qResp, err := casedApp.Query(context.Background(), &abcitypes.RequestQuery{
		Path: "/casedigestor/jurisprudence/v1/decisions",
		Data: []byte("G.R. No. 212345"),
	})
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	if qResp.Code != 0 {
		t.Fatalf("expected query code 0, got %d (log: %s)", qResp.Code, qResp.Log)
	}

	var queried app.DecisionRecord
	if err := json.Unmarshal(qResp.Value, &queried); err != nil {
		t.Fatalf("failed to unmarshal query result: %v", err)
	}
	if queried.GRNumber != "G.R. No. 212345" || queried.Ponente != "Leonen, J." {
		t.Fatalf("mismatched queried decision: %+v", queried)
	}
}
