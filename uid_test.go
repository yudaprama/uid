package uid

import (
	"bytes"
	"testing"
)

func TestDigits(t *testing.T) {
	if len(digits) != base {
		t.Fatalf("digits length does not match base module")
	}
}

func TestGlobalUIDInit(t *testing.T) {
	if globalUID == nil {
		t.Fatalf("Expected globalUID to be non-nil\n")
	}
	if globalUID.pre == nil || len(globalUID.pre) != preLen {
		t.Fatalf("Expected prefix to be initialized\n")
	}
	if globalUID.seq == 0 {
		t.Fatalf("Expected seq to be non-zero\n")
	}
}

func TestUIDRollover(t *testing.T) {
	globalUID.seq = maxSeq
	// copy
	oldPre := append([]byte{}, globalUID.pre...)
	Next()
	if bytes.Equal(globalUID.pre, oldPre) {
		t.Fatalf("Expected new pre, got the old one\n")
	}
}

func TestGUIDLen(t *testing.T) {
	uid := Next()
	if len(uid) != totalLen {
		t.Fatalf("Expected len of %d, got %d\n", totalLen, len(uid))
	}
}

func TestProperPrefix(t *testing.T) {
	min := byte(255)
	max := byte(0)
	for i := 0; i < len(digits); i++ {
		if digits[i] < min {
			min = digits[i]
		}
		if digits[i] > max {
			max = digits[i]
		}
	}
	total := 100000
	for i := 0; i < total; i++ {
		n := New()
		for j := 0; j < preLen; j++ {
			if n.pre[j] < min || n.pre[j] > max {
				t.Fatalf("Iter %d. Valid range for bytes prefix: [%d..%d]\nIncorrect prefix at pos %d: %v (%s)",
					i, min, max, j, n.pre, string(n.pre))
			}
		}
	}
}

func BenchmarkUIDSpeed(b *testing.B) {
	n := New()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		n.Next()
	}
}

func BenchmarkGlobalUIDSpeed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Next()
	}
}
