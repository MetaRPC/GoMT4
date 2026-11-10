# Protobuf Inspector — Interactive Type Explorer

> An offline developer tool that explores the MT4 gRPC protobuf **types, fields and enums** — search, inspect message structures, trace enum usage, and export to JSON, all **without connecting to MT4**.

It reads the descriptors the generated `package` registers at startup (Go's protobuf reflection), so it always matches the exact types your code compiles against.

---

## 🏁 Running it

From the `examples` module:

```bash
cd examples

go run ./protobuf_inspector               # interactive REPL
go run ./protobuf_inspector OpenedOrderInfo   # one-shot: inspect a type
go run ./protobuf_inspector stats             # one-shot: statistics
```

Interactive prompt:

```
===========================================================
MT4 PROTOBUF TYPES INSPECTOR
===========================================================

  Discovered 255 protobuf message types
  Discovered 26 enum types

  Type 'help' to see available commands
  Type 'list' to see all types
===========================================================

>
```

---

## 🧰 Commands

| Command | Description | Example |
|---------|-------------|---------|
| `list` / `ls` | All message types | `list` |
| `list --request` / `--reply` | Only `*Request` / `*Reply` types | `list --reply` |
| `<TypeName>` | Inspect a type's fields | `OpenedOrderInfo` |
| `search <text>` / `find <text>` | Search types by name | `search Order` |
| `field <name>` | Types containing a field | `field ticket` |
| `findtype <type>` | All fields of a given type | `findtype double` |
| `enum <name>` | An enum's values | `enum OrderSendOperationType` |
| `whereenum <name>` | Where an enum is used | `whereenum OpenedOrderType` |
| `export <TypeName>` | Export a type to JSON | `export OrderSendData` |
| `stats` | Statistics | `stats` |
| `help` / `?` | Help | `help` |
| `exit` / `quit` / `q` | Exit | `exit` |

---

## Practical examples

### 1) Inspect a type's fields

```
> OpenedOrderInfo
```
```
TYPE: OpenedOrderInfo
===========================================================
FIELDS (19):
-----------------------------------------------------------
  #  1 comment                          : string
  #  2 commision                        : double
  #  3 expiration_time                  : Timestamp
  #  4 lots                             : double
  #  5 magic_number                     : int32
  #  6 open_price                       : double
  #  7 profit                           : double
  # 12 ticket                           : int32
  # 14 order_type                       : OpenedOrderType
  ...
```

The field number is the protobuf tag; the type is the Go/proto type (`repeated …` and `map<…>` are shown where relevant). This is exactly what the generated getters return — `GetTicket()`, `GetOrderType()`, etc.

### 2) Which types have a `ticket` field?

```
> field ticket
```
```
Types containing a field like "ticket":
-----------------------------------------------------------
  HistoryOrderInfo                   .ticket : int32
  OnOpenedOrdersProfitOrderInfo      .ticket : int32
  OnOpenedOrdersTicketsData          .position_tickets : repeated int32
  OnTradeOrderInfo                   .ticket : int32
  OpenedOrderInfo                    .ticket : int32
```

### 3) See an enum's values

```
> enum OrderSendOperationType
```
```
ENUM: OrderSendOperationType (6 values)
-----------------------------------------------------------
     0  OC_OP_BUY
     1  OC_OP_SELL
     2  OC_OP_BUYLIMIT
     3  OC_OP_SELLLIMIT
     4  OC_OP_BUYSTOP
     5  OC_OP_SELLSTOP
```

### 4) Where is an enum used?

```
> whereenum OpenedOrderType
```
```
"OpenedOrderType" is used by:
-----------------------------------------------------------
  HistoryOrderInfo                   .order_type
  OpenedOrderInfo                    .order_type
```

### 5) Export a type to JSON (for docs/tooling)

```
> export OrderSendData
```
```json
{
  "type": "OrderSendData",
  "fields": [
    { "number": 1, "name": "ticket", "type": "int32" },
    { "number": 2, "name": "volume", "type": "double" },
    ...
  ]
}
```

---

## 🧩 How it works

- A blank import of the generated package (`_ "github.com/MetaRPC/GoMT4/package"`) registers every message and enum descriptor into protobuf's global registry on `init()`.
- The inspector walks `protoregistry.GlobalFiles`, collecting descriptors from the MT4 API packages (`mt4_term_api`, `mrpc_admin`) and skipping the protobuf/gRPC runtime's own types (`google.*`, `grpc.*`).
- Everything is **static reflection** — no terminal connection, no network.

Source: [`examples/protobuf_inspector/main.go`](https://github.com/MetaRPC/GoMT4/blob/main/examples/protobuf_inspector/main.go).

---

**See also:** [MT4Account API](../MT4Account/Overview.md) · [Glossary](GLOSSARY.md)
