package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	RankLimit = 500
)

var Verbose bool

// Specification
//	search target is current directory.
//	source target is first argument.
//	search filename extension is ".jpg" only.
//	exclude source target from search target files.
func main() {
	var directory string
	flag.StringVar(&directory, "d", ".", "target searching directory")
	flag.BoolVar(&Verbose, "v", false, "output verbose")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: Required at least an arguments")
		os.Exit(1)
	}

	target := args[0]

	d, err := os.Open(directory)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	colorBasedAnalyze(target, d)

}

func colorBasedAnalyze(target string, targetDir *os.File) {
	target_h, err := GetPartedHistogram(target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	verboseLog("---> Loading Target %s\n", target)
	//target_h.InfoShow()

	distanceMap := make(map[string]float64)

	fInfos, err := targetDir.Readdir(-1)
	for _, fInfo := range fInfos {
		if target == fInfo.Name() {
			continue
		}
		if strings.HasSuffix(fInfo.Name(), ".jpg") {
			h, err := GetPartedHistogram(fInfo.Name())
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			d, err := Distance(Normalize(HistToVector(target_h)),
				Normalize(HistToVector(h)))
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			distanceMap[fInfo.Name()] = d
			verboseLog("---> Loading %s\n", fInfo.Name())
			//h.InfoShow()
		}
	}

	// rank array for sorting map
	rank := make(PairList, len(distanceMap))
	var idx int
	for k, v := range distanceMap {
		rank[idx] = Pair{
			Filename:   k,
			Similarity: v,
		}
		idx++
	}
	sort.Sort(rank)

	for i, v := range rank {
		if i >= RankLimit {
			break
		}
		if Verbose {
			fmt.Printf("⭐️ %2d: %0.5f\t%s\n", i+1, v.Similarity, v.Filename)
		} else {
			fmt.Println(v.Filename)
		}
	}
	//fmt.Println(rank)
}

func verboseLog(format string, a ...interface{}) (int, error) {
	if Verbose {
		return fmt.Printf(format, a)
	}
	return 0, nil
}

// code for sorting map
// --------------------------------
type Pair struct {
	Filename   string
	Similarity float64
}

type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}

func (p PairList) Less(i, j int) bool {
	return p[i].Similarity < p[j].Similarity
}

func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// --------------------------------
