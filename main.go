package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
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

type VerseDetails struct {
	VfileID string
}

func main() {

	tmpl := template.Must(template.ParseFiles("main.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		v := VerseDetails{
			VfileID: r.FormValue("vfileid"),
		}

		// do something with v
		_ = v
		//		fmt.Println("v is \n", v)
		//		fmt.Println("vfileid is \n", v.VfileID)

		//Read entire file
		dat, err := ioutil.ReadFile(v.VfileID + ".txt")
		check(err)

		//Read and save each verse by Topic
		verses := strings.Split(string(dat), "TOPIC:")

		//Show all verses - Topic, Reference and Verse
		versecount := len(verses)
		CharPerDay := len(dat) / 30
		var topic, reference, verse string

		//Show the total number of characters of all the verses
		//fmt.Fprintln(w, "total char/30 is", CharPerDay)
		//fmt.Fprintln(w, "versecount is ", versecount)
		//		fmt.Fprintln(w, "Start program \n")

		//Show verses for current day or otherwise
		now := time.Now()
		DayOfMonth := now.Day()
		//		fmt.Fprintln(w, "DayOfMonth is ", DayOfMonth)
		StartDayCount := 0
		DayCount := 1
		var TotalCharNow int
		DayMatch := false

		//do until DayCount == DayOfMonth
		for !DayMatch {
			for j := 1; j < versecount; j++ {
				topic, reference, verse = TopRefVerse(verses[j])
				//				fmt.Fprintln(w, "j, DayCount: ", j, DayCount)
				//If the total characters so far are at the stopping point
				if TotalCharNow > (StartDayCount + CharPerDay) {
					DayCount++
					StartDayCount = TotalCharNow
					if DayCount > DayOfMonth {
						//						fmt.Fprint(w, "DayCount > DayOfMonth")
						break
					}
				}
				TotalCharNow = TotalCharNow + len(verses[j])
				if DayCount == DayOfMonth {
					//	fmt.Fprintln(w, "DayCount == DayOfMonth")
					topic, reference, verse = TopRefVerse(verses[j])
					fmt.Fprint(w, "Topic: ", topic)
					fmt.Fprint(w, "Reference: ", reference)
					fmt.Fprint(w, "Verse: ", verse)
					DayMatch = true
				}
			}
		}

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":8080", nil)
}

//Parse Topic, Reference and Verse from whole verse text
func TopRefVerse(wholeverse string) (string, string, string) {
	TopicRefVerse := strings.Split(wholeverse, "REFERENCE:")
	RefVerse := strings.Split(TopicRefVerse[1], "VERSE:")
	return TopicRefVerse[0], RefVerse[0], RefVerse[1]
}
