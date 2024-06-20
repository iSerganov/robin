package robin

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNext(t *testing.T) {
	asrt := assert.New(t)
	tests := []struct {
		name     string
		items    []WItem[string]
		load     int32
		expected []string
	}{
		{
			name: "4 items 24 iterations total weight 14",
			items: []WItem[string]{
				{item: "A", weight: 8},
				{item: "B", weight: 3},
				{item: "C", weight: 2},
				{item: "D", weight: 1},
			},
			load: 14,
			expected: []string{"A", "A", "A", "A", "A", "A", "A", "A", "B", "B", "B", "C", "C", "D",
				"A", "A", "A", "A", "A", "A", "A", "A", "B", "B"},
		},
		{
			name: "3 items 15 iterations total weight 12",
			items: []WItem[string]{
				{item: "NICE", weight: 6},
				{item: "AWESOME", weight: 5},
				{item: "THRILLING", weight: 1},
			},
			load: 12,
			expected: []string{"NICE", "NICE", "NICE", "NICE", "NICE", "NICE", "AWESOME",
				"AWESOME", "AWESOME", "AWESOME", "AWESOME", "THRILLING", "NICE", "NICE", "NICE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(it *testing.T) {
			wrr := &WRR[string]{
				items: tt.items,
			}
			wrr.cl.Store(tt.load)
			for i, e := range tt.expected {
				asrt.Equal(e, wrr.Next(), "The two words should be the same. Iteration %d", i)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		initial  []WItem[string]
		item     string
		weight   int
		expected []WItem[string]
	}{
		{
			name:     "Add item with higher weight",
			initial:  []WItem[string]{{item: "A", weight: 2}},
			item:     "B",
			weight:   3,
			expected: []WItem[string]{{item: "B", weight: 3}, {item: "A", weight: 2}},
		},
		{
			name:     "Add item with lower weight",
			initial:  []WItem[string]{{item: "A", weight: 3}},
			item:     "B",
			weight:   2,
			expected: []WItem[string]{{item: "A", weight: 3}, {item: "B", weight: 2}},
		},
		{
			name:     "Add item with equal weight",
			initial:  []WItem[string]{{item: "A", weight: 2}},
			item:     "B",
			weight:   2,
			expected: []WItem[string]{{item: "A", weight: 2}, {item: "B", weight: 2}},
		},
		{
			name:     "Add multiple items with different weights",
			initial:  []WItem[string]{{item: "A", weight: 2}},
			item:     "B",
			weight:   3,
			expected: []WItem[string]{{item: "B", weight: 3}, {item: "A", weight: 2}},
		},
		{
			name:     "Add item to an empty list",
			initial:  []WItem[string]{},
			item:     "A",
			weight:   2,
			expected: []WItem[string]{{item: "A", weight: 2}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrr := &WRR[string]{
				items: tt.initial,
			}
			wrr.Add(tt.item, tt.weight)

			if !reflect.DeepEqual(wrr.items, tt.expected) {
				t.Errorf("Add() = %v, want %v", wrr.items, tt.expected)
			}
		})
	}
}

func TestReset(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Reset multiple times"},
		{name: "Reset after adding items"},
		{name: "Reset after Next() calls"},
		{name: "Reset after multiple Next() calls"},
		{name: "Reset after Reset() call"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrr := &WRR[string]{
				items: []WItem[string]{
					{item: "A", weight: 2},
					{item: "B", weight: 3},
					{item: "C", weight: 1},
				},
			}
			wrr.cl.Store(6)

			// Perform some operations before resetting
			wrr.Next()
			wrr.Add("D", 4)

			// Call Reset multiple times
			for i := 0; i < 5; i++ {
				wrr.Reset()
			}

			// Check that Reset() does not panic
			assert.NotPanics(t, func() { wrr.Reset() })

			// Verify that the WRR is reset correctly
			assert.Equal(t, 0, int(wrr.cl.Load()))
			assert.Equal(t, 0, int(wrr.cr.Load()))
			assert.Empty(t, wrr.items)
		})
	}
}
