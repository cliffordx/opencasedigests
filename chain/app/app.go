// Package app provides the Case Digestor ABCI application and consensus state machine.
package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	"cosmossdk.io/log"
	abcitypes "github.com/cometbft/cometbft/abci/types"
)

const (
	// AppName is the identifier for the Case Digestor appchain.
	AppName = "casedigestord"
	// AppVersion is the semantic version of the protocol state machine.
	AppVersion = "0.1.0"
	// DefaultChainID is the identifier for the institutional consortium testnet.
	DefaultChainID = "casedigestor-testnet-1"
	// MaxZeroFeeGas is the hard gas ceiling enforced by the ante handler.
	MaxZeroFeeGas int64 = 500_000
)

// CaseDigestorApp implements the CometBFT abci.Application interface.
// It manages the state transitions for Supreme Court decision anchoring
// and peer-attested case digest lifecycles with zero gas execution.
type CaseDigestorApp struct {
	abcitypes.BaseApplication

	logger log.Logger
	mu     sync.RWMutex

	// State stores
	chainID      string
	height       int64
	appHash      []byte
	decisions    map[string]*DecisionRecord
	digests      map[string]*DigestRecord
	attestations map[string][]*AttestationRecord
	citationEdges map[string][]*CitationEdgeRecord
}

// DecisionRecord stores an anchored Supreme Court promulgation in memory.
type DecisionRecord struct {
	GRNumber       string `json:"gr_number"`
	Title          string `json:"title"`
	PromulgationDate string `json:"promulgation_date"`
	DivisionOrEnBanc string `json:"division_or_enbanc"`
	Ponente        string `json:"ponente"`
	IPFSCID        string `json:"ipfs_cid"`
	SHA256Checksum string `json:"sha256_checksum"`
	AnchoredBy     string `json:"anchored_by"`
	BlockHeight    int64  `json:"block_height"`
}

// DigestRecord stores a case digest record in memory.
type DigestRecord struct {
	DigestID       string `json:"digest_id"`
	GRNumber       string `json:"gr_number"`
	CurrentStage   string `json:"current_stage"`
	IPFSCID        string `json:"ipfs_cid"`
	SHA256Checksum string `json:"sha256_checksum"`
	SubmittedBy    string `json:"submitted_by"`
	AttestationCount int   `json:"attestation_count"`
	LastRevisedBy  string `json:"last_revised_by,omitempty"`
}

// AttestationRecord stores an attorney peer attestation.
type AttestationRecord struct {
	Attorney              string `json:"attorney"`
	DigestID              string `json:"digest_id"`
	DoctrineAccuracyScore int    `json:"doctrine_accuracy_score"`
	CitationIntegrityScore int   `json:"citation_integrity_score"`
	ReviewNotesIPFSCID    string `json:"review_notes_ipfs_cid,omitempty"`
}

// CitationEdgeRecord stores a doctrinal citation link between decisions.
type CitationEdgeRecord struct {
	Author              string `json:"author"`
	TargetDigestID      string `json:"target_digest_id"`
	SupersedingGRNumber string `json:"superseding_gr_number"`
	Relation            string `json:"relation"`
	IPFSProofCID        string `json:"ipfs_proof_cid"`
}

// TxEnvelope represents an incoming state-changing transaction.
type TxEnvelope struct {
	TypeURL string          `json:"type_url"`
	Sender  string          `json:"sender"`
	Gas     int64           `json:"gas"`
	Payload json.RawMessage `json:"payload"`
}

// NewCaseDigestorApp initializes a new Case Digestor ABCI application instance.
func NewCaseDigestorApp(logger log.Logger) *CaseDigestorApp {
	return &CaseDigestorApp{
		logger:        logger,
		chainID:       DefaultChainID,
		height:        0,
		appHash:       make([]byte, 32),
		decisions:     make(map[string]*DecisionRecord),
		digests:       make(map[string]*DigestRecord),
		attestations:  make(map[string][]*AttestationRecord),
		citationEdges: make(map[string][]*CitationEdgeRecord),
	}
}

