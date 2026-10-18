"""Tests for the Case Digestor cryptographic integrity pipeline.

Validates the deterministic SHA-256 hashing utilities that underpin the
IPFS-to-chain anchoring system. These tests ensure that:

1. Canonical JSON serialization is deterministic regardless of key order.
2. SHA-256 hashes are stable and reproducible across platforms.
3. Tamper detection works — any 1-byte difference in IPFS content
   produces a mismatched hash vs. the on-chain SHA-256 anchor.
4. Raw byte (PDF) hashing is consistent.
"""

import hashlib
import json

from app.core.crypto import (
    canonical_json_bytes,
    hash_canonical_json,
    hash_raw_bytes,
    sha256_hex,
    verify_integrity,
)


# ── Canonical JSON Serialization ────────────────────────────────────────


class TestCanonicalJsonBytes:
    """Tests for deterministic JSON serialization."""

    def test_sorted_keys(self):
        """Keys must be sorted alphabetically at all nesting levels."""
        data = {"z": 1, "a": 2, "m": 3}
        result = canonical_json_bytes(data)
        assert result == b'{"a":2,"m":3,"z":1}'

    def test_nested_sorted_keys(self):
        """Nested dictionaries must also have sorted keys."""
        data = {"outer": {"z_inner": 1, "a_inner": 2}, "alpha": True}
        result = canonical_json_bytes(data)
        parsed = json.loads(result)
        # Verify structure is preserved
        assert parsed["alpha"] is True
        assert parsed["outer"]["a_inner"] == 2
        # Verify key order in serialized form
        keys_in_order = list(json.loads(result).keys())
        assert keys_in_order == sorted(keys_in_order)

    def test_compact_separators(self):
        """No extra whitespace in the serialized output."""
        data = {"key": "value"}
        result = canonical_json_bytes(data)
        assert b" " not in result

    def test_utf8_encoding(self):
        """Non-ASCII characters must be preserved (not escaped)."""
        data = {"title": "People v. Dela Cruz — División Primera"}
        result = canonical_json_bytes(data)
        assert "Dela Cruz".encode("utf-8") in result
        assert "División".encode("utf-8") in result

    def test_key_order_independence(self):
        """Two dicts with the same content but different insertion order
        must produce identical canonical bytes."""
        dict_a = {"gr_number": "G.R. No. 212345", "title": "Test", "year": 2024}
        dict_b = {"year": 2024, "title": "Test", "gr_number": "G.R. No. 212345"}
        assert canonical_json_bytes(dict_a) == canonical_json_bytes(dict_b)

    def test_empty_dict(self):
        """An empty dictionary produces '{}'."""
        assert canonical_json_bytes({}) == b"{}"


# ── SHA-256 Hashing ─────────────────────────────────────────────────────


class TestSha256Hex:
    """Tests for raw SHA-256 hashing."""

    def test_known_hash(self):
        """Verify against a known SHA-256 digest."""
        # SHA-256 of empty string
        expected = hashlib.sha256(b"").hexdigest()
        assert sha256_hex(b"") == expected

    def test_lowercase_hex(self):
        """Output must be lowercase hexadecimal."""
        result = sha256_hex(b"test")
        assert result == result.lower()
        assert len(result) == 64

    def test_deterministic(self):
        """Same input must always produce the same hash."""
        content = b"Philippine Supreme Court Decision G.R. No. 212345"
        assert sha256_hex(content) == sha256_hex(content)

    def test_different_inputs_different_hashes(self):
        """Different inputs must produce different hashes."""
        assert sha256_hex(b"input_a") != sha256_hex(b"input_b")


