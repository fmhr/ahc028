package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strings"
	"time"
)

var start time.Time

func main() {
	log.SetFlags(log.Lshortfile)
	start = time.Now()
	solver()
	log.Printf("time=%f\n", time.Since(start).Seconds())
}

type Point struct {
	y, x int
}

func (p Point) String() string {
	return fmt.Sprintf("%d %d", p.y, p.x)
}

func solver() {
	var N, M int
	fmt.Scan(&N, &M)
	var start Point
	fmt.Scan(&start.y, &start.x)
	keyboard := make([][]byte, N)
	for i := 0; i < N; i++ {
		fmt.Scan(&keyboard[i])
	}
	words := make([]string, M)
	for i := 0; i < M; i++ {
		fmt.Scan(&words[i])
	}
	initial := make([]string, M)
	copy(initial, words)
	log.Println(N, M, start)
	// wordsの順番を変更して、最終的な文字列を最小にする
	// 文字の重複によって、文字列は縮む
	// shortest Superstring problem
	// 前後の文字列と繋がっていない文字は順番を変更できるので、キーボードの位置を考慮して、順番を変える
	result := shortestSuperstring(words)
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
	log.Printf("len=%d\n", len(result))
	//log.Println(minRoot(result[0], points))
	//greedyRoot(result[0], points)
	str := ""
	for i := 0; i < len(result); i++ {
		str += result[i]
		//log.Println(result[i])
	}
	log.Printf("total=%d\n", len(str))
	rtn := SARoot(str, points)
	for i := 0; i < len(rtn); i++ {
		fmt.Println(points[str[i]-'A'][rtn[i]])
	}
}

func shortestSuperstring(words []string) []string {
	initial := make([]string, len(words))
	copy(initial, words)
	best := make([]string, len(words))
	copy(best, words)
	bestStrLen := cntStringsLen(best)
	for z := 0; z < 5; z++ {
		words := make([]string, len(initial))
		copy(words, initial)
		rand.Shuffle(len(words), func(i, j int) {
			words[i], words[j] = words[j], words[i]
		})
		for {
			var restart bool
			// k > 0 か k > 1　で考える
			for k := 4; k > 0; k-- {
				for i := 0; i < len(words); i++ {
					for j := 0; j < len(words); j++ {
						if i == j {
							continue
						}
						if len(words[i]) > len(words[j]) {
							if strings.Contains(words[i], words[j]) {
								words[j] = words[len(words)-1]
								words = words[:len(words)-1]
								restart = true
								break
							}
						} else if len(words[i]) < len(words[j]) {
							if strings.Contains(words[j], words[i]) {
								words[i] = words[len(words)-1]
								words = words[:len(words)-1]
								restart = true
								break
							}
						}
						if words[i][len(words[i])-k:] == words[j][:k] {
							newWord := words[i] + words[j][k:]
							words[i] = newWord
							words[j] = words[len(words)-1]
							words = words[:len(words)-1]
							restart = true
							break
						}
					}
					if restart {
						break
					}
				}
				if restart {
					break
				}
			}
			if !restart {
				break
			}
		}
		if cntStringsLen(words) < bestStrLen {
			bestStrLen = cntStringsLen(words)
			best = make([]string, len(words))
			copy(best, words)
		}
	}
	return best
}

func cntStringsLen(words []string) int {
	length := 0
	for i := 0; i < len(words); i++ {
		length += len(words[i])
	}
	return length
}

func greedyRoot(word string, points [26][]Point) []int {
	root := make([]int, len(word))
	for i := 0; i < len(word); i++ {
		root[i] = rand.Intn(len(points[word[i]-'A']))
	}
	best := rootLength(word, root, points)
	for i := 0; i < 100; i++ {
		w := rand.Intn(len(word))
		old := root[w]
		n := rand.Intn(len(points[word[w]-'A']))
		root[w] = n
		newLength := rootLength(word, root, points)
		if newLength < best {
			best = newLength
		} else {
			root[w] = old
		}
		//log.Println(best, newLength, root)
	}
	return root
}

const start_temp = 5.0
const end_temp = 0.0
const maxTimeSeconsds = 1.9

func SARoot(word string, points [26][]Point) []int {
	iterations := 0
	wordNum := make([]int, len(word))
	for i := 0; i < len(word); i++ {
		wordNum[i] = int(word[i] - 'A')
	}
	currentSolution := make([]int, len(word))
	for i := 0; i < len(word); i++ {
		currentSolution[i] = rand.Intn(len(points[wordNum[i]]))
	}
	bestSolution := make([]int, len(word))
	copy(bestSolution, currentSolution)
	newSolution := make([]int, len(word))
	best := float64(rootLength(word, currentSolution, points))
	for {
		currentTime := time.Since(start).Seconds()
		if currentTime > maxTimeSeconsds {
			break
		}
		copy(newSolution, currentSolution)
		w := rand.Intn(len(word))
		n := rand.Intn(len(points[wordNum[w]]))
		newSolution[w] = n
		currentEnergy := float64(rootLength(word, currentSolution, points))
		newEnergy := float64(rootLength(word, newSolution, points))

		// 受理確率を計算
		temp := start_temp + (end_temp-start_temp)*currentTime/maxTimeSeconsds
		acceptanceProbability := math.Exp((currentEnergy - newEnergy) / temp)
		if newEnergy <= currentEnergy || rand.Float64() < acceptanceProbability {
			currentEnergy = newEnergy
			copy(currentSolution, newSolution)
		}

		if currentEnergy < best {
			best = currentEnergy
			copy(bestSolution, newSolution)
		}
		//log.Println(best, newEnergy, currentEnergy, temp, acceptanceProbability)
		iterations++
	}
	log.Printf("iterations=%d\n", iterations)
	return bestSolution
}

func rootLength(word string, root []int, points [26][]Point) int {
	length := 0
	for i := 0; i < len(word)-1; i++ {
		length += distance(points[word[i]-'A'][root[i]], points[word[i+1]-'A'][root[i+1]])
	}
	return length
}

func distance(p1, p2 Point) int {
	return abs(p1.y-p2.y) + abs(p1.x-p2.x)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
