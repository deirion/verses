package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

//error check for file read
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/,ah/health", healthCheckHandler)
	log.Print("listening on port:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))

}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	//Read entire file
	dat, err := ioutil.ReadFile("emaw_verses.txt")
	check(err)

	//Read and save each verse by Topic
	//	onestring := make([]string, 0)
	//	totalchar := 0
	verses := strings.Split(string(dat), "TOPIC:")

	//Show all verses - Topic, Reference and Verse
	versecount := len(verses)
	totalchar := len(dat)
	CharPerDay := len(dat) / 30
	var topic, reference, verse string

	//Show the total number of characters of all the verses
	fmt.Fprintln(w, "Total characters in dat is ", totalchar)
	fmt.Fprintln(w, "total char/30 is", CharPerDay)
	fmt.Fprintln(w, "Total verses in verses is ", versecount)

	//Show verses for current day or otherwise
	now := time.Now()
	DayOfMonth := now.Day()
	//	DayOfMonth := 2
	StartDayCount := 0
	DayCount := 1
	var TotalCharNow int
	DayMatch := false

	//do until DayCount == DayOfMonth
	for !DayMatch {
		for j := 1; j < versecount; j++ {
			fmt.Println("j, present verse is ", j)
			topic, reference, verse = TopRefVerse(verses[j])
			fmt.Println("Reference: ", reference)
			fmt.Println("(versecount - j), (30 - DayOfMonth)", (versecount - j), (30 - DayOfMonth))

			//If the total characters so far are at the stopping point
			if TotalCharNow > (StartDayCount + CharPerDay) {
				DayCount++
				StartDayCount = TotalCharNow
				if DayCount > DayOfMonth {
					break
				}
			}
			TotalCharNow = TotalCharNow + len(verses[j])
			if DayCount == DayOfMonth {
				//				fmt.Fprint(w, "Total char through Verse# ", j, " is:", TotalCharNow, " \n")
				//				fmt.Fprint(w, "Verse# ", j, ": \n")
				topic, reference, verse = TopRefVerse(verses[j])
				fmt.Fprint(w, "Topic: ", topic)
				fmt.Fprint(w, "Reference: ", reference)
				fmt.Fprint(w, "Verse: ", verse)
				DayMatch = true

			}
		}
	}

	//Read and save each line
	//		lines := strings.Split(string(dat), "\n")
	//Display 1st five lines
	//		for j := 0; j <= 4; j++ {
	//			fmt.Fprint(w, "line - ", j, lines[j], "\n")
	//		}

	//Display 1st five characters of 1st line
	//		for j := 0; j <= 4; j++ {
	//			fmt.Fprint(w, "char - ", j, lines[0][j], "\n")
	//		}
} //end of func handle

//Parse Topic, Reference and Verse from whole verse text
func TopRefVerse(wholeverse string) (string, string, string) {
	TopicRefVerse := strings.Split(wholeverse, "REFERENCE:")
	RefVerse := strings.Split(TopicRefVerse[1], "VERSE:")
	return TopicRefVerse[0], RefVerse[0], RefVerse[1]
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "ok")
}
