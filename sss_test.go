package main

import (
	"os"
	"testing"
)

func BenchmarkSSS(b *testing.B) {
	file, err := os.Open("tools/in/0000.txt")
	if err != nil {
		b.Error(err)
	}
	defer file.Close()
	read(file)

	var points [26][]Point
	for i := 0; i < 26; i++ {
		for j := 0; j < N; j++ {
			for k := 0; k < N; k++ {
				if keyboard[j][k] == byte('A'+i) {
					points[i] = append(points[i], Point{j, k})
				}
			}
		}
	}
	result := shortestSuperstring(words, points)
	_ = result
	//	str := beamSearchOrder(result, points, startPoint)
	//
	// _ = str
	// log.Println(str)
}
