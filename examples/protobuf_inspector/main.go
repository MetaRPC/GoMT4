// Command protobuf_inspector is an interactive, offline explorer for the MT4
// gRPC protobuf types. It reads the descriptors registered by the generated
// package (no MT4 connection needed) and lets you search types, fields and
// enums, inspect message structures, and export them to JSON.
//
// Run it from the examples module:
//
//	go run ./protobuf_inspector            # interactive REPL
//	go run ./protobuf_inspector OpenedOrderInfo   # one-shot: inspect a type
//	go run ./protobuf_inspector stats
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	// Blank import: the generated package registers every message/enum
	// descriptor into the global registry on init, which is all we read.
	_ "github.com/MetaRPC/GoMT4/package"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type inspector struct {
	msgs  map[string]protoreflect.MessageDescriptor
	enums map[string]protoreflect.EnumDescriptor
}

func main() {
	in := load()

	// One-shot mode: `go run . <command ...>` runs a single command and exits.
	if len(os.Args) > 1 {
		in.run(strings.Join(os.Args[1:], " "))
		return
	}

	in.header()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return // EOF
		}
		cmd := strings.TrimSpace(line)
		switch cmd {
		case "exit", "quit", "q":
			return
		case "":
			continue
		default:
			in.run(cmd)
		}
	}
}

// ---- discovery ----

func load() *inspector {
	in := &inspector{
		msgs:  map[string]protoreflect.MessageDescriptor{},
		enums: map[string]protoreflect.EnumDescriptor{},
	}
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if !isMT4Package(string(fd.Package())) {
			return true // skip google.*, grpc.* and other runtime descriptors
		}
		in.collectMessages(fd.Messages())
		in.collectEnums(fd.Enums())
		return true
	})
	return in
}

// isMT4Package keeps only the MetaRPC MT4 API proto packages, filtering out the
// protobuf/gRPC runtime's own descriptor types.
func isMT4Package(pkg string) bool {
	return strings.HasPrefix(pkg, "mt4_term_api") || strings.HasPrefix(pkg, "mrpc")
}

func (in *inspector) collectMessages(mds protoreflect.MessageDescriptors) {
	for i := 0; i < mds.Len(); i++ {
		md := mds.Get(i)
		if md.IsMapEntry() {
			continue // synthetic map-entry messages are noise
		}
		in.msgs[string(md.Name())] = md
		in.collectMessages(md.Messages()) // nested types
		in.collectEnums(md.Enums())       // nested enums
	}
}

func (in *inspector) collectEnums(eds protoreflect.EnumDescriptors) {
	for i := 0; i < eds.Len(); i++ {
		ed := eds.Get(i)
		in.enums[string(ed.Name())] = ed
	}
}

// ---- command dispatch ----

func (in *inspector) run(cmd string) {
	parts := strings.Fields(cmd)
	head := parts[0]
	arg := strings.TrimSpace(strings.TrimPrefix(cmd, head))

	switch head {
	case "help", "?":
		in.help()
	case "list", "ls":
		in.list(arg)
	case "search", "find":
		in.search(arg)
	case "field":
		in.field(arg)
	case "findtype":
		in.findType(arg)
	case "enum":
		in.enum(arg)
	case "whereenum":
		in.whereEnum(arg)
	case "export":
		in.export(arg)
	case "stats":
		in.stats()
	default:
		// Bare word → treat it as a type name to inspect.
		in.inspect(head)
	}
}

// ---- commands ----

func (in *inspector) header() {
	line := strings.Repeat("=", 59)
	fmt.Println(line)
	fmt.Println("MT4 PROTOBUF TYPES INSPECTOR")
	fmt.Println(line)
	fmt.Printf("\n  Discovered %d protobuf message types\n", len(in.msgs))
	fmt.Printf("  Discovered %d enum types\n\n", len(in.enums))
	fmt.Println("  Type 'help' to see available commands")
	fmt.Println("  Type 'list' to see all types")
	fmt.Println(line)
}

