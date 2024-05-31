package main

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"testing"
)

func TestFetchWasmExecJS(t *testing.T) {
	if err := os.RemoveAll(tempDir); err != nil {
		t.Fatal(err)
	}
	v, err := goVersion(".")
	if err != nil {
		t.Fatalf("get go version: %v", err)
	}

	cache, err := getCachedWasmExecJS(v)
	if err == nil {
		t.Fatal("cache should be missing")
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("cache file has an error: %v", err)
	}

	src, err := fetchWasmExecJS(v)
	if err != nil {
		t.Fatalf("fetch from source: %v", err)
	}
	if len(src) == 0 {
		t.Fatal("cannot fetch from source")
	}
	t.Log(string(src[:100]))

	cache, err = getCachedWasmExecJS(v)
	if err != nil {
		t.Fatalf("cannot get cache: %v", err)
	}
	if !bytes.Equal(src, cache) {
		t.Fatal("source != cache")
	}
	t.Log(string(cache[:100]))
}
