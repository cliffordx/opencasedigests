"""Async IPFS (Kubo) RPC client for the Case Digestor Protocol.

Provides content-addressed pinning of Supreme Court decisions (PDFs) and
structured case digest JSON payloads to a local or remote Kubo IPFS node.
All content is pinned with CIDv1 addressing and returns both the IPFS CID
and the SHA-256 checksum for on-chain anchoring.
"""

import logging
from dataclasses import dataclass

import httpx

from app.core.crypto import canonical_json_bytes, sha256_hex

logger = logging.getLogger(__name__)


@dataclass(frozen=True)
class PinResult:
    """Result of an IPFS pin operation.

    Attributes:
        ipfs_cid: The CIDv1 content identifier assigned by IPFS.
        sha256_checksum: Hex-encoded SHA-256 hash of the raw bytes that were pinned.
        size_bytes: Size of the pinned content in bytes.
    """

    ipfs_cid: str
    sha256_checksum: str
    size_bytes: int


class IPFSClient:
    """Async client for the Kubo IPFS HTTP RPC API.

    Connects to a Kubo node's RPC endpoint (default: http://127.0.0.1:5001)
    and provides methods to pin Supreme Court decisions and case digests
    with CIDv1 content addressing.

    Args:
        kubo_rpc_url: Base URL of the Kubo IPFS RPC API.
        timeout_seconds: Request timeout in seconds.
    """

    def __init__(
        self,
        kubo_rpc_url: str = "http://127.0.0.1:5001",
        timeout_seconds: float = 60.0,
    ) -> None:
        self._base_url = kubo_rpc_url.rstrip("/")
        self._timeout = timeout_seconds

    async def pin_raw_bytes(self, content: bytes, filename: str = "document") -> PinResult:
        """Pin raw byte content (e.g., a PDF) to IPFS with CIDv1 addressing.

        Args:
            content: The raw bytes to pin.
            filename: Optional filename for the IPFS object.

        Returns:
            PinResult with the CIDv1 identifier and SHA-256 checksum.

        Raises:
            httpx.HTTPStatusError: If the Kubo RPC returns a non-2xx status.
        """
        checksum = sha256_hex(content)

        async with httpx.AsyncClient(timeout=self._timeout) as client:
            response = await client.post(
                f"{self._base_url}/api/v0/add",
                params={
                    "cid-version": "1",
                    "pin": "true",
                    "quieter": "true",
                },
                files={"file": (filename, content)},
            )
            response.raise_for_status()

        result = response.json()
        cid = result["Hash"]
        size = int(result.get("Size", len(content)))

        logger.info(
            "Pinned %s to IPFS: cid=%s sha256=%s size=%d",
            filename,
            cid,
            checksum,
            size,
        )

        return PinResult(ipfs_cid=cid, sha256_checksum=checksum, size_bytes=size)

    async def pin_json(self, data: dict, filename: str = "digest.json") -> PinResult:
        """Pin a structured JSON payload to IPFS using canonical serialization.

        The JSON is serialized in canonical form (sorted keys, compact separators,
        UTF-8 encoding) to ensure deterministic content addressing.

        Args:
            data: The dictionary to serialize and pin.
            filename: Optional filename for the IPFS object.

        Returns:
            PinResult with the CIDv1 identifier and SHA-256 checksum.

        Raises:
            httpx.HTTPStatusError: If the Kubo RPC returns a non-2xx status.
        """
        raw = canonical_json_bytes(data)
        return await self.pin_raw_bytes(raw, filename=filename)

    async def cat(self, cid: str) -> bytes:
        """Retrieve content from IPFS by CID.

        Args:
            cid: The IPFS content identifier to fetch.

        Returns:
            The raw bytes of the IPFS object.

        Raises:
            httpx.HTTPStatusError: If the Kubo RPC returns a non-2xx status.
        """
        async with httpx.AsyncClient(timeout=self._timeout) as client:
            response = await client.post(
                f"{self._base_url}/api/v0/cat",
                params={"arg": cid},
            )
            response.raise_for_status()

        return response.content
