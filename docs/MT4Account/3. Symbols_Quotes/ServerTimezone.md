# ServerTimezone

> The broker server's UTC offset — how far server time is from UTC, and how that was determined.

**API Information:**

* **Low-level API:** `MT4Account.ServerTimezone(ctx)` (defined in [`mt4/MT4Account.go`](https://github.com/MetaRPC/GoMT4/blob/main/mt4/MT4Account.go))
* **gRPC service:** `mt4_term_api.ServerTimezone`
* **RPC:** `ServerTimezone(ServerTimezoneRequest) → ServerTimezoneReply`
* **Payload:** `ServerTimezoneData`

## 💬 Just the essentials

* **What it is.** The server's offset from UTC (in seconds) plus its provenance.
* **Why you need it.** To convert `*_utc` vs server timestamps and align history correctly.
* **When to use.** When displaying/aligning times, or reconciling server vs UTC fields.

---

## 🎯 Purpose

* Read the broker's UTC offset in seconds
* Know how the offset was determined (live vs last-validated)
* Convert between server time and UTC consistently

---

## Method Signature

```go
func (a *MT4Account) ServerTimezone(ctx context.Context) (*pb.ServerTimezoneData, error)
```

**Request message:**

```protobuf
message ServerTimezoneRequest {}   // no arguments
```

**Payload message:**

```protobuf
message ServerTimezoneData {
  string server = ...;
  int32  utc_offset_seconds = ...;
  SERVER_TIMEZONE_SOURCE source = ...;
  Timestamp validated_utc = ...;
}
```

---

## 🔽 Input

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `ctx` | `context.Context` | yes | Cancellation / timeout. The call takes no other arguments. |

---

## ⬆️ Output — `*pb.ServerTimezoneData`

| Field | Proto Type | Go getter | Description |
|-------|-----------|-----------|-------------|
| `server` | string | `GetServer() string` | Broker server name |
| `utc_offset_seconds` | int32 | `GetUtcOffsetSeconds() int32` | Offset from UTC in seconds (positive = ahead of UTC) |
| `source` | SERVER_TIMEZONE_SOURCE | `GetSource()` | How the offset was determined |
| `validated_utc` | Timestamp | `GetValidatedUtc()` | When the offset was last validated |

> On failure the call returns a non-nil `error` (text `"API error (code=…): …"`) — see the [Return Codes Reference](../../All_Guides/RETURN_CODES_REFERENCE.md).

---

## 🧱 Related enums (from proto)

### `SERVER_TIMEZONE_SOURCE` — field `source`

| Proto value | Value | Meaning |
|-------------|-------|---------|
| `SERVER_TIMEZONE_SOURCE_UNKNOWN` | 0 | Not yet determined |
| `SERVER_TIMEZONE_SOURCE_LIVE` | 1 | Learned from a live quote this session |
| `SERVER_TIMEZONE_SOURCE_LAST_VALIDATED` | 2 | From the last validated value |

> In Go the generated constant repeats the prefix: `pb.SERVER_TIMEZONE_SOURCE_SERVER_TIMEZONE_SOURCE_LIVE`. See the [ENUMs Reference](../../All_Guides/ENUMS_USAGE_REFERENCE.md).

---

## 🧩 Notes & Tips

* **Offset is informational.** The API already provides `*_utc` fields on time-bearing payloads; you rarely need to convert by hand.
* **Multiple of 15 minutes.** Broker offsets are validated to 15-minute steps.
* **Source matters.** `UNKNOWN` means the terminal hasn't learned the offset yet (e.g. market closed since start).

---

## 🔗 Usage Examples

### 1) Print the offset

```go
func tz(ctx context.Context, acc *mt4.MT4Account) error {
    d, err := acc.ServerTimezone(ctx)
    if err != nil { return err }
    fmt.Printf("%s  UTC%+.1fh  (source %v)\n",
        d.GetServer(), float64(d.GetUtcOffsetSeconds())/3600.0, d.GetSource())
    return nil
}
```

### 2) Convert a server time to UTC

```go
func toUTC(ctx context.Context, acc *mt4.MT4Account, serverTime time.Time) (time.Time, error) {
    d, err := acc.ServerTimezone(ctx)
    if err != nil { return time.Time{}, err }
    return serverTime.Add(-time.Duration(d.GetUtcOffsetSeconds()) * time.Second), nil
}
```

### 3) Guard on an unknown offset

```go
import pb "github.com/MetaRPC/GoMT4/package"

func offsetKnown(ctx context.Context, acc *mt4.MT4Account) (bool, error) {
    d, err := acc.ServerTimezone(ctx)
    if err != nil { return false, err }
    return d.GetSource() != pb.SERVER_TIMEZONE_SOURCE_SERVER_TIMEZONE_SOURCE_UNKNOWN, nil
}
```

---

## ⚠️ Pitfalls

* **`UNKNOWN` source.** Don't convert times when the offset isn't known yet — wait for a live quote.
* **DST steps.** History from before the broker's last DST change can be off by the DST step — the current offset is applied uniformly.

---

## 🧪 Testing Suggestions

* **Happy path (demo, market open).** `GetSource()` is `LIVE`; `GetUtcOffsetSeconds()` is a multiple of 900.
* **Conversion.** `toUTC` of a server time equals the matching `*_utc` field on another payload.
* **Unknown.** Right after connect on a closed market, `offsetKnown` may be false.

---

## 📚 See also

* [AccountSummary](../1.%20Account_Information/AccountSummary.md) · [QuoteHistory](./QuoteHistory.md)
* **Back to:** [3. Symbols & Quotes — Overview](./Symbols_Quotes.Overview.md) · [MT4Account — Overview](../Overview.md)
