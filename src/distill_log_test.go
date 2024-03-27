package main

import (
	"log"
	"testing"
	"time"
)

func TestDistillReader(t *testing.T) {
	dlogReader, err := NewDistillReader("../output/dlog.log")
	if err != nil {
		t.Error(err)
		return
	}
	start := time.Now()
	line := dlogReader.Search("115", 0, 10)
	cost_ms := time.Since(start).Milliseconds()
	log.Printf("search over, cost_ms:%d", cost_ms)
	if len(line) == 0 {
		t.Errorf("search fail")
	}
}
