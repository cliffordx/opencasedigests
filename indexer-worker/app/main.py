"""Case Digestor Indexer Worker — FastAPI Application Entrypoint.

Provides the IPFS pinning, SHA-256 verification, and CometBFT WebSocket
indexing microservice for the Case Digestor Protocol hybrid anchoring pipeline.
"""

from fastapi import FastAPI

app = FastAPI(
    title="Case Digestor Indexer Worker",
    description=(
        "Hybrid anchoring microservice for IPFS pinning, SHA-256 verification, "
        "and CometBFT WebSocket-driven indexing of Philippine Supreme Court "
        "case digests and peer attestations."
    ),
    version="0.1.0",
)


@app.get("/healthz")
async def healthz() -> dict[str, str]:
    """Health check endpoint for container orchestration and monitoring."""
    return {"status": "ok", "service": "casedigestor-indexer-worker"}