func (in *inspector) help() {
	fmt.Print(`
AVAILABLE COMMANDS:
-----------------------------------------------------------
  list / ls              - List all protobuf types
  list --request         - List only *Request types
  list --reply           - List only *Reply types
  <TypeName>             - Inspect a specific type's fields
  search <text>          - Search types by name (alias: find)
  field <name>           - Find types containing a field
  findtype <type>        - Find all fields of a given type
  enum <name>            - Show an enum's values
  whereenum <name>       - Find where an enum is used
  export <TypeName>      - Export a type to JSON
  stats                  - Show statistics
  help / ?               - Show this help
  exit / quit / q        - Exit

EXAMPLES:
  > OpenedOrderInfo      - Show all fields of OpenedOrderInfo
  > search Order         - Find all types with 'Order'
  > field ticket         - Find types with a 'ticket' field
  > findtype double      - Find all double fields
  > enum OrderSendOperationType
  > whereenum OpenedOrderType
  > export OrderSendData
`)
}

func (in *inspector) sortedTypeNames() []string {
	names := make([]string, 0, len(in.msgs))
	for n := range in.msgs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func (in *inspector) list(filter string) {
	names := in.sortedTypeNames()
	title := fmt.Sprintf("ALL TYPES (%d)", len(names))
	match := func(string) bool { return true }
	switch strings.TrimSpace(filter) {
	case "--request":
		title = "REQUEST TYPES"
		match = func(n string) bool { return strings.HasSuffix(n, "Request") }
	case "--reply":
		title = "REPLY TYPES"
		match = func(n string) bool { return strings.HasSuffix(n, "Reply") }
	}
	fmt.Printf("\n%s\n%s\n", title, strings.Repeat("-", 59))
	i := 0
	for _, n := range names {
		if match(n) {
			i++
			fmt.Printf("  %3d. %s\n", i, n)
		}
	}
	if i == 0 {
		fmt.Println("  (none)")
	}
}

func (in *inspector) inspect(name string) {
	md, ok := in.msgs[name]
	if !ok {
		fmt.Printf("Type %q not found. Try 'search %s'.\n", name, name)
		return
	}
	fields := md.Fields()
	fmt.Printf("\nTYPE: %s\n%s\n", md.Name(), strings.Repeat("=", 59))
	fmt.Printf("FIELDS (%d):\n%s\n", fields.Len(), strings.Repeat("-", 59))
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		fmt.Printf("  #%3d %-32s : %s\n", f.Number(), f.Name(), fieldType(f))
	}
}

func (in *inspector) search(q string) {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		fmt.Println("usage: search <text>")
		return
	}
	fmt.Printf("\nTypes matching %q:\n%s\n", q, strings.Repeat("-", 59))
	n := 0
	for _, name := range in.sortedTypeNames() {
		if strings.Contains(strings.ToLower(name), q) {
			n++
			fmt.Printf("  %s\n", name)
		}
	}
	fmt.Printf("  (%d result%s)\n", n, plural(n))
}

func (in *inspector) field(name string) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		fmt.Println("usage: field <name>")
		return
	}
	fmt.Printf("\nTypes containing a field like %q:\n%s\n", name, strings.Repeat("-", 59))
	n := 0
	for _, tn := range in.sortedTypeNames() {
		fields := in.msgs[tn].Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.Get(i)
			if strings.Contains(strings.ToLower(string(f.Name())), name) {
				n++
				fmt.Printf("  %-34s .%s : %s\n", tn, f.Name(), fieldType(f))
				break
			}
		}
	}
	fmt.Printf("  (%d type%s)\n", n, plural(n))
}

func (in *inspector) findType(typ string) {
	typ = strings.TrimSpace(typ)
	if typ == "" {
		fmt.Println("usage: findtype <type>  (e.g. double, string, int32, Timestamp)")
		return
	}
	fmt.Printf("\nFields of type %q:\n%s\n", typ, strings.Repeat("-", 59))
	n := 0
	for _, tn := range in.sortedTypeNames() {
		fields := in.msgs[tn].Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.Get(i)
			if strings.EqualFold(baseType(f), typ) {
				n++
				fmt.Printf("  %-34s .%s : %s\n", tn, f.Name(), fieldType(f))
			}
		}
	}
	fmt.Printf("  (%d field%s)\n", n, plural(n))
}

