//
//  stats_ws.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

import (
	"runtime"
	"time"
	"bytes"
	"net/http"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"github.com/ortuman/mercury/logger"
	"github.com/ortuman/mercury/config"
)

type goStats struct {
	Time int64 `json:"time"`

	// runtime
	GoVersion    string `json:"go_version"`
	GoOs         string `json:"go_os"`
	GoArch       string `json:"go_arch"`
	CpuNum       int    `json:"cpu_num"`
	GoroutineNum int    `json:"goroutine_num"`
	Gomaxprocs   int    `json:"gomaxprocs"`
	CgoCallNum   int64  `json:"cgo_call_num"`

	// memory
	MemoryAlloc      uint64 `json:"memory_alloc"`
	MemoryTotalAlloc uint64 `json:"memory_total_alloc"`
	MemorySys        uint64 `json:"memory_sys"`
	MemoryLookups    uint64 `json:"memory_lookups"`
	MemoryMallocs    uint64 `json:"memory_mallocs"`
	MemoryFrees      uint64 `json:"memory_frees"`

	// Stack
	StackInUse uint64 `json:"memory_stack"`

	// heap
	HeapAlloc    uint64 `json:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys"`
	HeapIdle     uint64 `json:"heap_idle"`
	HeapInuse    uint64 `json:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released"`
	HeapObjects  uint64 `json:"heap_objects"`

	// garbage collection
	GcNext           uint64    `json:"gc_next"`
	GcLast           uint64    `json:"gc_last"`
	GcNum            uint32    `json:"gc_num"`
	GcPerSecond      float64   `json:"gc_per_second"`
	GcPausePerSecond float64   `json:"gc_pause_per_second"`
	GcPause          []float64 `json:"gc_pause"`
}

var lastSampleTime time.Time
var lastPauseNs uint64 = 0
var lastNumGc uint32 = 0

func NewStatsWS() *restful.WebService {
	s := new(restful.WebService).Path("/v1/stats")

	s.Route(s.GET("/push").To(getPushStats))
	s.Route(s.GET("/sys").To(getSysStats))

	return s
}

func getPushStats(_ *restful.Request, response *restful.Response) {
	logger.Infof("stats_ws: retrieving push stats...")
	writeStats(globalSender.Stats(), response)
}

func getSysStats(_ *restful.Request, response *restful.Response) {
	logger.Infof("stats_ws: retrieving go stats...")

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	now := time.Now()

	var gcPausePerSecond float64

	if lastPauseNs > 0 {
		pauseSinceLastSample := mem.PauseTotalNs - lastPauseNs
		gcPausePerSecond = float64(pauseSinceLastSample) / float64(time.Millisecond)
	}

	lastPauseNs = mem.PauseTotalNs

	countGc := int(mem.NumGC - lastNumGc)

	var gcPerSecond float64

	if lastNumGc > 0 {
		diff := float64(countGc)
		diffTime := now.Sub(lastSampleTime).Seconds()
		gcPerSecond = diff / diffTime
	}

	gcPause := make([]float64, countGc)

	if countGc > 0 {
		if countGc > 256 {
			// lagging GC pause times
			countGc = 256
		}

		for i := 0; i < countGc; i++ {
			idx := int((mem.NumGC-uint32(i))+255) % 256
			pause := float64(mem.PauseNs[idx])
			gcPause[i] = pause / float64(time.Millisecond)
		}
	}

	lastNumGc = mem.NumGC
	lastSampleTime = time.Now()

	stats := &goStats{
		Time:         now.UnixNano(),
		GoVersion:    runtime.Version(),
		GoOs:         runtime.GOOS,
		GoArch:       runtime.GOARCH,
		CpuNum:       runtime.NumCPU(),
		GoroutineNum: runtime.NumGoroutine(),
		Gomaxprocs:   runtime.GOMAXPROCS(0),
		CgoCallNum:   runtime.NumCgoCall(),

		// memory
		MemoryAlloc:      mem.Alloc,
		MemoryTotalAlloc: mem.TotalAlloc,
		MemorySys:        mem.Sys,
		MemoryLookups:    mem.Lookups,
		MemoryMallocs:    mem.Mallocs,
		MemoryFrees:      mem.Frees,

		// Stack
		StackInUse: mem.StackInuse,

		// heap
		HeapAlloc:    mem.HeapAlloc,
		HeapSys:      mem.HeapSys,
		HeapIdle:     mem.HeapIdle,
		HeapInuse:    mem.HeapInuse,
		HeapReleased: mem.HeapReleased,
		HeapObjects:  mem.HeapObjects,

		// garbage collection
		GcNext:           mem.NextGC,
		GcLast:           mem.LastGC,
		GcNum:            mem.NumGC,
		GcPerSecond:      gcPerSecond,
		GcPausePerSecond: gcPausePerSecond,
		GcPause:          gcPause,
	}

	writeStats(stats, response)
}

func writeStats(stats interface{}, response *restful.Response) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(stats); err != nil {
		logger.Errorf("stats_ws: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Server", config.ServiceName + "/" + config.ServiceVersion + " (" + runtime.GOOS + ")")
	response.Header().Set("Content-Type", "application/json")
	response.Write(buf.Bytes())
}
