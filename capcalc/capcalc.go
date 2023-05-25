// Package capcalc provides helper to calculate channels capacity with specific strategy.
package capcalc

import "math"

// Sum calcs total capacity of all channels.
func Sum[T any](channels ...chan T) int {
	if len(channels) == 0 {
		return 0
	}
	sum := 0
	for i := range channels {
		sum += cap(channels[i])
	}
	return sum
}

// Mult calcs total capacity of all channels and multiplies it by N.
func Mult[T any](n int, channels ...chan T) int {
	return Sum(channels...) * n
}

// Min calcs minimal capacity of all channels.
func Min[T any](channels ...chan T) int {
	if len(channels) == 0 {
		return 0
	}
	min := math.MaxInt
	for i := range channels {
		c := cap(channels[i])
		if c < min {
			min = c
		}
	}
	return min
}

// Max calcs maximal capacity of all channels.
func Max[T any](channels ...chan T) int {
	if len(channels) == 0 {
		return 0
	}
	max := 0
	for i := range channels {
		c := cap(channels[i])
		if c > max {
			max = c
		}
	}
	return max
}
