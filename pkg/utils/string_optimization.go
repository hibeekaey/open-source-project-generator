package utils

import (
	"strings"
	"sync"
)

// StringPool provides string pooling to reduce memory allocations
type StringPool struct {
	pool sync.Pool
}

// NewStringPool creates a new string pool
func NewStringPool() *StringPool {
	return &StringPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]string, 0, 16) // Pre-allocate capacity for 16 strings
			},
		},
	}
}

// Get retrieves a string slice from the pool
func (sp *StringPool) Get() []string {
	return sp.pool.Get().([]string)
}

// Put returns a string slice to the pool after clearing it
func (sp *StringPool) Put(s []string) {
	s = s[:0] // Clear the slice but keep capacity
	sp.pool.Put(s)
}

// StringBuilder provides efficient string building with pooled buffers
type StringBuilder struct {
	buf strings.Builder
}

// NewStringBuilder creates a new string builder
func NewStringBuilder() *StringBuilder {
	return &StringBuilder{}
}

// WriteString writes a string to the builder
func (sb *StringBuilder) WriteString(s string) {
	sb.buf.WriteString(s)
}

// WriteByte writes a byte to the builder
func (sb *StringBuilder) WriteByte(b byte) {
	sb.buf.WriteByte(b)
}

// String returns the built string
func (sb *StringBuilder) String() string {
	return sb.buf.String()
}

// Reset resets the builder for reuse
func (sb *StringBuilder) Reset() {
	sb.buf.Reset()
}

// Len returns the current length
func (sb *StringBuilder) Len() int {
	return sb.buf.Len()
}

// Cap returns the current capacity
func (sb *StringBuilder) Cap() int {
	return sb.buf.Cap()
}

// StringBuilderPool provides pooled string builders
type StringBuilderPool struct {
	pool sync.Pool
}

// NewStringBuilderPool creates a new string builder pool
func NewStringBuilderPool() *StringBuilderPool {
	return &StringBuilderPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &StringBuilder{}
			},
		},
	}
}

// Get retrieves a string builder from the pool
func (sbp *StringBuilderPool) Get() *StringBuilder {
	sb := sbp.pool.Get().(*StringBuilder)
	sb.Reset()
	return sb
}

// Put returns a string builder to the pool
func (sbp *StringBuilderPool) Put(sb *StringBuilder) {
	sb.Reset()
	sbp.pool.Put(sb)
}

// OptimizedStringOperations provides optimized string operations
type OptimizedStringOperations struct {
	builderPool *StringBuilderPool
	stringPool  *StringPool
}

// NewOptimizedStringOperations creates a new optimized string operations instance
func NewOptimizedStringOperations() *OptimizedStringOperations {
	return &OptimizedStringOperations{
		builderPool: NewStringBuilderPool(),
		stringPool:  NewStringPool(),
	}
}

// JoinStrings efficiently joins strings with a separator
func (oso *OptimizedStringOperations) JoinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	sb := oso.builderPool.Get()
	defer oso.builderPool.Put(sb)

	sb.WriteString(strs[0])
	for i := 1; i < len(strs); i++ {
		sb.WriteString(sep)
		sb.WriteString(strs[i])
	}

	return sb.String()
}

// ConcatenateStrings efficiently concatenates multiple strings
func (oso *OptimizedStringOperations) ConcatenateStrings(strs ...string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}

	sb := oso.builderPool.Get()
	defer oso.builderPool.Put(sb)

	for _, s := range strs {
		sb.WriteString(s)
	}

	return sb.String()
}

// ReplaceMultiple efficiently performs multiple string replacements
func (oso *OptimizedStringOperations) ReplaceMultiple(s string, replacements map[string]string) string {
	if len(replacements) == 0 {
		return s
	}

	result := s
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	return result
}

// SplitAndTrim splits a string and trims whitespace from each part
func (oso *OptimizedStringOperations) SplitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := oso.stringPool.Get()
	defer oso.stringPool.Put(result)

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	// Return a copy since we're returning the slice to the pool
	finalResult := make([]string, len(result))
	copy(finalResult, result)
	return finalResult
}

// ContainsAny efficiently checks if a string contains any of the given substrings
func (oso *OptimizedStringOperations) ContainsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// HasPrefixAny efficiently checks if a string has any of the given prefixes
func (oso *OptimizedStringOperations) HasPrefixAny(s string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// HasSuffixAny efficiently checks if a string has any of the given suffixes
func (oso *OptimizedStringOperations) HasSuffixAny(s string, suffixes []string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(s, suffix) {
			return true
		}
	}
	return false
}

// Global instances for convenience
var (
	globalStringOps   = NewOptimizedStringOperations()
	globalStringPool  = NewStringPool()
	globalBuilderPool = NewStringBuilderPool()
)

// Convenience functions using global instances

// JoinStrings efficiently joins strings using the global instance
func JoinStrings(strs []string, sep string) string {
	return globalStringOps.JoinStrings(strs, sep)
}

// ConcatenateStrings efficiently concatenates strings using the global instance
func ConcatenateStrings(strs ...string) string {
	return globalStringOps.ConcatenateStrings(strs...)
}

// ReplaceMultiple performs multiple replacements using the global instance
func ReplaceMultiple(s string, replacements map[string]string) string {
	return globalStringOps.ReplaceMultiple(s, replacements)
}

// SplitAndTrim splits and trims using the global instance
func SplitAndTrim(s, sep string) []string {
	return globalStringOps.SplitAndTrim(s, sep)
}

// GetStringBuilder gets a string builder from the global pool
func GetStringBuilder() *StringBuilder {
	return globalBuilderPool.Get()
}

// PutStringBuilder returns a string builder to the global pool
func PutStringBuilder(sb *StringBuilder) {
	globalBuilderPool.Put(sb)
}

// GetStringSlice gets a string slice from the global pool
func GetStringSlice() []string {
	return globalStringPool.Get()
}

// PutStringSlice returns a string slice to the global pool
func PutStringSlice(s []string) {
	globalStringPool.Put(s)
}
