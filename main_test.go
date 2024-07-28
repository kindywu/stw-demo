package main

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v3"
)

// go test -bench=BenchmarkHttpHandler -benchmem -benchtime=10s -v
// go test -bench=BenchmarkHttpJson -benchmem -benchtime=10s -v
// go test -bench=BenchmarkHttpArray -benchmem -benchtime=10s -v

// go test -bench=BenchmarkFiberHandler -benchmem -benchtime=10s -v
// go test -bench=BenchmarkFiberJson -benchmem -benchtime=10s -v
// go test -bench=BenchmarkFiberArray -benchmem -benchtime=10s -v

var once sync.Once

func BenchmarkHttpHandler(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		w.Body.Reset()
		httpHandler(w, req)
	}
	printStats(b)
}

func BenchmarkHttpJson(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		w.Body.Reset()
		httpHandlerJson(w, req)
	}
	printStats(b)
}

func BenchmarkHttpArray(b *testing.B) {
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		w.Body.Reset() //注释这句会导致w的内存一直无法释放,从而导致GC次数变少
		httpHandlerArray(w, req)
	}
	printStats(b)
}

func BenchmarkFiberHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)
	app := fiber.New()
	app.Get("/", fiberHandler)
	for n := 0; n < b.N; n++ {
		app.Test(req)
	}

	printStats(b)
}

func BenchmarkFiberJson(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)
	app := fiber.New()
	app.Get("/", fiberHandlerJson)
	for n := 0; n < b.N; n++ {
		app.Test(req)
	}

	printStats(b)
}

func BenchmarkFiberArray(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)
	app := fiber.New()
	app.Get("/", fiberHandlerArray)
	for n := 0; n < b.N; n++ {
		app.Test(req)
	}

	printStats(b)
}

func printStats(b *testing.B) {
	once.Do(func() {
		// 禁用GC，total_pause为零
		// debug.SetGCPercent(-1)
	})

	// 获取基准测试后的GC统计信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	numGC := m.NumGC               //累计GC次数
	pauseTotalNs := m.PauseTotalNs //累计GC SWT时间
	alloc := m.Alloc               //堆内存分配
	totalAlloc := m.TotalAlloc     //累计堆内存分配

	b.Logf("test N: %d, GC: %v, total_pause: %.2fms, alloc: %.2fMB, total_alloc: %.2fMB",
		b.N,
		numGC,
		toMs(pauseTotalNs),
		toMb(alloc),
		toMb(totalAlloc))
}

func toMb(b uint64) float64 {
	return float64(b) / 1024.0 / 1024.0
}

func toMs(b uint64) float64 {
	return float64(b) / 1000.0 / 1000.0
}
