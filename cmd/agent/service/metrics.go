package service

import (
	"math/rand"
	"runtime"

	"github.com/zhupanovdm/go-runtime-monitor/internal/encoder"
)

const (
	Alloc         = "Alloc"
	BuckHashSys   = "BuckHashSys"
	GCCPUFraction = "GCCPUFraction"
	GCSys         = "GCSys"
	HeapAlloc     = "HeapAlloc"
	HeapIdle      = "HeapIdle"
	HeapInuse     = "HeapInuse"
	HeapObjects   = "HeapObjects"
	HeapReleased  = "HeapReleased"
	HeapSys       = "HeapSys"
	LastGC        = "LastGC"
	Lookups       = "Lookups"
	MCacheInuse   = "MCacheInuse"
	MCacheSys     = "MCacheSys"
	MSpanInuse    = "MSpanInuse"
	MSpanSys      = "MSpanSys"
	Mallocs       = "Mallocs"
	NextGC        = "NextGC"
	NumForcedGC   = "NumForcedGC"
	NumGC         = "NumGC"
	OtherSys      = "OtherSys"
	PauseTotalNs  = "PauseTotalNs"
	StackInuse    = "StackInuse"
	StackSys      = "StackSys"
	Sys           = "Sys"
	PollCount     = "PollCount"
	RandomValue   = "RandomValue"
)

func MetricsReader() func(subscriber chan<- encoder.Encoder) {
	var pollCounter int64
	return func(subscriber chan<- encoder.Encoder) {
		read(&pollCounter, subscriber)
	}
}

func read(pollCounter *int64, subscriber chan<- encoder.Encoder) {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	*pollCounter++

	rand.Seed(*pollCounter)

	subscriber <- encoder.NewGaugeI(Alloc, stats.Alloc)
	subscriber <- encoder.NewGaugeI(BuckHashSys, stats.BuckHashSys)
	subscriber <- encoder.NewGaugeF(GCCPUFraction, stats.GCCPUFraction)
	subscriber <- encoder.NewGaugeI(GCSys, stats.GCSys)
	subscriber <- encoder.NewGaugeI(HeapAlloc, stats.HeapAlloc)
	subscriber <- encoder.NewGaugeI(HeapIdle, stats.HeapIdle)
	subscriber <- encoder.NewGaugeI(HeapInuse, stats.HeapInuse)
	subscriber <- encoder.NewGaugeI(HeapObjects, stats.HeapObjects)
	subscriber <- encoder.NewGaugeI(HeapReleased, stats.HeapReleased)
	subscriber <- encoder.NewGaugeI(HeapSys, stats.HeapSys)
	subscriber <- encoder.NewGaugeI(LastGC, stats.LastGC)
	subscriber <- encoder.NewGaugeI(Lookups, stats.Lookups)
	subscriber <- encoder.NewGaugeI(MCacheInuse, stats.MCacheInuse)
	subscriber <- encoder.NewGaugeI(MCacheSys, stats.MCacheSys)
	subscriber <- encoder.NewGaugeI(MSpanInuse, stats.MSpanInuse)
	subscriber <- encoder.NewGaugeI(MSpanSys, stats.MSpanSys)
	subscriber <- encoder.NewGaugeI(Mallocs, stats.Mallocs)
	subscriber <- encoder.NewGaugeI(NextGC, stats.NextGC)
	subscriber <- encoder.NewGaugeI(NumForcedGC, uint64(stats.NumForcedGC))
	subscriber <- encoder.NewGaugeI(NumGC, uint64(stats.NumGC))
	subscriber <- encoder.NewGaugeI(OtherSys, stats.OtherSys)
	subscriber <- encoder.NewGaugeI(PauseTotalNs, stats.PauseTotalNs)
	subscriber <- encoder.NewGaugeI(StackInuse, stats.StackInuse)
	subscriber <- encoder.NewGaugeI(StackSys, stats.StackSys)
	subscriber <- encoder.NewGaugeI(Sys, stats.Sys)
	subscriber <- encoder.NewCounter(PollCount, *pollCounter)
	subscriber <- encoder.NewGaugeF(RandomValue, rand.Float64())
}
