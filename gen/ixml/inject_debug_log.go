// Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
// All Rights Reserved.
//
// Post-process c-for-go output: inject input/output debug logging around
// every C. call in generated ixml.go wrappers. Run via:
//
//	go run gen/ixml/inject_debug_log.go -- pkg/ixml/ixml.go
package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	funcHeaderRe = regexp.MustCompile(`(?m)^func (\w+)\(([^)]*)\) Return \{`)
	cCallLineRe  = regexp.MustCompile(`(?m)^(\t__ret := C\.[^\n]+\n)`)
	retAssignRe  = regexp.MustCompile(`(?m)^(\t__v := \(Return\)\(__ret\))\n(\treturn __v\n)`)
	oldDebugRe   = regexp.MustCompile(`(?m)^\tdebugLogCgoCall\([^\n]*\n`)
	inputDebugRe = regexp.MustCompile(`debugLogCgoInput\(`)
	// Strip previously injected input/output blocks (for upgrade re-run after partial edits).
	inputBlockRe = regexp.MustCompile(`(?m)^\tif IsDebugLogEnabled\(\) \{\n\t\tdebugLogCgoInput\([^\n]*\n\t\}\n`)
	outputBlockRe = regexp.MustCompile(`(?m)^\tif IsDebugLogEnabled\(\) \{\n\t\tdebugLogCgoOutput\([^\n]*\n\t\}\n`)
)

func main() {
	args := make([]string, 0, len(os.Args)-1)
	for _, a := range os.Args[1:] {
		if a == "--" {
			continue
		}
		args = append(args, a)
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [--] <path/to/ixml.go>\n", os.Args[0])
		os.Exit(1)
	}
	path := args[0]
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
		os.Exit(1)
	}

	out, injected, skipped, err := injectDebugLog(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "inject: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(path, out, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", path, err)
		os.Exit(1)
	}
	fmt.Printf("inject_debug_log: %s injected=%d skipped=%d\n", path, injected, skipped)
}

func injectDebugLog(src []byte) (out []byte, injected, skipped int, err error) {
	locs := funcHeaderRe.FindAllSubmatchIndex(src, -1)
	if len(locs) == 0 {
		return src, 0, 0, fmt.Errorf("no Return wrappers found")
	}

	var buf bytes.Buffer
	prev := 0
	for i, loc := range locs {
		name := string(src[loc[2]:loc[3]])
		paramsRaw := string(src[loc[4]:loc[5]])
		allParams, ptrParams := parseParams(paramsRaw)

		headerEnd := loc[1]
		bodyStart := headerEnd
		bodyEnd := len(src)
		if i+1 < len(locs) {
			bodyEnd = locs[i+1][0]
		}

		buf.Write(src[prev:bodyStart])
		body := src[bodyStart:bodyEnd]

		if !bytes.Contains(body, []byte("__ret := C.")) {
			buf.Write(body)
			prev = bodyEnd
			continue
		}

		// Idempotent: already has new-style input logging.
		if inputDebugRe.Match(body) && !oldDebugRe.Match(body) {
			skipped++
			buf.Write(body)
			prev = bodyEnd
			continue
		}

		// Upgrade path: strip old/new debug so we can re-inject cleanly.
		body = oldDebugRe.ReplaceAll(body, nil)
		body = inputBlockRe.ReplaceAll(body, nil)
		body = outputBlockRe.ReplaceAll(body, nil)

		cLoc := cCallLineRe.FindSubmatchIndex(body)
		if cLoc == nil {
			return nil, injected, skipped, fmt.Errorf("function %s: missing __ret := C. line", name)
		}
		rLoc := retAssignRe.FindSubmatchIndex(body)
		if rLoc == nil {
			return nil, injected, skipped, fmt.Errorf("function %s: missing __v/return pattern", name)
		}

		var b bytes.Buffer
		b.Write(body[:cLoc[2]])
		writeInputBlock(&b, name, allParams)
		b.Write(body[cLoc[2]:rLoc[3]]) // through __v line
		b.WriteByte('\n')
		writeOutputBlock(&b, name, ptrParams)
		b.Write(body[rLoc[4]:]) // return __v ...

		buf.Write(b.Bytes())
		injected++
		prev = bodyEnd
	}
	buf.Write(src[prev:])
	return buf.Bytes(), injected, skipped, nil
}

func parseParams(raw string) (all []string, ptrs []string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		fields := strings.Fields(part)
		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		typ := fields[len(fields)-1]
		all = append(all, name)
		if strings.HasPrefix(typ, "*") {
			ptrs = append(ptrs, name)
		}
	}
	return all, ptrs
}

func writeInputBlock(b *bytes.Buffer, name string, params []string) {
	b.WriteString("\tif IsDebugLogEnabled() {\n")
	b.WriteString("\t\tdebugLogCgoInput(")
	b.WriteString(fmt.Sprintf("%q", name))
	for _, p := range params {
		b.WriteString(fmt.Sprintf(", %q, %s", p, p))
	}
	b.WriteString(")\n")
	b.WriteString("\t}\n")
}

func writeOutputBlock(b *bytes.Buffer, name string, ptrParams []string) {
	b.WriteString("\tif IsDebugLogEnabled() {\n")
	b.WriteString("\t\tdebugLogCgoOutput(")
	b.WriteString(fmt.Sprintf("%q, __v", name))
	for _, p := range ptrParams {
		b.WriteString(fmt.Sprintf(", %q, %s", p, p))
	}
	b.WriteString(")\n")
	b.WriteString("\t}\n")
}
