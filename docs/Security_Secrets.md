---
description: "Protect broker credentials and MetaRPC API keys in GoMT4: TLS, secret storage, safe logs, trial access, and demo-first verification."
---

# Protect credentials in GoMT4

GoMT4 connects your Go application to the MetaRPC MT4 endpoint over gRPC with TLS.
The broker password and MetaRPC API key serve different purposes: the password
connects a broker account; the API key identifies your MetaRPC service account.
Neither belongs in a public repository, a screenshot, or an issue report.

## Supply an API key without committing it

`NewMT4Account` accepts an API key explicitly. If that argument is empty, the SDK
checks `MRPC_API_KEY`, then falls back to `TRIAL`. An unknown key is refused;
retrying the same invalid key will not resolve the error. Trial behavior is not a
substitute for authenticating a production integration.

Inject the key from your deployment secret store or set `MRPC_API_KEY` in the
process environment. Avoid including a real value in a shell command that will
be saved to command history. The examples also support configuration files;
keep a private local copy and review the staged diff before every commit.

Treat environment variables as a delivery mechanism, not an isolation boundary:
administrators, debuggers and processes with sufficient access may read them.
Restrict access to the application identity and its configuration.

## Keep transport verification enabled

Use the production endpoint `mt4.mrpc.pro:443` with the SDK's TLS configuration.
Do not disable certificate verification to work around an endpoint or certificate
error. Check the hostname, system clock and certificate chain instead.

Each request carries `apikey` and terminal `id` metadata. Terminal identifiers
are useful for diagnosis, but they do not replace API-key authorization. Do not
publish complete request metadata or connection messages containing passwords.

## Share useful diagnostics safely

Include the SDK commit, method name, timestamp with timezone, gRPC/API error code,
and whether the account is a demo account. Redact passwords, API keys, access
tokens and unnecessary account identifiers. Avoid logging an entire account
object or configuration structure: those structures can contain credentials.

If a secret has been exposed, replace it through the relevant provider, update
your deployment and confirm that the old value no longer works. Removing a
secret from the latest Git commit does not remove it from Git history.

## Test with a disposable demo account

Begin with account data and quotes, then test order entry and cleanup on a demo
account. Inspect the code before running a trading example. A successful login
does not prove that a trade opened, closed, or had the expected size.

- [Connect your first project](All_Guides/Your_First_Project.md)
- [Diagnose connection and trading failures](Troubleshooting_FAQ.md)
- [Review the MT4Account API](MT4Account/Overview.md)
