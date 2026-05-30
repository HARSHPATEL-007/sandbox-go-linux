# goboxd

goboxd (Go Sandbox Daemon) is an HTTP service designed to execute untrusted, multi-language code safely within isolated nsjail sandboxes. It provides zero-code language extensibility, strict resource boundary enforcement, and deterministic queue-backed concurrency.

## Core Architecture

The service coordinates three primary steps for every `POST /run` request:
1. **Ingress Validation:** Blocks malicious compiler flags, path traversals, and oversized payloads at the HTTP layer.
2. **Bounded Queueing:** Jobs enter a finite worker pool. When capacity is reached, jobs queue rather than crashing the system or dropping connections.
3. **Isolated Execution:** Code is written to a unique volatile directory and executed via an `nsjail` process utilizing Linux namespaces, cgroups, and strict syscall restrictions.

---

## Project Structure and Document Paths

Detailed documentation lives in the `docs/` directory:

* **[`docs/architecture.md`](docs/architecture.md):** System topology, request lifecycles, and worker-pool queuing mechanics.
* **[`docs/api.md`](docs/api.md):** Complete HTTP contract specifications, payload fields, and error states.
* **[`docs/languages.md`](docs/languages.md):** Runtime configurations and instructions for writing plugin-and-play YAML definitions.
* **[`docs/security.md`](docs/security.md):** Audit logs of the 7 closed security threats and host mitigation strategies.
* **[`docs/benchmarks.md`](docs/benchmarks.md):** Latency tracking ($p_{50}$, $p_{95}$, $p_{99}$) under sustained parallel load (1 to 100 clients).

---

## Quick Start

### Prerequisites
* Docker
* Docker Compose

### Building and Running the Sandbox

All interactions are managed via the `Makefile`. Do not run bare Go commands.

1. **Build and start the sandbox container:**
   ```bash
   make docker-run