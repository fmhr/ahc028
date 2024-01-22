package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"sort"
	"time"
)

var start time.Time

// ./bin/main -cpuprofile cpuprof < tools/in/0000.txt
// go tool pprof -http=localhost:8888 bin/main cpuprof
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

//var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func main() {
	log.SetFlags(log.Lshortfile)
	///////////////////////////////
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	/////////////////////////////////////

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

var N, M int
var startPoint Point
var keyboard [][]byte
var words []string

func read(r io.Reader) {
	fmt.Fscan(r, &N, &M)
	fmt.Fscan(r, &startPoint.y, &startPoint.x)
	keyboard = make([][]byte, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(r, &keyboard[i])
	}
	words = make([]string, M)
	for i := 0; i < M; i++ {
		fmt.Fscan(r, &words[i])
	}
}

func solver() {
	read(os.Stdin)
	// wordsの順番を変更して、最終的な文字列を最小にする
	// 文字の重複によって、文字列は縮む
	// shortest Superstring problem
	// 前後の文字列と繋がっていない文字は順番を変更できるので、キーボードの位置を考慮して、順番を変える
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
	//log.Printf("len=%d\n", len(result))
	//str := greedyOrder(result, points, startPoint)
	str := beamSearchOrder(result, points, startPoint)
	rtn2, _ := dpRoot(str, points, startPoint)
	score(rtn2, startPoint)
	for i := 0; i < len(rtn2); i++ {
		fmt.Println(rtn2[i])
	}
}

var X int = 1

func shortestSuperstring(words []string, points [26][]Point) []string {
	for {
		var restart bool
		// k > 0 か k > 1　で考える
		for k := 4; k > 0; k-- {
			for i := 0; i < len(words); i++ {
				for j := 0; j < len(words); j++ {
					if i == j {
						continue
					}
					if words[i][len(words[i])-k:] == words[j][:k] {
						_, costi := dpRootCache(words[i], points)
						_, costj := dpRootCache(words[j], points)
						newWord := words[i] + words[j][k:]
						_, cost := dpRootCache(newWord, points)
						if costi+costj+X >= cost {
							//log.Println(words[i], words[j], costi, costj, cost)
							words[i] = newWord
							words[j] = words[len(words)-1]
							words = words[:len(words)-1]
							restart = true
							break
						}
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
	return words
}

// greedyで順番を決める
//func greedyOrder(words []string, points [26][]Point, start Point) string {
//size := len(words)
//minCost := math.MaxInt32
//minWord := ""
//minWordIndex := 0
//for i := 0; i < len(words); i++ {
//_, cst := dpRoot(words[i], points, start)
//_, baseCst := dpRootCache(words[i], points)
//cst -= baseCst
//if cst < minCost {
//minCost = cst
//minWordIndex = i
//minWord = words[i]
//}
//}
//rtn := minWord
//words[size-1], words[minWordIndex] = words[minWordIndex], words[size-1]
//words = words[:size-1]
//for len(words) > 0 {
//minCost = math.MaxInt32
//minWord = ""
//minWordIndex = 0
//for i := 0; i < len(words); i++ {
//_, cst := dpRootCache(rtn+words[i], points)
//_, baseCst := dpRootCache(words[i], points)
//cst -= baseCst
//if cst < minCost {
//minCost = cst
//minWordIndex = i
//minWord = words[i]
//}
//}
//words[len(words)-1], words[minWordIndex] = words[minWordIndex], words[len(words)-1]
//words = words[:len(words)-1]
//rtn += minWord
//}
//return rtn
//}

type Node struct {
	used [200]bool
	str  string
	cost int
}

func goalCheck(n *Node, m int) bool {
	for i := 0; i < m; i++ {
		if !n.used[i] {
			return false
		}
	}
	return true
}

func generateNodes(n Node, points [26][]Point, words []string) []Node {
	nodes := make([]Node, 0, 200)
	for i := 0; i < len(words); i++ {
		if n.used[i] {
			continue
		}
		_, cst := dpRootCache(n.str+words[i], points)
		_, baseCst := dpRootCache(words[i], points)
		cst -= baseCst
		var str string
		if len(n.str) > 1 && n.str[len(n.str)-1] == words[i][0] {
			str = n.str + words[i][1:]
		} else {
			str = n.str + words[i]
		}
		node := Node{n.used, str, n.cost + cst}
		node.used[i] = true
		nodes = append(nodes, node)
	}
	return nodes
}

func beamSearchOrder(words []string, points [26][]Point, start Point) string {
	beamWidth := 1
	initialNode := Node{[200]bool{}, "", 0}
	nodes := make([]Node, 0, 200)
	nodes = append(nodes, initialNode)
	nodesSub := make([]Node, 0, 200)
	for len(nodes) > 0 {
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].cost < nodes[j].cost
		})
		nodesSub = nodesSub[:0]
		for i := 0; i < min(beamWidth, len(nodes)); i++ {
			nextNodes := generateNodes(nodes[i], points, words)
			for j := 0; j < len(nextNodes); j++ {
				nodesSub = append(nodesSub, nextNodes[j])
			}
		}
		nodes = make([]Node, len(nodesSub))
		copy(nodes, nodesSub)
		if goalCheck(&nodes[0], len(words)) {
			break
		}
	}
	for i := 0; i < len(nodes); i++ {
		_, cst := dpRoot(nodes[i].str, points, start)
		nodes[i].cost = cst
	}
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].cost < nodes[j].cost
	})
	return nodes[0].str
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

