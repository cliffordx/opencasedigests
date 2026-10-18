/**
 * Gasless Transaction Helper for Case Digestor Protocol
 *
 * Wraps curator and counsel transaction messages inside a
 * cosmos.authz.v1beta1.MsgExec envelope with the Case Digestor
 * Master Authority as fee granter, enabling zero-fee execution
 * for law students, paralegals, and attorneys.
 *
 * This module uses @cosmjs/stargate and @cosmjs/proto-signing to
 * construct, sign, and broadcast gasless transactions through the
 * x/authz + x/feegrant sponsorship mechanism.
 *
 * @module gasless-tx
 */

import { Registry, GeneratedType } from "@cosmjs/proto-signing";
import {
  SigningStargateClient,
  StargateClient,
  StdFee,
  DeliverTxResponse,
} from "@cosmjs/stargate";
import { EncodeObject } from "@cosmjs/proto-signing";

// ── Configuration ──────────────────────────────────────────────────────

/**
 * The Bech32 address of the Case Digestor Master Authority that sponsors
 * all zero-fee transactions via x/feegrant.
 *
 * This address must have an active feegrant allowance for the grantee
 * (the user submitting the transaction).
 */
const MASTER_AUTHORITY_ADDRESS =
  process.env.NEXT_PUBLIC_MASTER_AUTHORITY_ADDRESS ??
  "casedigestor1masterauthorityplaceholder";

/**
 * Default gas limit for zero-fee transactions.
 * Must not exceed the ZeroGasAgencyDecorator's MaxZeroFeeGas (500,000).
 */
const DEFAULT_GAS_LIMIT = "500000";

/**
 * The chain's native denomination. For Case Digestor Protocol,
 * the token is used only for validator staking, never by end users.
 */
const DENOM = process.env.NEXT_PUBLIC_CHAIN_DENOM ?? "ucasedigestor";

// ── Type URLs ──────────────────────────────────────────────────────────

const TYPE_URLS = {
  // cosmos.authz.v1beta1
  msgExec: "/cosmos.authz.v1beta1.MsgExec",

  // casedigestor.attestations.v1
  msgSubmitDigestDraft:
    "/casedigestor.attestations.v1.MsgSubmitDigestDraft",
  msgSubmitCuratorRevision:
    "/casedigestor.attestations.v1.MsgSubmitCuratorRevision",
  msgAttestDigest: "/casedigestor.attestations.v1.MsgAttestDigest",
  msgLinkSubsequentRuling:
    "/casedigestor.attestations.v1.MsgLinkSubsequentRuling",

  // casedigestor.jurisprudence.v1
  msgAnchorDecision:
    "/casedigestor.jurisprudence.v1.MsgAnchorDecision",
} as const;

// ── Message Interfaces ─────────────────────────────────────────────────

/** Parameters for submitting a new case digest draft. */
export interface SubmitDigestDraftParams {
  /** Bech32 address of the submitter or AI coordinator. */
  creator: string;
  /** Supreme Court docket number (e.g., "G.R. No. 212345"). */
  grNumber: string;
  /** IPFS CIDv1 of the structured digest JSON. */
  ipfsCid: string;
  /** Hex-encoded SHA-256 checksum of the raw file bytes. */
  sha256Checksum: string;
}

/** Parameters for submitting a curator revision of an existing digest. */
export interface SubmitCuratorRevisionParams {
  /** Bech32 address of the community curator. */
  curator: string;
  /** ID of the digest being revised. */
  digestId: string;
  /** IPFS CIDv1 of the revised digest JSON. */
  revisedIpfsCid: string;
  /** Hex-encoded SHA-256 checksum of the revised file bytes. */
  revisedSha256Checksum: string;
  /** Optional IPFS CID of revision notes. */
  revisionNotesIpfsCid?: string;
}

/** Parameters for a verified attorney's peer attestation. */
export interface AttestDigestParams {
  /** Bech32 address of the verified attorney. */
  attorney: string;
  /** ID of the digest being attested. */
  digestId: string;
  /** Rubric score for doctrinal accuracy (1-5). */
  doctrineAccuracyScore: number;
  /** Rubric score for citation integrity (1-5). */
  citationIntegrityScore: number;
  /** Optional IPFS CID of detailed peer review notes. */
  reviewNotesIpfsCid?: string;
}

// ── Zero-Fee StdFee Builder ────────────────────────────────────────────