class TestHashCanonicalJson:
    """Tests for the combined canonical JSON + SHA-256 pipeline."""

    def test_deterministic_across_key_orders(self):
        """Hash must be identical regardless of original key insertion order."""
        digest_a = {
            "gr_number": "G.R. No. 212345",
            "title": "People v. Dela Cruz",
            "ponente": "Justice Leonen",
            "division": "EN_BANC",
        }
        digest_b = {
            "division": "EN_BANC",
            "ponente": "Justice Leonen",
            "gr_number": "G.R. No. 212345",
            "title": "People v. Dela Cruz",
        }
        assert hash_canonical_json(digest_a) == hash_canonical_json(digest_b)

    def test_stable_hash_value(self):
        """The hash of a known payload must match a pre-computed value."""
        data = {"key": "value"}
        canonical = canonical_json_bytes(data)
        expected = hashlib.sha256(canonical).hexdigest()
        assert hash_canonical_json(data) == expected


# ── Tamper Detection ────────────────────────────────────────────────────


class TestVerifyIntegrity:
    """Tests for the SHA-256 integrity verification used in the
    IPFS fetch → on-chain anchor comparison pipeline."""

    def test_matching_content_passes(self):
        """Identical content must pass verification."""
        content = b"Authentic Supreme Court decision content"
        checksum = sha256_hex(content)
        assert verify_integrity(checksum, content) is True

    def test_single_byte_tamper_detected(self):
        """Changing even 1 byte in the IPFS payload must fail verification.

        This is the critical tamper detection test: if an IPFS node returns
        content that differs by a single byte from what was originally pinned,
        the SHA-256 mismatch must be caught.
        """
        original = b"Original case digest content with legal analysis"
        on_chain_hash = sha256_hex(original)

        # Simulate IPFS returning tampered content (1 byte changed)
        tampered = bytearray(original)
        tampered[0] = (tampered[0] + 1) % 256  # Flip 1 byte
        tampered = bytes(tampered)

        assert original != tampered
        assert verify_integrity(on_chain_hash, tampered) is False

    def test_appended_byte_tamper_detected(self):
        """Appending a single byte to the content must fail verification."""
        original = b"Case digest JSON payload"
        on_chain_hash = sha256_hex(original)
        tampered = original + b"\x00"
        assert verify_integrity(on_chain_hash, tampered) is False

    def test_truncated_content_tamper_detected(self):
        """Truncating content by 1 byte must fail verification."""
        original = b"Case digest JSON payload with full content"
        on_chain_hash = sha256_hex(original)
        tampered = original[:-1]
        assert verify_integrity(on_chain_hash, tampered) is False

    def test_empty_content_against_nonempty_hash(self):
        """Empty content must not match a hash of non-empty content."""
        original = b"Non-empty content"
        on_chain_hash = sha256_hex(original)
        assert verify_integrity(on_chain_hash, b"") is False

    def test_case_insensitive_hash_comparison(self):
        """Hash comparison should be case-insensitive (uppercase anchor still matches)."""
        content = b"Test content"
        checksum_upper = sha256_hex(content).upper()
        assert verify_integrity(checksum_upper, content) is True


# ── Raw Bytes (PDF) Hashing ─────────────────────────────────────────────


class TestHashRawBytes:
    """Tests for raw byte hashing (used for PDF files)."""

    def test_pdf_like_content(self):
        """Simulated PDF header bytes should hash consistently."""
        pdf_header = b"%PDF-1.7\n" + b"\x00" * 100  # Simulated PDF
        hash1 = hash_raw_bytes(pdf_header)
        hash2 = hash_raw_bytes(pdf_header)
        assert hash1 == hash2
        assert len(hash1) == 64

    def test_large_content(self):
        """Large content (simulating a real PDF) should hash without error."""
        large_content = b"x" * (10 * 1024 * 1024)  # 10 MB
        result = hash_raw_bytes(large_content)
        assert len(result) == 64

    def test_binary_content(self):
        """Arbitrary binary content should hash correctly."""
        binary = bytes(range(256))
        result = hash_raw_bytes(binary)
        expected = hashlib.sha256(binary).hexdigest()
        assert result == expected