// 一番短いルートを探す
const sizeN = 16 // N=15
func dpRoot(word string, points [26][]Point, startP Point) ([]Point, int) {
	dp := make([][sizeN][sizeN]int, len(word))
	root := make([][sizeN][sizeN]Point, len(word))
	for i := 0; i < len(word); i++ {
		for j := 0; j < sizeN; j++ {
			for k := 0; k < sizeN; k++ {
				dp[i][j][k] = math.MaxInt32
			}
		}
	}
	if startP.y != -1 && word[0] == keyboard[startP.y][startP.x] {
		dp[0][startP.y][startP.x] = 0
	} else {
		for i := 0; i < len(points[word[0]-'A']); i++ {
			dp[0][points[word[0]-'A'][i].y][points[word[0]-'A'][i].x] = 0
		}
	}
	for l := 1; l < len(word); l++ {
		a := word[l-1] - 'A'
		b := word[l] - 'A'
		for i := 0; i < len(points[a]); i++ {
			for j := 0; j < len(points[b]); j++ {
				cost := distance(points[a][i], points[b][j])
				if dp[l-1][points[a][i].y][points[a][i].x]+cost < dp[l][points[b][j].y][points[b][j].x] {
					dp[l][points[b][j].y][points[b][j].x] = dp[l-1][points[a][i].y][points[a][i].x] + cost
					root[l][points[b][j].y][points[b][j].x] = points[a][i]
				}
			}
		}
	}
	minCostIndex := 0
	minCost := dp[len(word)-1][points[word[len(word)-1]-'A'][0].y][points[word[len(word)-1]-'A'][0].x]
	for i := 1; i < len(points[word[len(word)-1]-'A']); i++ {
		if dp[len(word)-1][points[word[len(word)-1]-'A'][i].y][points[word[len(word)-1]-'A'][i].x] < minCost {
			minCostIndex = i
			minCost = dp[len(word)-1][points[word[len(word)-1]-'A'][i].y][points[word[len(word)-1]-'A'][i].x]
		}
	}
	rootPoint := make([]Point, len(word))
	rootPoint[len(word)-1] = points[word[len(word)-1]-'A'][minCostIndex]
	for i := len(word) - 2; i >= 0; i-- {
		rootPoint[i] = root[i+1][rootPoint[i+1].y][rootPoint[i+1].x]
	}
	return rootPoint, minCost
}

type DpRootCache struct {
	root []Point
	cost int
}

var dpRootCacheMap map[string]DpRootCache

func dpRootCache(word string, points [26][]Point) ([]Point, int) {
	if dpRootCacheMap == nil {
		dpRootCacheMap = make(map[string]DpRootCache)
	}
	if cache, ok := dpRootCacheMap[word]; ok {
		return cache.root, cache.cost
	} else {
		root, cost := dpRoot(word, points, Point{-1, -1})
		dpRootCacheMap[word] = DpRootCache{root, cost}
		return root, cost
	}
}

func score(ans []Point, start Point) {
	score := 10000
	cost := distance(start, ans[0]) + 1
	for i := 0; i < len(ans)-1; i++ {
		cost += distance(ans[i], ans[i+1]) + 1
	}
	log.Printf("score=%d cost=%d\n", score-cost, cost)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
