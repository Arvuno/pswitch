package logx

import "testing"

func TestResolveColorEnabledDefaultsToCapability(t *testing.T) {
	got := resolveColorEnabled(nil, true)
	if !got {
		t.Fatal("expected color to be enabled when terminal supports it")
	}

	got = resolveColorEnabled(nil, false)
	if got {
		t.Fatal("expected color to be disabled when terminal does not support it")
	}
}

func TestResolveColorEnabledHonorsOverride(t *testing.T) {
	on := true
	off := false

	if !resolveColorEnabled(&on, false) {
		t.Fatal("expected explicit true override to enable color")
	}
	if resolveColorEnabled(&off, true) {
		t.Fatal("expected explicit false override to disable color")
	}
}
