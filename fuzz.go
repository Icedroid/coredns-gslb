// +build fuzz

package gslb

import (
	"github.com/coredns/coredns/plugin/pkg/fuzz"
)

// Fuzz fuzzes cache.
func Fuzz(data []byte) int {
	w := Gslb{}
	return fuzz.Do(w, data)
}
