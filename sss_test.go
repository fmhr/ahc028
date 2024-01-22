package main

import (
	"log"
	"os"
	"testing"
)

func readSample() {
	file, err := os.Open("tools/in/0000.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	read(file)
}

func BenchmarkSSS(b *testing.B) {
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
}

func BenchmarkDpRoot(b *testing.B) {
	readSample()
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
	result, n := dpRoot("ACDGEATPHEPP", points, Point{-1, -1}, false)
	//log.Println(n)
	//log.Println(result)
	_, _ = result, n
}

// go test -bench BeamSearch -cpuprofile cpu.out -benchmem
// go tool pprof -http=":8080" cpu.out
func BenchmarkBeamSearch(b *testing.B) {
	readSample()
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
	b.ResetTimer()
	str := beamSearchOrder(words, points, startPoint)
	_ = str
}
