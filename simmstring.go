package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/xrash/smetrics"
)

// TODO Refactor big hairy main
// TODO make it work for numcases > 1
// Example
// go run simmstring.go -source resources\test_strings.txt -target resources\ref_strings.txt -v
func main() {
	//
	// Jaro Winkler parameters
	//
	boostThreshold := 0.7 // boostThreshold = minimum score for a string that gets boosted. This value was set to 0.7 in Winkler's papers.
	prefixSize := 4       // prefixSize = size of the initial prefix considered. This value was set to 4 in Winkler's papers.

	//
	// ARGUMENTS
	//
	var src string
	var tar string
	var verbose bool
	nummatches := 1

	flag.StringVar(&src, "source", "", "(Required) Path to file with source strings.")
	flag.StringVar(&tar, "target", "", "(Required) Path to file with target strings.")
	flag.BoolVar(&verbose, "v", false, "(Optional) Verbose mode.")

	flag.Parse()

	if src == "" || tar == "" {
		printUsage()
		os.Exit(1)
	}

	//Read source strings from a file
	source, err := os.Open(src)
	exitOnError(err, "Oops cannot find source")

	target, err := os.Open(tar)
	exitOnError(err, "Oops cannot find target")

	defer target.Close()
	targetScanner := bufio.NewScanner(target)

	//Get unique target strings
	//Use a map to get the set of unique target strings
	targetStrings := make(map[string]bool)
	for targetScanner.Scan() {
		targetStr := targetScanner.Text()
		if verbose {
			if targetStrings[targetStr] {
				fmt.Print(targetStr, " is already in targetStrings map\n") // debug
			} else {
				fmt.Print("Adding ", targetStr, " to targetStrings map\n") // debug
			}
		}
		targetStrings[targetStr] = true
	}

	if err := targetScanner.Err(); err != nil {
		// TODO do sth with error
	}

	defer source.Close()
	sourceScanner := bufio.NewScanner(source)

	//bestMatches := [][]string{}
	for sourceScanner.Scan() {
		sourceStr := sourceScanner.Text()
		//as long as we keep bbestMatches sorted we know lowerLim/upperLim match
		//but we do not want to sort it each and every time
		var bbestMatches []SimPair
		//bestMatch stuff will be obsoleted by bbestMatch stuff
		//bestMatchDist := 0.0
		//bestMatchStr := ""
		//lowerLimMatch
		lowerLimDist := 1.0
		//start loop target strings
		for targetStr := range targetStrings {
			dist := smetrics.JaroWinkler(sourceStr, targetStr, boostThreshold, prefixSize)
			if len(bbestMatches) < nummatches {
				bbestMatches = append(bbestMatches, SimPair{sourceStr, targetStr, dist})
				if dist < lowerLimDist {
					lowerLimDist = dist
				}
			}
			if len(bbestMatches) == nummatches {
				// at this point there must be as many targets as nummatches
				if dist > lowerLimDist {
					// if nummatches = 1
					//  set new match as lowerLimDist
					// if nummatches > 1
					//  sort bbestMatches
					//  evict worst match
					//  second worst match as lowerLimDist
					if nummatches == 1 {
						// TODO can we just set bbestMatches anew?
						bbestMatches = []SimPair{SimPair{sourceStr, targetStr, dist}}
						lowerLimDist = dist
					}
				}
			}
			//if dist > bestMatchDist {
			//	bestMatchDist = dist
			//	bestMatchStr = targetStr
			//}
			//if verbose {
			//	fmt.Print(sourceStr, " ", targetStr, " ", dist, "\n")
			//}
		}
		// end loop target strings
		//print results to std out for each source str
		for i := range bbestMatches {
			out := bbestMatches[i]
			fmt.Print(out.String())
		}
		//bestMatches = append(bestMatches, []string{sourceStr, bestMatchStr})
		//if verbose {
		//	fmt.Print("best match ", bestMatchStr, "\n")
		//}
	}

	if err := sourceScanner.Err(); err != nil {
		// TODO do sth with error
	}

	/* test sorting of SimPair
	var sp []*SimPair
	sp = append(sp, &SimPair{Source: "a,a", Target: "b", Distance: 2.5})
	sp = append(sp, &SimPair{Source: "a", Target: "c", Distance: 1.5})
	fmt.Println(sp[0].String())
	fmt.Println(sp)
	sort.Sort(ByDistance(sp))
	fmt.Println(sp)
	*/

	//At this point we have the best matches
	// TODO write out to a file
	// or just rely on piping to output ... ?
	// what is more efficient ?
	//for i := range bestMatches {
	//	fmt.Print(bestMatches[i][0], ",", bestMatches[i][1], "\n")
	//}

}

// SimPair is pair of strings, a Source and a Target string, along with the Distance between that pair of strings
type SimPair struct {
	Source   string
	Target   string
	Distance float64
}

// String writes SimPair as CSV string
func (s SimPair) String() string {
	record := []string{s.Source, s.Target, strconv.FormatFloat(s.Distance, 'f', -1, 64)}
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	w.Write(record)
	w.Flush()
	str := b.String()
	return str
}

// ByDistance is for sorting SimPairs by distance
type ByDistance []*SimPair

// Len method of sort interface for ByDistance
func (d ByDistance) Len() int { return len(d) }

// Swap method of sort interface for ByDistance
func (d ByDistance) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

// Less method of sort interface for ByDistance
func (d ByDistance) Less(i, j int) bool { return d[i].Distance < d[j].Distance }

// TODO it's a bit ugly that exit on error prints usage
func exitOnError(e error, msg string) {
	if e != nil {
		fmt.Print(msg, "\n\n")
		printUsage()
		os.Exit(1)
	}
}

// printUsage prints help output
func printUsage() {
	println("simmstring version 0.1-SNAPSHOT")
	println("")
	println("Usage:")
	flag.PrintDefaults()
	println("")
	println("Examples:")
	println("  simmstring -source resoures\\test_strings.txt -target resources\\ref_strings.txt")
}