/**
 * Build a zero-fee StdFee with the Master Authority as granter.
 *
 * The `granter` field tells the chain to deduct fees from the Master
 * Authority's feegrant allowance instead of the signer's account.
 * Combined with the ZeroGasAgencyDecorator, this means the actual
 * fee amount is 0 tokens.
 */
function buildZeroFee(gasLimit: string = DEFAULT_GAS_LIMIT): StdFee {
  return {
    amount: [{ denom: DENOM, amount: "0" }],
    gas: gasLimit,
    granter: MASTER_AUTHORITY_ADDRESS,
  };
}

// ── MsgExec Wrapper ────────────────────────────────────────────────────

/**
 * Wrap an inner message inside a cosmos.authz.v1beta1.MsgExec envelope.
 *
 * The MsgExec is executed by the grantee on behalf of the granter,
 * using a pre-existing authz authorization. The fee granter in the
 * StdFee covers the (zero) gas cost via x/feegrant.
 *
 * @param grantee - Bech32 address of the user executing the transaction.
 * @param innerMsg - The EncodeObject representing the inner message.
 * @returns An EncodeObject for the MsgExec wrapper.
 */
function wrapInMsgExec(
  grantee: string,
  innerMsg: EncodeObject,
): EncodeObject {
  return {
    typeUrl: TYPE_URLS.msgExec,
    value: {
      grantee,
      msgs: [innerMsg],
    },
  };
}

// ── Gasless Transaction Functions ──────────────────────────────────────

/**
 * Submit a new case digest draft as a gasless transaction.
 *
 * The MsgSubmitDigestDraft is wrapped in a MsgExec with the Master
 * Authority as fee granter — the user pays zero fees.
 *
 * @param client - A connected SigningStargateClient instance.
 * @param params - The digest draft parameters.
 * @returns The broadcast transaction result.
 */
export async function submitDigestDraftGasless(
  client: SigningStargateClient,
  params: SubmitDigestDraftParams,
): Promise<DeliverTxResponse> {
  const innerMsg: EncodeObject = {
    typeUrl: TYPE_URLS.msgSubmitDigestDraft,
    value: {
      creator: params.creator,
      grNumber: params.grNumber,
      ipfsCid: params.ipfsCid,
      sha256Checksum: params.sha256Checksum,
    },
  };

  const execMsg = wrapInMsgExec(params.creator, innerMsg);
  const fee = buildZeroFee();

  return client.signAndBroadcast(params.creator, [execMsg], fee);
}

/**
 * Submit a curator revision as a gasless transaction.
 *
 * Wraps MsgSubmitCuratorRevision inside MsgExec with zero fees.
 *
 * @param client - A connected SigningStargateClient instance.
 * @param params - The curator revision parameters.
 * @returns The broadcast transaction result.
 */
export async function submitCuratorRevisionGasless(
  client: SigningStargateClient,
  params: SubmitCuratorRevisionParams,
): Promise<DeliverTxResponse> {
  const innerMsg: EncodeObject = {
    typeUrl: TYPE_URLS.msgSubmitCuratorRevision,
    value: {
      curator: params.curator,
      digestId: params.digestId,
      revisedIpfsCid: params.revisedIpfsCid,
      revisedSha256Checksum: params.revisedSha256Checksum,
      revisionNotesIpfsCid: params.revisionNotesIpfsCid ?? "",
    },
  };

  const execMsg = wrapInMsgExec(params.curator, innerMsg);
  const fee = buildZeroFee();

  return client.signAndBroadcast(params.curator, [execMsg], fee);
}

/**
 * Submit a peer attestation as a gasless transaction.
 *
 * Wraps MsgAttestDigest inside MsgExec with zero fees. Only verified
 * attorneys (ROLE_COUNSEL_REVIEWER) may call this.
 *
 * @param client - A connected SigningStargateClient instance.
 * @param params - The attestation parameters.
 * @returns The broadcast transaction result.
 */
export async function attestDigestGasless(
  client: SigningStargateClient,
  params: AttestDigestParams,
): Promise<DeliverTxResponse> {
  const innerMsg: EncodeObject = {
    typeUrl: TYPE_URLS.msgAttestDigest,
    value: {
      attorney: params.attorney,
      digestId: params.digestId,
      doctrineAccuracyScore: params.doctrineAccuracyScore,
      citationIntegrityScore: params.citationIntegrityScore,
      reviewNotesIpfsCid: params.reviewNotesIpfsCid ?? "",
    },
  };

  const execMsg = wrapInMsgExec(params.attorney, innerMsg);
  const fee = buildZeroFee();

  return client.signAndBroadcast(params.attorney, [execMsg], fee);
}