// Info returns basic node information to CometBFT during handshake.
func (app *CaseDigestorApp) Info(_ context.Context, _ *abcitypes.RequestInfo) (*abcitypes.ResponseInfo, error) {
	app.mu.RLock()
	defer app.mu.RUnlock()

	return &abcitypes.ResponseInfo{
		Data:             AppName,
		Version:          AppVersion,
		AppVersion:       1,
		LastBlockHeight:  app.height,
		LastBlockAppHash: app.appHash,
	}, nil
}

// InitChain initializes genesis state from the genesis.json document.
func (app *CaseDigestorApp) InitChain(_ context.Context, req *abcitypes.RequestInitChain) (*abcitypes.ResponseInitChain, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.chainID = req.ChainId
	app.height = 0
	app.appHash = sha256Sum([]byte(app.chainID))

	app.logger.Info("initialized Case Digestor Protocol state machine", "chain_id", app.chainID)
	return &abcitypes.ResponseInitChain{
		AppHash: app.appHash,
	}, nil
}

// CheckTx validates incoming transactions for the mempool, enforcing the
// ZeroGasAgencyDecorator rules and hard gas ceiling.
func (app *CaseDigestorApp) CheckTx(_ context.Context, req *abcitypes.RequestCheckTx) (*abcitypes.ResponseCheckTx, error) {
	var tx TxEnvelope
	if err := json.Unmarshal(req.Tx, &tx); err != nil {
		return &abcitypes.ResponseCheckTx{
			Code: 1,
			Log:  fmt.Sprintf("invalid tx format: %v", err),
		}, nil
	}

	// Enforce gas ceiling
	if tx.Gas > MaxZeroFeeGas {
		return &abcitypes.ResponseCheckTx{
			Code: 2,
			Log:  fmt.Sprintf("gas limit %d exceeds max zero fee gas ceiling %d", tx.Gas, MaxZeroFeeGas),
		}, nil
	}

	// Validate allowed message types
	switch tx.TypeURL {
	case "/casedigestor.jurisprudence.v1.MsgAnchorDecision",
		"/casedigestor.attestations.v1.MsgSubmitDigestDraft",
		"/casedigestor.attestations.v1.MsgSubmitCuratorRevision",
		"/casedigestor.attestations.v1.MsgAttestDigest",
		"/casedigestor.attestations.v1.MsgLinkSubsequentRuling":
		// Whitelisted message type passes ante validation
		return &abcitypes.ResponseCheckTx{
			Code:      0,
			GasWanted: tx.Gas,
			GasUsed:   10_000,
		}, nil
	default:
		return &abcitypes.ResponseCheckTx{
			Code: 3,
			Log:  fmt.Sprintf("message type %s not eligible for zero-fee execution", tx.TypeURL),
		}, nil
	}
}

// FinalizeBlock processes transactions in a decided block, mutating state
// and emitting events for the indexer worker to ingest.
func (app *CaseDigestorApp) FinalizeBlock(_ context.Context, req *abcitypes.RequestFinalizeBlock) (*abcitypes.ResponseFinalizeBlock, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	app.height = req.Height
	txResults := make([]*abcitypes.ExecTxResult, len(req.Txs))

	for i, rawTx := range req.Txs {
		var tx TxEnvelope
		if err := json.Unmarshal(rawTx, &tx); err != nil {
			txResults[i] = &abcitypes.ExecTxResult{
				Code: 1,
				Log:  fmt.Sprintf("decode error: %v", err),
			}
			continue
		}

		result := app.executeTx(&tx, req.Height)
		txResults[i] = result
	}

	// Compute state root hash
	stateBlob := fmt.Sprintf("%s:%d:%d:%d", app.chainID, app.height, len(app.decisions), len(app.digests))
	app.appHash = sha256Sum([]byte(stateBlob))

	return &abcitypes.ResponseFinalizeBlock{
		TxResults: txResults,
		AppHash:   app.appHash,
	}, nil
}

