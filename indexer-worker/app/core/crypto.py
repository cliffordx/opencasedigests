"""Deterministic SHA-256 hashing utilities for Case Digestor Protocol.

Provides two canonical hashing modes:

1. **Canonical JSON hashing** — For structured case digest payloads.
   Serializes with sorted keys and UTF-8 encoding to ensure identical hashes
   regardless of key insertion order or platform-specific serialization.

2. **Raw bytes hashing** — For PDF files and binary content.
   Hashes the raw byte stream directly.

Both modes produce lowercase hex-encoded SHA-256 digests (64 characters).
"""

import hashlib
import json
from typing import Any


def canonical_json_bytes(data: dict[str, Any]) -> bytes:
    """Serialize a dictionary to canonical JSON bytes.

    Canonical form guarantees deterministic hashing by:
    - Sorting keys at all nesting levels
    - Using compact separators (no trailing whitespace)
    - Encoding as UTF-8

    Args:
        data: The dictionary to serialize. Must be JSON-serializable.

    Returns:
        UTF-8 encoded bytes of the canonical JSON representation.

    Raises:
        TypeError: If the data contains non-JSON-serializable values.
    """
    return json.dumps(
        data,
        sort_keys=True,
        separators=(",", ":"),
        ensure_ascii=False,
    ).encode("utf-8")


def sha256_hex(raw: bytes) -> str:
    """Compute the SHA-256 hash of raw bytes.

    Args:
        raw: The byte buffer to hash (canonical JSON bytes or raw PDF content).

    Returns:
        Lowercase hex-encoded SHA-256 digest (64 characters).
    """
    return hashlib.sha256(raw).hexdigest()


def hash_canonical_json(data: dict[str, Any]) -> str:
    """Compute the SHA-256 hash of a dictionary's canonical JSON form.

    This is the primary hashing function for structured case digest payloads.
    It combines canonical JSON serialization with SHA-256 hashing.

    Args:
        data: The dictionary to hash.

    Returns:
        Lowercase hex-encoded SHA-256 digest of the canonical JSON.
    """
    return sha256_hex(canonical_json_bytes(data))


def hash_raw_bytes(content: bytes) -> str:
    """Compute the SHA-256 hash of raw byte content (e.g., a PDF file).

    Args:
        content: The raw bytes to hash.

    Returns:
        Lowercase hex-encoded SHA-256 digest.
    """
    return sha256_hex(content)


def verify_integrity(expected_hash: str, actual_content: bytes) -> bool:
    """Verify that content matches its expected SHA-256 checksum.

    Used to detect tampering when content fetched from IPFS is compared
    against the on-chain SHA-256 anchor.

    Args:
        expected_hash: The expected lowercase hex SHA-256 digest (64 chars).
        actual_content: The raw bytes to verify.

    Returns:
        True if the computed hash matches the expected hash, False otherwise.
    """
    computed = sha256_hex(actual_content)
    return computed == expected_hash.lower()
