/*
Copyright (c) 2024, Shanghai Iluvatar CoreX Semiconductor Co., Ltd.
All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ixml

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"unsafe"
)

const (
	debugMaxDepth      = 2
	debugMaxArrayElems = 4
	debugMaxStringLen  = 240
)

var debugLogEnabled atomic.Bool

func init() {
	level := strings.ToLower(strings.TrimSpace(os.Getenv("GO_IXML_LOG_LEVEL")))
	if level == "debug" {
		debugLogEnabled.Store(true)
	}
}

func SetDebugLogEnabled(enabled bool) {
	debugLogEnabled.Store(enabled)
}

func IsDebugLogEnabled() bool {
	return debugLogEnabled.Load()
}

// debugLogCgoInput logs Go wrapper args before the C call.
// Pointer args are not dereferenced (may be uninitialized out-params).
func debugLogCgoInput(name string, kv ...any) {
	if !debugLogEnabled.Load() {
		return
	}
	log.Printf("[go-ixml][debug] %s in %s", name, formatDebugKV(kv, false))
}

// debugLogCgoOutput logs out-params (pointer args, dereferenced) and Return after the C call.
func debugLogCgoOutput(name string, ret Return, kv ...any) {
	if !debugLogEnabled.Load() {
		return
	}
	if len(kv) == 0 {
		log.Printf("[go-ixml][debug] %s out -> %s(%d)", name, ret.String(), ret)
		return
	}
	log.Printf("[go-ixml][debug] %s out %s -> %s(%d)", name, formatDebugKV(kv, true), ret.String(), ret)
}

func formatDebugKV(kv []any, deref bool) string {
	if len(kv) == 0 {
		return ""
	}
	parts := make([]string, 0, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		key, _ := kv[i].(string)
		parts = append(parts, fmt.Sprintf("%s=%s", key, formatDebugArg(key, kv[i+1], deref)))
	}
	return strings.Join(parts, " ")
}

func formatDebugArg(key string, v any, deref bool) string {
	if v == nil {
		return "<nil>"
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer, reflect.UnsafePointer:
		if rv.IsNil() {
			return "<nil>"
		}
		if !deref {
			return fmt.Sprintf("%p", v)
		}
		elem := rv.Elem()
		switch elem.Kind() {
		case reflect.Uint8, reflect.Int8:
			return formatBytePtrString(unsafe.Pointer(rv.Pointer()), 256)
		}
		if s, ok := formatNamedDebugValue(key, elem); ok {
			return s
		}
		return truncateDebugString(formatCompactValue(elem, 0))
	default:
		if s, ok := formatNamedDebugValue(key, rv); ok {
			return s
		}
		return truncateDebugString(formatCompactValue(rv, 0))
	}
}

// formatNamedDebugValue applies known semantic encodings (e.g. CUDA version 10020 -> 10.2).
func formatNamedDebugValue(key string, rv reflect.Value) (string, bool) {
	switch key {
	case "CudaDriverVersion":
		var ver int64
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ver = rv.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ver = int64(rv.Uint())
		default:
			return "", false
		}
		major := ver / 1000
		minor := (ver % 1000) / 10
		return fmt.Sprintf("%d(%d.%d)", ver, major, minor), true
	}
	return "", false
}

func formatCompactValue(rv reflect.Value, depth int) string {
	if !rv.IsValid() {
		return "<invalid>"
	}
	switch rv.Kind() {
	case reflect.Pointer, reflect.UnsafePointer:
		if rv.IsNil() {
			return "<nil>"
		}
		if depth >= debugMaxDepth {
			return fmt.Sprintf("%p", rv.Interface())
		}
		elem := rv.Elem()
		if elem.Kind() == reflect.Struct && !hasExportableDebugFields(elem) {
			return fmt.Sprintf("%p", rv.Interface())
		}
		return formatCompactValue(elem, depth+1)
	case reflect.Struct:
		parts := make([]string, 0, rv.NumField())
		for i := 0; i < rv.NumField(); i++ {
			f := rv.Field(i)
			if !f.CanInterface() {
				continue
			}
			name := rv.Type().Field(i).Name
			if depth >= debugMaxDepth && isDebugComplexKind(f.Kind()) {
				continue
			}
			if isDebugZeroValue(f) && isDebugComplexKind(f.Kind()) {
				continue
			}
			parts = append(parts, name+"="+formatCompactValue(f, depth+1))
		}
		return "{" + strings.Join(parts, " ") + "}"
	case reflect.Array, reflect.Slice:
		n := rv.Len()
		if n == 0 {
			return "[]"
		}
		if s, ok := formatByteArrayAsString(rv); ok {
			return s
		}
		shown := 0
		parts := make([]string, 0, debugMaxArrayElems)
		for i := 0; i < n && shown < debugMaxArrayElems; i++ {
			elem := rv.Index(i)
			if isDebugZeroValue(elem) {
				continue
			}
			parts = append(parts, formatCompactValue(elem, depth+1))
			shown++
		}
		if shown == 0 {
			return fmt.Sprintf("[%d]{}", n)
		}
		if shown < n {
			return fmt.Sprintf("[%d]{%s ...}", n, strings.Join(parts, ", "))
		}
		return fmt.Sprintf("[%d]{%s}", n, strings.Join(parts, ", "))
	case reflect.String:
		s := rv.String()
		if len(s) > 64 {
			return fmt.Sprintf("%q...", s[:64])
		}
		return fmt.Sprintf("%q", s)
	default:
		if rv.CanInterface() {
			return fmt.Sprintf("%v", rv.Interface())
		}
		return rv.Kind().String()
	}
}

// formatByteArrayAsString treats [N]byte / [N]int8 C-string buffers as quoted text.
func formatByteArrayAsString(rv reflect.Value) (string, bool) {
	if rv.Len() == 0 {
		return "", false
	}
	elemKind := rv.Type().Elem().Kind()
	switch elemKind {
	case reflect.Uint8, reflect.Int8:
	default:
		return "", false
	}
	b := make([]byte, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		var c byte
		if elemKind == reflect.Uint8 {
			c = byte(rv.Index(i).Uint())
		} else {
			c = byte(rv.Index(i).Int())
		}
		if c == 0 {
			break
		}
		if c < 0x20 || c > 0x7e {
			return "", false
		}
		b = append(b, c)
	}
	return fmt.Sprintf("%q", string(b)), true
}

func isDebugComplexKind(k reflect.Kind) bool {
	switch k {
	case reflect.Struct, reflect.Array, reflect.Slice, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface:
		return true
	default:
		return false
	}
}

func hasExportableDebugFields(rv reflect.Value) bool {
	if rv.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).CanInterface() {
			return true
		}
	}
	return false
}

func isDebugZeroValue(rv reflect.Value) bool {
	if !rv.IsValid() {
		return true
	}
	switch rv.Kind() {
	case reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return rv.IsNil()
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			f := rv.Field(i)
			if !f.CanInterface() {
				continue
			}
			if !isDebugZeroValue(f) {
				return false
			}
		}
		return true
	case reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			if !isDebugZeroValue(rv.Index(i)) {
				return false
			}
		}
		return true
	default:
		z := reflect.Zero(rv.Type())
		return reflect.DeepEqual(rv.Interface(), z.Interface())
	}
}

func truncateDebugString(s string) string {
	if len(s) <= debugMaxStringLen {
		return s
	}
	return s[:debugMaxStringLen] + "..."
}

func formatBytePtrString(p unsafe.Pointer, max int) string {
	if p == nil || max <= 0 {
		return ""
	}
	b := make([]byte, 0, 64)
	for i := 0; i < max; i++ {
		c := *(*byte)(unsafe.Pointer(uintptr(p) + uintptr(i)))
		if c == 0 {
			break
		}
		b = append(b, c)
	}
	return string(b)
}