// executeTx executes a validated transaction and produces emitted events.
func (app *CaseDigestorApp) executeTx(tx *TxEnvelope, height int64) *abcitypes.ExecTxResult {
	switch tx.TypeURL {
	case "/casedigestor.jurisprudence.v1.MsgAnchorDecision":
		var dec DecisionRecord
		if err := json.Unmarshal(tx.Payload, &dec); err != nil {
			return &abcitypes.ExecTxResult{Code: 1, Log: err.Error()}
		}
		dec.BlockHeight = height
		dec.AnchoredBy = tx.Sender
		app.decisions[dec.GRNumber] = &dec

		event := abcitypes.Event{
			Type: "casedigestor.jurisprudence.v1.EventDecisionAnchored",
			Attributes: []abcitypes.EventAttribute{
				{Key: "gr_number", Value: dec.GRNumber, Index: true},
				{Key: "ipfs_cid", Value: dec.IPFSCID, Index: true},
				{Key: "sha256_checksum", Value: dec.SHA256Checksum, Index: true},
				{Key: "anchored_by", Value: dec.AnchoredBy, Index: true},
				{Key: "block_height", Value: fmt.Sprintf("%d", height), Index: false},
			},
		}

		return &abcitypes.ExecTxResult{
			Code:      0,
			GasUsed:   50_000,
			GasWanted: tx.Gas,
			Events:    []abcitypes.Event{event},
		}

	case "/casedigestor.attestations.v1.MsgSubmitDigestDraft":
		var dig DigestRecord
		if err := json.Unmarshal(tx.Payload, &dig); err != nil {
			return &abcitypes.ExecTxResult{Code: 1, Log: err.Error()}
		}
		dig.SubmittedBy = tx.Sender
		dig.CurrentStage = "DIGEST_STAGE_COMMUNITY_REVIEW"
		app.digests[dig.DigestID] = &dig

		event := abcitypes.Event{
			Type: "casedigestor.attestations.v1.EventDigestSubmitted",
			Attributes: []abcitypes.EventAttribute{
				{Key: "digest_id", Value: dig.DigestID, Index: true},
				{Key: "gr_number", Value: dig.GRNumber, Index: true},
				{Key: "ipfs_cid", Value: dig.IPFSCID, Index: true},
				{Key: "sha256_checksum", Value: dig.SHA256Checksum, Index: true},
				{Key: "submitted_by", Value: dig.SubmittedBy, Index: true},
				{Key: "stage", Value: dig.CurrentStage, Index: true},
			},
		}

		return &abcitypes.ExecTxResult{
			Code:      0,
			GasUsed:   40_000,
			GasWanted: tx.Gas,
			Events:    []abcitypes.Event{event},
		}

	case "/casedigestor.attestations.v1.MsgAttestDigest":
		var att AttestationRecord
		if err := json.Unmarshal(tx.Payload, &att); err != nil {
			return &abcitypes.ExecTxResult{Code: 1, Log: err.Error()}
		}
		att.Attorney = tx.Sender
		app.attestations[att.DigestID] = append(app.attestations[att.DigestID], &att)

		// Advance stage if qualifications met
		updatedStage := "DIGEST_STAGE_COMMUNITY_REVIEW"
		totalAttestations := len(app.attestations[att.DigestID])
		if totalAttestations >= 2 && att.DoctrineAccuracyScore >= 4 && att.CitationIntegrityScore >= 4 {
			updatedStage = "DIGEST_STAGE_BAR_ATTESTED"
			if dig, ok := app.digests[att.DigestID]; ok {
				dig.CurrentStage = updatedStage
				dig.AttestationCount = totalAttestations
			}
		}

		event := abcitypes.Event{
			Type: "casedigestor.attestations.v1.EventDigestAttested",
			Attributes: []abcitypes.EventAttribute{
				{Key: "digest_id", Value: att.DigestID, Index: true},
				{Key: "attorney", Value: att.Attorney, Index: true},
				{Key: "doctrine_accuracy_score", Value: fmt.Sprintf("%d", att.DoctrineAccuracyScore), Index: false},
				{Key: "citation_integrity_score", Value: fmt.Sprintf("%d", att.CitationIntegrityScore), Index: false},
				{Key: "updated_stage", Value: updatedStage, Index: true},
				{Key: "total_attestations", Value: fmt.Sprintf("%d", totalAttestations), Index: false},
			},
		}

		return &abcitypes.ExecTxResult{
			Code:      0,
			GasUsed:   60_000,
			GasWanted: tx.Gas,
			Events:    []abcitypes.Event{event},
		}

	case "/casedigestor.attestations.v1.MsgLinkSubsequentRuling":
		var edge CitationEdgeRecord
		if err := json.Unmarshal(tx.Payload, &edge); err != nil {
			return &abcitypes.ExecTxResult{Code: 1, Log: err.Error()}
		}
		edge.Author = tx.Sender
		app.citationEdges[edge.TargetDigestID] = append(app.citationEdges[edge.TargetDigestID], &edge)

		event := abcitypes.Event{
			Type: "casedigestor.attestations.v1.EventCitationLinked",
			Attributes: []abcitypes.EventAttribute{
				{Key: "target_digest_id", Value: edge.TargetDigestID, Index: true},
				{Key: "superseding_gr_number", Value: edge.SupersedingGRNumber, Index: true},
				{Key: "relation", Value: edge.Relation, Index: true},
				{Key: "author", Value: edge.Author, Index: true},
			},
		}

		return &abcitypes.ExecTxResult{
			Code:      0,
			GasUsed:   35_000,
			GasWanted: tx.Gas,
			Events:    []abcitypes.Event{event},
		}

	default:
		return &abcitypes.ExecTxResult{
			Code: 3,
			Log:  fmt.Sprintf("unsupported message type %s", tx.TypeURL),
		}
	}
}

