package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/xrash/smetrics"
)

// TODO Refactor big hairy main - Apply boy scout rule
func main() {
	fmt.Println("Simmer down...")
	boostThreshold := 0.7 // boostThreshold = minimum score for a string that gets boosted. This value was set to 0.7 in Winkler's papers.
	prefixSize := 4       // prefixSize = size of the initial prefix considered. This value was set to 4 in Winkler's papers.
	verbose := true

	// TODO Parse arguments here

	//Read source strings from a file
	source, err := os.Open("resources/test_strings.txt")
	exitOnError(err, "Oops cannot find source")

	target, err := os.Open("resources/ref_strings.txt")
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

	bestMatches := [][]string{}
	for sourceScanner.Scan() {
		sourceStr := sourceScanner.Text()
		bestMatchDist := 0.0
		bestMatchStr := ""
		for targetStr := range targetStrings {
			dist := smetrics.JaroWinkler(sourceStr, targetStr, boostThreshold, prefixSize)
			if dist > bestMatchDist {
				bestMatchDist = dist
				bestMatchStr = targetStr
			}
			if verbose {
				fmt.Print(sourceStr, " ", targetStr, " ", dist, "\n")
			}
		}
		bestMatches = append(bestMatches, []string{sourceStr, bestMatchStr})
		if verbose {
			fmt.Print("best match ", bestMatchStr, "\n")
		}
	}

	if err := sourceScanner.Err(); err != nil {
		// TODO do sth with error
	}

	//At this point we have the best matches
	// TODO write out to a file
	for i := range bestMatches {
		fmt.Print(bestMatches[i][0], ",", bestMatches[i][1], "\n")
	}

}

func exitOnError(e error, msg string) {
	if e != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}
