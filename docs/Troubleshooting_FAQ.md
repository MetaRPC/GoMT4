---
description: "Diagnose GoMT4 connection errors, rejected API keys, missing quotes, unknown trade outcomes and stream interruptions without duplicating an order."
---

# Diagnose a GoMT4 integration

Work through the connection, account, symbol and order result separately. A request
reaching MetaRPC does not mean that the broker accepted it, and a timeout does not
prove that the broker did nothing. Start on a disposable demo account.

## Authentication and connection

Check the endpoint (`mt4.mrpc.pro:443`), API key, MT4 broker login, broker server
and password. A MetaTrader 5 login cannot authenticate an MT4 terminal. Confirm
that you used the trading password when the test requires order operations.

`API_KEY_INVALID` or `API_KEY_REQUIRED` needs a credential/configuration correction;
the client does not retry a refused API key. Connect before reading the account.
After connection, the SDK uses the terminal instance GUID returned by the server.
See [Getting started](All_Guides/GETTING_STARTED.md) and
[credential handling](Security_Secrets.md).

## Quotes and symbol names

Use the exact symbol offered by the broker, including suffixes. Select it in
Market Watch and check whether a fresh quote exists. A connected terminal can
still have a closed market or an unavailable instrument.

- [List symbols](MT4Account/3.%20Symbols_Quotes/Symbols.md)
- [Select a symbol](MT4Account/3.%20Symbols_Quotes/SymbolSelect.md)
- [Inspect contract and volume rules](MT4Account/3.%20Symbols_Quotes/SymbolParams.md)
- [Read a quote](MT4Account/3.%20Symbols_Quotes/Quote.md)

## Timeouts and uncertain trade outcomes

Read calls can retry eligible transient failures with bounded backoff and respect
their context deadline. Trade operations are different: the current SDK does not
automatically replay order send, modify, close or close-by calls after an unknown
outcome. Blindly retrying can create a second trade or repeat a modification.

After a timeout or lost connection, inspect open orders and history before
deciding what to do next. Reconcile using the ticket when known, plus symbol,
direction, volume and time. Magic/comment values can help diagnosis but are not
an idempotency guarantee. Do not implement an outer retry loop that defeats this
protection.

- [Open orders](MT4Account/2.%20Orders_History/OpenedOrders.md)
- [Order history](MT4Account/2.%20Orders_History/OrdersHistory.md)
- [Broker return codes](All_Guides/RETURN_CODES_REFERENCE.md)

## Streams and reconnects

Cancel streams when leaving the screen or shutting down the application. If a
stream stops, inspect its error and connection state, reconnect as appropriate,
then reconcile account state before acting on new events. Do not assume that
every event was delivered while disconnected. See
[gRPC stream management](All_Guides/GRPC_STREAM_MANAGEMENT.md).

## Logs and health checks

Record method, duration, error code, attempt count and a correlation identifier.
Exclude credentials and full request metadata. A successful account-summary read
is a useful connection check; it does not establish quote freshness, market hours
or permission to trade. Track those separately.

For a reproducible report include the SDK commit, UTC time, demo/live mode,
method and redacted error. Compare an account-summary request, one quote request
and the specific failing operation instead of submitting a complete secret-bearing
terminal dump.

## Latency and resource use

Measure connection time separately from subsequent API calls. Reuse the client
for a session, set operation deadlines, cancel unused subscriptions and avoid
unbounded polling/retry loops. Use bulk quote methods where appropriate. Do not
increase retry counts to hide a persistent credential or broker rejection.

## Example configuration and cleanup

Run the examples from the repository workspace; the interactive menu groups them
by SDK layer. Confirm the selected example's trading behavior and account before
running it. On completion, close or cancel the demo orders you created and use
the SDK's disconnect/connection cleanup path. See the
[first project](All_Guides/Your_First_Project.md) and
[project layout](All_Guides/PROJECT_MAP.md).
