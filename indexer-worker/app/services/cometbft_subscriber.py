"""CometBFT WebSocket subscriber for the Case Digestor Protocol.

Connects to a CometBFT node's WebSocket endpoint and subscribes to
on-chain events (primarily EventDigestAttested). When an event is received,
the subscriber:

1. Extracts the IPFS CID and SHA-256 checksum from the event attributes.
2. Fetches the digest payload from IPFS.
3. Verifies the SHA-256 hash against the on-chain anchor.
4. Invokes a callback for downstream processing (e.g., Supabase pgvector indexing).

This implements the Phase 1 "Hybrid Anchoring" reactive indexing pipeline:
Chain Event → IPFS Fetch → SHA-256 Verify → PostgreSQL/pgvector Index.
"""

import asyncio
import json
import logging
from collections.abc import Awaitable, Callable
from dataclasses import dataclass

import websockets

from app.core.crypto import verify_integrity
from app.services.ipfs_client import IPFSClient

logger = logging.getLogger(__name__)

# CometBFT event query for digest attestation events
SUBSCRIBE_QUERY = "tm.event='Tx' AND EventDigestAttested.digest_id EXISTS"


@dataclass(frozen=True)
class VerifiedDigestEvent:
    """A digest attestation event that has been fetched from IPFS and
    verified against its on-chain SHA-256 anchor.

    Attributes:
        digest_id: The unique identifier of the attested digest.
        attorney: Bech32 address of the attesting attorney.
        ipfs_cid: The IPFS CIDv1 of the digest content.
        sha256_checksum: The on-chain SHA-256 anchor.
        content: The raw bytes fetched from IPFS (verified).
        doctrine_accuracy_score: Rubric score from the attestation.
        citation_integrity_score: Rubric score from the attestation.
        updated_stage: The digest stage after attestation.
        total_attestations: Total number of attestations after this one.
    """

    digest_id: str
    attorney: str
    ipfs_cid: str
    sha256_checksum: str
    content: bytes
    doctrine_accuracy_score: int
    citation_integrity_score: int
    updated_stage: str
    total_attestations: int


# Type alias for the callback invoked after successful verification.
OnVerifiedCallback = Callable[[VerifiedDigestEvent], Awaitable[None]]