func (in *inspector) enum(name string) {
	name = strings.TrimSpace(name)
	ed, ok := in.enums[name]
	if !ok {
		// case-insensitive fallback
		for en, e := range in.enums {
			if strings.EqualFold(en, name) {
				ed, ok = e, true
				break
			}
		}
	}
	if !ok {
		fmt.Printf("Enum %q not found.\n", name)
		return
	}
	vals := ed.Values()
	fmt.Printf("\nENUM: %s (%d values)\n%s\n", ed.Name(), vals.Len(), strings.Repeat("-", 59))
	for i := 0; i < vals.Len(); i++ {
		v := vals.Get(i)
		fmt.Printf("  %4d  %s\n", v.Number(), v.Name())
	}
}

func (in *inspector) whereEnum(name string) {
	name = strings.TrimSpace(name)
	if name == "" {
		fmt.Println("usage: whereenum <EnumName>")
		return
	}
	fmt.Printf("\n%q is used by:\n%s\n", name, strings.Repeat("-", 59))
	n := 0
	for _, tn := range in.sortedTypeNames() {
		fields := in.msgs[tn].Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.Get(i)
			if f.Kind() == protoreflect.EnumKind && strings.EqualFold(string(f.Enum().Name()), name) {
				n++
				fmt.Printf("  %-34s .%s\n", tn, f.Name())
			}
		}
	}
	fmt.Printf("  (%d field%s)\n", n, plural(n))
}

func (in *inspector) export(name string) {
	md, ok := in.msgs[strings.TrimSpace(name)]
	if !ok {
		fmt.Printf("Type %q not found.\n", name)
		return
	}
	type fieldJSON struct {
		Number int32  `json:"number"`
		Name   string `json:"name"`
		Type   string `json:"type"`
	}
	out := struct {
		Type   string      `json:"type"`
		Fields []fieldJSON `json:"fields"`
	}{Type: string(md.Name())}
	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		out.Fields = append(out.Fields, fieldJSON{int32(f.Number()), string(f.Name()), fieldType(f)})
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

func (in *inspector) stats() {
	totalFields := 0
	kinds := map[string]int{}
	for _, md := range in.msgs {
		fields := md.Fields()
		totalFields += fields.Len()
		for i := 0; i < fields.Len(); i++ {
			kinds[baseType(fields.Get(i))]++
		}
	}
	reqs, reps := 0, 0
	for n := range in.msgs {
		if strings.HasSuffix(n, "Request") {
			reqs++
		}
		if strings.HasSuffix(n, "Reply") {
			reps++
		}
	}
	fmt.Printf("\nSTATISTICS\n%s\n", strings.Repeat("-", 59))
	fmt.Printf("  Message types : %d  (%d Request, %d Reply)\n", len(in.msgs), reqs, reps)
	fmt.Printf("  Enum types    : %d\n", len(in.enums))
	fmt.Printf("  Total fields  : %d\n\n", totalFields)
	fmt.Println("  Fields by type:")
	type kv struct {
		k string
		v int
	}
	var list []kv
	for k, v := range kinds {
		list = append(list, kv{k, v})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].v > list[j].v })
	for _, e := range list {
		fmt.Printf("    %-20s %d\n", e.k, e.v)
	}
}

// ---- helpers ----

// baseType is the field's element type without repeated/map decoration.
func baseType(f protoreflect.FieldDescriptor) string {
	switch f.Kind() {
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return string(f.Message().Name())
	case protoreflect.EnumKind:
		return string(f.Enum().Name())
	default:
		return f.Kind().String()
	}
}

// fieldType is the display type, including repeated / map decoration.
func fieldType(f protoreflect.FieldDescriptor) string {
	if f.IsMap() {
		return fmt.Sprintf("map<%s, %s>", f.MapKey().Kind().String(), baseType(f.MapValue()))
	}
	if f.IsList() {
		return "repeated " + baseType(f)
	}
	return baseType(f)
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
