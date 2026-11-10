# MT4 for Beginners: Your First Steps

Welcome! This guide is for complete beginners who have never used MetaTrader 4. We'll set up the terminal and a free demo account, then point GoMT4 at it.

---

## What is MetaTrader 4 (MT4)?

**MetaTrader 4** is a widely-used trading platform. It lets you:

- Trade forex, metals, commodities and more
- Analyze markets with charts and indicators
- Run automated strategies
- Monitor positions in real time

Think of MT4 as the "cockpit" for the markets. **GoMT4** lets you control that cockpit programmatically from Go, over gRPC.

---

## Demo vs Live: start safely

### Why start with a demo account?

A **demo account** is a risk-free practice account with virtual money. We **strongly recommend** it while you learn:

- ✅ **Zero risk** — virtual funds (typically $10,000–$100,000)
- ✅ **Real market conditions** — live prices and behavior
- ✅ **Learn the platform** before risking real money
- ✅ **Test your code** — perfect for developing GoMT4 bots

### Live account

A **live account** uses real money. Only move to live when you understand MT4, have tested your strategy thoroughly on demo, and accept the risks.

---

## Quick start: creating a demo account

The fastest path is to create a demo account directly in the MT4 terminal.

### Step 1️⃣ — Download and install the MT4 terminal

1. Visit the official MetaQuotes site: [https://www.metatrader4.com/en/download](https://www.metatrader4.com/en/download)
2. Download the build for your OS.
3. Install and launch it.

---

### Step 2️⃣ — Launch MT4 and find the Navigator panel

On first launch MT4 usually shows a popup to open a new account. If it doesn't, find the **Navigator** panel.

![Navigator Panel Location](../Guide_Images/1.%20NAV.png)

The **Navigator** is on the **left**, below **Market Watch**. It has sections: Accounts, Indicators, Expert Advisors, and others.

---

### Step 3️⃣ — Open a new demo account

1. In **Navigator**, find **Accounts**.
2. **Right-click** "Accounts".
3. Choose **"Open an account"**.

This opens the broker server search.

---

### Step 4️⃣ — Select the MetaQuotes demo server

In the broker search, for your first account use the simplest option:

- Look for **"MetaQuotes-Demo"** / **"MetaQuotes Ltd"** (use the search field if needed).
- Select it and click **Next**.

![MetaQuotes Server Selection](../Guide_Images/2.%20META.png)

---

### Step 5️⃣ — Choose "Open a new demo account"

![Open New Demo Account](../Guide_Images/3.%20OPEN.png)

Select **"Open a new demo account"** and click **Next**.

---

### Step 6️⃣ — Fill in your details

- **Name / Last Name / Email**
- **Phone** — optional for demo
- **Deposit** — virtual amount (e.g. $10,000)
- **Leverage** — 1:100 is a good start
- **Currency** — USD, EUR, …

Accept the terms and click **Next**.

![Personal Information Form](../Guide_Images/4.%20creator.png)

---

### Step 7️⃣ — Save your credentials

MT4 creates the account and shows:

- **Login** — account number (e.g. 591129415)
- **Password** — master (trading) password
- **Investor Password** — read-only
- **Server** — e.g. "MetaQuotes-Demo"

**⚠️ Save these now.** You'll need **Login**, **Password** and **Server** for GoMT4.

![Account Credentials](../Guide_Images/5.%20OKK.png)

Click **Finish**.

---

### Step 8️⃣ — Configure GoMT4 with your credentials

Open **`examples/config/config.json`** and fill in what you saved:

```json
{
  "Login": 591129415,
  "Password": "YourPasswordHere",
  "Server": "MetaQuotes-Demo",
  "DefaultSymbol": "EURUSD",
  "ApiKey": "your-api-key"
}
```

| Key | Meaning |
|-----|---------|
| `Login` | Your login number (`0` = auto-provision a fresh demo on `Server`) |
| `Password` | Your master password |
| `Server` | Exact broker server name from MT4 (e.g. `MetaQuotes-Demo`) |
| `DefaultSymbol` | Default symbol the examples trade |
| `ApiKey` | Your MetaRPC key, or `"TRIAL"` / empty to use the `MRPC_API_KEY` env var |

> Prefer env vars? `MRPC_API_KEY` is picked up when `ApiKey` is empty/`TRIAL`. Never commit real credentials.

The gRPC endpoint (`mt4.mrpc.pro:443`) is built into the SDK — you don't set a host/port in the file.

---

### Step 9️⃣ — Verify the MT4 connection

Back in the terminal:

1. MT4 logs you into the new demo account automatically.
2. Bottom-right corner should show a **green connection indicator** and your **balance** (e.g. $10,000).

You're ready to test.

---

### Step 🔟 — Run GoMT4

With `config.json` filled in, run the examples from a clone:

```bash
cd GoMT4/examples
go run .            # interactive menu of grouped examples
```

Pick `11` to read account info, `31` for the MT4Sugar basics, and so on. Trading demos are off by default — enable them (demo only!) with `MT4_ENABLE_TRADING=1`.

Trades GoMT4 sends appear in your MT4 terminal in real time, and the demo balance moves accordingly. **Experiment freely — it's virtual money.**

---

## 🧩 MT4 password types

### Master password
- **Full access** — open/close trades, change settings.
- **This is the one GoMT4 uses** (`Password` in `config.json`).

### Investor password
- **View-only** — see trades/balance/history, **cannot trade**.
- Good for sharing performance or monitoring tools.
- Set it in **Tools → Options → Server → Change** (Investor Password).

---

## Choosing a broker (optional)

The MetaQuotes demo is fine for practice. Later you might pick a broker for live trading, better spreads, or regional support.

**Not financial advice — do your own research.** Popular MT4 brokers include IC Markets, Pepperstone, OANDA, XM, RoboForex, FxPro.

What to look for:

- ✅ **Regulation** (FCA, ASIC, CySEC, …)
- ✅ **MT4 support** (specifically MT4, not only MT5)
- ✅ **Free demo** (note: demos often expire after ~30 days of inactivity)
- ✅ **Good spreads** and responsive **support**

---

## What's next?

### 1. Connect GoMT4
👉 **[Getting Started with GoMT4](GETTING_STARTED.md)** · **[Your First Project](Your_First_Project.md)**

### 2. Understand the architecture
GoMT4 has three layers — start wherever fits:

- **[MT4Account](../MT4Account/Overview.md)** — raw gRPC, maximum control.
- **[MT4Service](../MT4Service/Overview.md)** — friendly wrappers, most common path.
- **[MT4Sugar](../MT4Sugar/Overview.md)** — one-liners, most beginner-friendly.

---

## Important reminders

### 🛡️ Security
- Never share your master password. Use the investor password for monitoring tools. Keep credentials in a password manager.

### 📊 Risk management
Even on demo, build good habits: always use a stop loss, risk no more than 1–2% per trade, understand position sizing (`MT4Sugar.CalcLotByRisk`).

### 🧪 Testing
Test all code on demo first. Verify strategies over weeks, not days. Demo results don't guarantee live results.

---

## Useful commands

```bash
cd examples
go run .                 # interactive example launcher
go run ./protobuf_inspector   # offline type explorer
go build ./...           # compile everything
go version               # check your Go toolchain
```

---

**Ready to code?** Head to [Getting Started](GETTING_STARTED.md) and connect GoMT4 to your demo account.

**Good luck! 🚀**
