package searchapp

import (
	"testing"
)

// TestLeftCutPadString
func TestLeftCutPadString(t *testing.T) {
	if leftCutPadString("abc", -1) != "" {
		t.Fatal("Incorrect left cut from abc to '' (negative)")
	}
	if leftCutPadString("abc", 0) != "" {
		t.Fatal("Incorrect left cut from abc to ''")
	}
	if leftCutPadString("abc", 1) != "…" {
		t.Fatal("Incorrect left cut from abc to …")
	}
	if leftCutPadString("abc", 2) != "…c" {
		t.Fatal("Incorrect left cut from abc to …c")
	}
	if leftCutPadString("abc", 3) != "abc" {
		t.Fatal("Incorrect left cut from abc to abc")
	}
	if leftCutPadString("abc", 5) != "  abc" {
		t.Fatal("Incorrect left pad from abc to '  abc'")
	}

	// unicode
	if leftCutPadString("♥♥♥♥", -1) != "" {
		t.Fatal("Incorrect left cut from ♥♥♥♥ to '' (negative)")
	}
	if leftCutPadString("♥♥♥♥", 0) != "" {
		t.Fatal("Incorrect left cut from ♥♥♥♥ to ''")
	}
	if leftCutPadString("♥♥♥♥", 1) != "…" {
		t.Fatal("Incorrect left cut from ♥♥♥♥ to …")
	}
	if leftCutPadString("♥♥♥♥", 2) != "…♥" {
		t.Fatal("Incorrect left cut from ♥♥♥♥ to …♥")
	}
	if leftCutPadString("♥♥♥♥", 4) != "♥♥♥♥" {
		t.Fatal("Incorrect left cut from ♥♥♥♥ to ♥♥♥♥")
	}
	if leftCutPadString("♥♥♥♥", 6) != "  ♥♥♥♥" {
		t.Fatal("Incorrect left pad from ♥♥♥♥ to '  ♥♥♥♥'")
	}
}

// TestRightCutPadString
func TestRightCutPadString(t *testing.T) {
	if rightCutPadString("abc", -1) != "" {
		t.Fatal("Incorrect right cut from abc to '' (negative)")
	}
	if rightCutPadString("abc", 0) != "" {
		t.Fatal("Incorrect right cut from abc to ''")
	}
	if rightCutPadString("abc", 1) != "…" {
		t.Fatal("Incorrect right cut from abc to …")
	}
	if rightCutPadString("abc", 2) != "a…" {
		t.Fatal("Incorrect right cut from abc to a…")
	}
	if rightCutPadString("abc", 3) != "abc" {
		t.Fatal("Incorrect right cut from abc to abc")
	}
	if rightCutPadString("abc", 5) != "abc  " {
		t.Fatal("Incorrect right pad from abc to 'abc  '")
	}

	// unicode
	if rightCutPadString("♥♥♥♥", -1) != "" {
		t.Fatal("Incorrect right cut from ♥♥♥♥ to '' (negative)")
	}
	if rightCutPadString("♥♥♥♥", 0) != "" {
		t.Fatal("Incorrect right cut from ♥♥♥♥ to ''")
	}
	if rightCutPadString("♥♥♥♥", 1) != "…" {
		t.Fatal("Incorrect right cut from ♥♥♥♥ to …")
	}
	if rightCutPadString("♥♥♥♥", 2) != "♥…" {
		t.Fatal("Incorrect right cut from ♥♥♥♥ to ♥…")
	}
	if rightCutPadString("♥♥♥♥", 4) != "♥♥♥♥" {
		t.Fatal("Incorrect right cut from ♥♥♥♥ to ♥♥♥♥")
	}
	if rightCutPadString("♥♥♥♥", 6) != "♥♥♥♥  " {
		t.Fatal("Incorrect right pad from ♥♥♥♥ to '♥♥♥♥  '")
	}
}
