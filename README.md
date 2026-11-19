# GoMT4 SDK

## Documentation checks and previews

Use Python 3.12 and the pinned documentation dependencies:

```sh
python -m pip install -r requirements-docs.txt
mkdocs build --strict
python scripts/verify-doc-routes.py
```

The route check verifies the generated HTML for 61 moved documentation URLs,
their canonical destinations, and the restored security and troubleshooting
guides. Preserve old public URLs when reorganizing documentation: add a specific
mapping in `mkdocs.yml`, or retain useful content at its original path. Do not
redirect unrelated guides to the homepage.

The **Documentation preview checks** workflow builds pull requests and the
`codex/traffic-doc-route-recovery` review branch, then uploads a downloadable
`gomt4-documentation-preview` artifact. It does not deploy GitHub Pages.
Production Pages deployment remains in the separate main/master workflow.

> Official MetaTrader 4 Go SDK for high-performance trading automation and market data streaming via gRPC.

[![Docs](https://img.shields.io/badge/docs-GoMT4-blue.svg)](https://metarpc.github.io/GoMT4/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

📄 **Full Documentation**: [https://metarpc.github.io/GoMT4/](https://metarpc.github.io/GoMT4/)

---

## 📦 Installation

```bash
# the SDK (MT4Account / MT4Service / MT4Sugar); pulls the generated protobuf package too
go get github.com/MetaRPC/GoMT4/mt4

# (the generated protobuf types live in a sibling module)
go get github.com/MetaRPC/GoMT4/package
```

---

## 🏃 How to Run Examples

Clone the repository and run the examples out-of-the-box:

```bash
git clone https://github.com/MetaRPC/GoMT4.git
cd GoMT4/examples
go mod download

# 1. Run with default TRIAL key:
go run main.go

# 2. Or pass your MetaRPC API key directly as an argument:
go run main.go your_api_key_here

# 3. Or use the MRPC_API_KEY environment variable:
export MRPC_API_KEY="your_api_key_here"        # Windows CMD: set MRPC_API_KEY=your_api_key_here
go run main.go                                 # Windows PowerShell: $env:MRPC_API_KEY="your_api_key_here"
```

---

## 🔑 API Key & Authentication

Connecting to MetaRPC production endpoints (`mt4.mrpc.pro:443`) requires an API key:

1. **Sign Up**: Create an account for free at [https://mrpc.pro/signup](https://mrpc.pro/signup).
2. **Generate API Key**: In your MetaRPC Portal dashboard at [https://mrpc.pro/my](https://mrpc.pro/my), go to **API Keys** to generate and copy your personal API token.
3. **Configure Connection**: Pass your API key / token along with the server address (`mt4.mrpc.pro:443`) in your connection settings.

---

## 🆔 Automatic Account ID & Authentication

MetaRPC endpoints route calls using a terminal session identifier (`id`):
- When connecting via `ConnectByServerName` / `ConnectByHostPort` (which issue the `ConnectEx` RPC under the hood), the server automatically generates and returns a session GUID (`terminalInstanceGuid`).
- The SDK automatically captures this session ID and attaches both the `id` and `APIKey` headers on all subsequent calls (`AccountSummary`, `OrderSend`, streaming, etc.) — no manual `GetId` or `curl` calls required.
- Pass your API key directly to the Account/Client constructor or via the `MRPC_API_KEY` environment variable.

## 🌐 Production Endpoints

| Environment | Host | Port | Protocol |
| :--- | :--- | :--- | :--- |
| **Production** | `mt4.mrpc.pro` | `443` | TLS / gRPC |
| **Direct API UI (Swagger)** | `https://mt4.mrpc.pro/apiui` | `443` | HTTPS |
| **Portal Dashboard** | `https://mrpc.pro/my` | `443` | HTTPS |
| **Registration / API Key** | `https://mrpc.pro/signup` | `443` | HTTPS |

---

## 📄 Documentation & Guides

Explore comprehensive documentation at [https://metarpc.github.io/GoMT4/](https://metarpc.github.io/GoMT4/):
- 🚀 **Quick Start & First Project**
- 🔑 **Authentication & API Keys**
- 📡 **Live Market Data & gRPC Streaming**
- 💼 **Account Management & Order Execution**
- 📊 **Return Codes & Error Handling Reference**
