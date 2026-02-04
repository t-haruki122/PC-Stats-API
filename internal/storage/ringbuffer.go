package storage

import (
	"sync"
	"time"

	"pc-stats-api/internal/collector"
)

// RingBuffer is a thread-safe fixed-size circular buffer for metric samples
type RingBuffer struct {
	mu       sync.RWMutex
	samples  []*collector.MetricSample
	capacity int
	head     int
	size     int
}

// NewRingBuffer creates a new ring buffer with the specified capacity
func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		samples:  make([]*collector.MetricSample, capacity),
		capacity: capacity,
		head:     0,
		size:     0,
	}
}

// Add adds a new sample to the buffer
func (rb *RingBuffer) Add(sample *collector.MetricSample) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	rb.samples[rb.head] = sample
	rb.head = (rb.head + 1) % rb.capacity

	if rb.size < rb.capacity {
		rb.size++
	}
}

// GetLatest returns the most recent sample
func (rb *RingBuffer) GetLatest() *collector.MetricSample {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.size == 0 {
		return nil
	}

	// Most recent is at (head - 1)
	index := (rb.head - 1 + rb.capacity) % rb.capacity
	return rb.samples[index]
}

// GetHistory returns samples from the last N seconds
func (rb *RingBuffer) GetHistory(seconds int) []*collector.MetricSample {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.size == 0 {
		return nil
	}

	cutoff := time.Now().Add(-time.Duration(seconds) * time.Second)
	result := make([]*collector.MetricSample, 0, rb.size)

	// Iterate from oldest to newest
	start := (rb.head - rb.size + rb.capacity) % rb.capacity
	for i := 0; i < rb.size; i++ {
		index := (start + i) % rb.capacity
		sample := rb.samples[index]

		if sample != nil && sample.Timestamp.After(cutoff) {
			result = append(result, sample)
		}
	}

	return result
}

// GetAll returns all samples in chronological order
func (rb *RingBuffer) GetAll() []*collector.MetricSample {
	rb.mu.RLock()
	defer rb.mu.RUnlock()

	if rb.size == 0 {
		return nil
	}

	result := make([]*collector.MetricSample, 0, rb.size)
	start := (rb.head - rb.size + rb.capacity) % rb.capacity

	for i := 0; i < rb.size; i++ {
		index := (start + i) % rb.capacity
		if rb.samples[index] != nil {
			result = append(result, rb.samples[index])
		}
	}

	return result
}

// Size returns the current number of samples in the buffer
func (rb *RingBuffer) Size() int {
	rb.mu.RLock()
	defer rb.mu.RUnlock()
	return rb.size
}
