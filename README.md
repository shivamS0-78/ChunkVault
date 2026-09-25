# ChunkVault (restic-clone)

A fast, secure, client-side encrypted, content-defined deduplicating backup engine written in Go.

Inspired by [restic](https://github.com/restic/restic), ChunkVault breaks files into variable-sized chunks using rolling Rabin fingerprints, encrypts each chunk locally with AES-256-GCM, bundles them into self-describing pack files, and stores them in pluggable storage backends.

---

## Key Features

- **Content-Defined Chunking (CDC):** Uses a Rabin rolling hash (`512 KB` to `8 MB` variable chunks) to maintain deduplication efficiency even when files shift by single bytes.
- **Client-Side Encryption (Zero-Knowledge):** Master key derived via **Argon2id** (`64 MB` memory, 3 iterations) and data authenticated & encrypted via **AES-256-GCM**. The storage backend never sees plaintext or keys.
- **Deduplication via Content Addressing:** Plaintext SHA-256 content IDs ensure identical chunks across files, directories, and snapshots are stored only once.
- **Pack Bundling:** Chunks are aggregated into ~`24 MB` pack files with binary footers to prevent inode exhaustion and eliminate cloud storage API call overhead.
- **Merkle Tree Directory Snapshots:** Subdirectories are represented as content-addressed tree blobs; unchanged directory trees deduplicate transitively with a single root hash check.
- **Pluggable Storage Backends:** Clean `Backend` storage interface (`Save`, `Load`, `Stat`, `List`, `Delete`) with local disk support (atomic write + directory fan-out) and planned S3/cloud support.
- **Atomic Commits & Crash Resilience:** All data blobs are immutable; snapshot manifests act as single commit points to prevent repository corruption on crashes.

---

## Architecture Overview

Every backup operation flows top-to-bottom through a modular five-stage pipeline:

```text
CLI (Cobra Commands)
        │
        ▼
Tree & Snapshot ─────────────  Directory walk, Merkle tree manifests
        │
        ▼
Repository Core ─────────────  Orchestrates pipeline & deduplication
        │
        ├───────────┬────────────┬─────────────┐
        ▼           ▼            ▼             ▼
     Chunker      Crypto       Packer        Index
   Rabin CDC    AEAD + KDF    Bundling    BoltDB Map
        │           │            │             │
        └───────────┴────────────┴─────────────┘
                         │
                         ▼
                Storage Backend (Pluggable)
                  ├── Local (Atomic + Fan-out)
                  └── S3 (Cloud Object Store)
```

---

## Repository On-Disk Layout

```text
repo/
├── config.json              # Unencrypted KDF parameters & encrypted master key
├── data/
│   ├── ab/
│   │   └── ab3f...pack      # Encrypted pack files fanned out by prefix
│   └── cd/
│       └── cd91...pack
├── index/
│   └── index.boltdb         # Local lookup cache (rebuildable from pack footers)
└── snapshots/
    └── 2026-08-15T09-12-00Z # Snapshot manifests (root tree ID, timestamp, metadata)
```

---

## Project Structure

```text
restic-clone/
├── cmd/
│   └── backup/              # CLI entry point (Cobra commands)
├── internal/
│   ├── backend/             # Storage backend abstraction
│   │   └── local/           # Local filesystem backend implementation
│   ├── chunker/             # Content-defined chunking (Rabin polynomial)
│   ├── crypto/              # Argon2id KDF, AES-256-GCM AEAD, SHA-256 IDs
│   ├── pack/                # Pack file builder & index footer serializer
│   ├── index/               # Chunk ID → pack offset mapping (BoltDB)
│   ├── repository/          # High-level pipeline coordinator
│   ├── snapshot/            # Snapshot manifests & metadata
│   └── tree/                # Directory walker & Merkle tree reconstruction
├── go.mod
├── sysDesign.md             # Detailed system design specification
└── README.md
```

---

## Core Components

| Component | Package | Responsibility |
|---|---|---|
| **Crypto** | `internal/crypto` | Argon2id key derivation, AES-256-GCM authenticated encryption, SHA-256 content IDs. |
| **Chunker** | `internal/chunker` | Rolling-hash chunker splitting byte streams at content boundaries. |
| **Packer** | `internal/pack` | Buffers chunks into ~24 MB packs and attaches self-describing footers. |
| **Backend** | `internal/backend` | Opaque blob storage interface (`Local` filesystem implemented, `S3` planned). |
| **Repository** | `internal/repository` | Coordinates chunking, dedup checking, encryption, and pack writing. |
| **Snapshot/Tree** | `internal/snapshot`, `internal/tree` | Traverses filesystem, builds tree manifests, and restores snapshot states. |

---

## Getting Started

### Prerequisites

- **Go 1.22+** (or later)

### Installation

Clone the repository and download dependencies:

```bash
git clone https://github.com/shivamS0-78/ChunkVault.git
cd ChunkVault
go mod download
```

### Running Tests

Run the test suite across all packages:

```bash
go test -v ./...
```

---

## License

MIT License