// Commit finalizes the block state and returns the current app hash.
func (app *CaseDigestorApp) Commit(_ context.Context, _ *abcitypes.RequestCommit) (*abcitypes.ResponseCommit, error) {
	app.mu.Lock()
	defer app.mu.Unlock()

	return &abcitypes.ResponseCommit{
		RetainHeight: 0,
	}, nil
}

// Query handles read-only queries for anchored decisions and digests.
func (app *CaseDigestorApp) Query(_ context.Context, req *abcitypes.RequestQuery) (*abcitypes.ResponseQuery, error) {
	app.mu.RLock()
	defer app.mu.RUnlock()

	switch req.Path {
	case "/casedigestor/jurisprudence/v1/decisions":
		grNumber := string(req.Data)
		dec, ok := app.decisions[grNumber]
		if !ok {
			return &abcitypes.ResponseQuery{Code: 1, Log: "decision not found"}, nil
		}
		data, _ := json.Marshal(dec)
		return &abcitypes.ResponseQuery{Code: 0, Value: data, Height: app.height}, nil

	case "/casedigestor/attestations/v1/digests":
		digestID := string(req.Data)
		dig, ok := app.digests[digestID]
		if !ok {
			return &abcitypes.ResponseQuery{Code: 1, Log: "digest not found"}, nil
		}
		data, _ := json.Marshal(dig)
		return &abcitypes.ResponseQuery{Code: 0, Value: data, Height: app.height}, nil

	default:
		return &abcitypes.ResponseQuery{Code: 2, Log: "unknown query path"}, nil
	}
}

// sha256Sum computes a 32-byte SHA-256 hash.
func sha256Sum(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}

// HexString returns a lowercase hex string for bytes.
func HexString(b []byte) string {
	return hex.EncodeToString(b)
}