class CometBFTSubscriber:
    """Async WebSocket listener for CometBFT on-chain events.

    Subscribes to EventDigestAttested events and orchestrates the
    IPFS fetch → SHA-256 verification → callback pipeline.

    Args:
        websocket_url: CometBFT WebSocket endpoint URL.
        ipfs_client: An IPFSClient instance for fetching content from IPFS.
        on_verified: Async callback invoked with a VerifiedDigestEvent
            after successful SHA-256 verification.
        reconnect_delay_seconds: Delay before reconnecting after a disconnect.
        max_reconnect_attempts: Maximum consecutive reconnection attempts
            before giving up. Set to 0 for unlimited retries.
    """

    def __init__(
        self,
        websocket_url: str = "ws://127.0.0.1:26657/websocket",
        ipfs_client: IPFSClient | None = None,
        on_verified: OnVerifiedCallback | None = None,
        reconnect_delay_seconds: float = 5.0,
        max_reconnect_attempts: int = 0,
    ) -> None:
        self._ws_url = websocket_url
        self._ipfs = ipfs_client or IPFSClient()
        self._on_verified = on_verified
        self._reconnect_delay = reconnect_delay_seconds
        self._max_reconnects = max_reconnect_attempts
        self._running = False

    async def start(self) -> None:
        """Start the WebSocket subscription loop with automatic reconnection.

        This method runs indefinitely, reconnecting on disconnection up to
        max_reconnect_attempts times (or forever if set to 0).
        """
        self._running = True
        consecutive_failures = 0

        while self._running:
            try:
                await self._subscribe_loop()
                consecutive_failures = 0
            except (
                websockets.exceptions.ConnectionClosed,
                websockets.exceptions.WebSocketException,
                OSError,
            ) as exc:
                consecutive_failures += 1
                logger.warning(
                    "CometBFT WebSocket disconnected (attempt %d): %s",
                    consecutive_failures,
                    exc,
                )

                if self._max_reconnects and consecutive_failures >= self._max_reconnects:
                    logger.error(
                        "Max reconnection attempts (%d) reached. Stopping subscriber.",
                        self._max_reconnects,
                    )
                    break

                logger.info("Reconnecting in %.1f seconds...", self._reconnect_delay)
                await asyncio.sleep(self._reconnect_delay)

    def stop(self) -> None:
        """Signal the subscriber to stop after the current iteration."""
        self._running = False

    async def _subscribe_loop(self) -> None:
        """Connect to CometBFT and process events until disconnection."""
        async with websockets.connect(self._ws_url) as ws:
            # Subscribe to digest attestation events
            subscribe_msg = json.dumps({
                "jsonrpc": "2.0",
                "method": "subscribe",
                "id": "casedigestor-indexer-1",
                "params": {"query": SUBSCRIBE_QUERY},
            })
            await ws.send(subscribe_msg)

            # Read the subscription confirmation
            confirmation = await ws.recv()
            logger.info("Subscribed to CometBFT events: %s", confirmation)

            # Process incoming events
            async for raw_message in ws:
                if not self._running:
                    break

                try:
                    await self._handle_message(raw_message)
                except Exception:
                    logger.exception("Error processing CometBFT event")

    async def _handle_message(self, raw_message: str | bytes) -> None:
        """Parse and process a single CometBFT WebSocket message."""
        message = json.loads(raw_message)

        # CometBFT wraps event data in result.data.value.TxResult.result.events
        result = message.get("result", {})
        if not result or "data" not in result:
            return  # Subscription confirmations or heartbeats

        events_data = result.get("data", {}).get("value", {})
        tx_result = events_data.get("TxResult", {}).get("result", {})
        events = tx_result.get("events", [])

        for event in events:
            if event.get("type") == "EventDigestAttested":
                attrs = self._parse_event_attributes(event.get("attributes", []))
                await self._process_attestation_event(attrs)

    @staticmethod
    def _parse_event_attributes(attributes: list[dict]) -> dict[str, str]:
        """Convert CometBFT event attributes to a flat dictionary."""
        result = {}
        for attr in attributes:
            key = attr.get("key", "")
            value = attr.get("value", "")
            result[key] = value
        return result

    async def _process_attestation_event(self, attrs: dict[str, str]) -> None:
        """Fetch digest from IPFS, verify SHA-256, and invoke callback."""
        digest_id = attrs.get("digest_id", "")
        ipfs_cid = attrs.get("ipfs_cid", "")
        sha256_checksum = attrs.get("sha256_checksum", "")
        attorney = attrs.get("attorney", "")

        if not all([digest_id, ipfs_cid, sha256_checksum]):
            logger.warning(
                "Incomplete attestation event attributes: digest_id=%s cid=%s",
                digest_id,
                ipfs_cid,
            )
            return

        logger.info(
            "Processing attestation event: digest_id=%s cid=%s attorney=%s",
            digest_id,
            ipfs_cid,
            attorney,
        )

        # Fetch content from IPFS
        content = await self._ipfs.cat(ipfs_cid)

        # Verify SHA-256 against the on-chain anchor
        if not verify_integrity(sha256_checksum, content):
            logger.error(
                "SHA-256 MISMATCH for digest %s! "
                "On-chain anchor: %s, IPFS content hash differs. "
                "Possible content tampering detected.",
                digest_id,
                sha256_checksum,
            )
            return

        logger.info("SHA-256 verified for digest %s: %s", digest_id, sha256_checksum)

        # Build the verified event and invoke callback
        if self._on_verified:
            verified_event = VerifiedDigestEvent(
                digest_id=digest_id,
                attorney=attorney,
                ipfs_cid=ipfs_cid,
                sha256_checksum=sha256_checksum,
                content=content,
                doctrine_accuracy_score=int(attrs.get("doctrine_accuracy_score", "0")),
                citation_integrity_score=int(attrs.get("citation_integrity_score", "0")),
                updated_stage=attrs.get("updated_stage", ""),
                total_attestations=int(attrs.get("total_attestations", "0")),
            )
            await self._on_verified(verified_event)
