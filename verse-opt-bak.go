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
	Tdays   string
	Wday    string
}

func main() {

	tmpl := template.Must(template.ParseFiles("verse-opt.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		v := VerseDetails{
			VfileID: r.FormValue("vfileid"),
			Tdays:   r.FormValue("tdays"),
			Wday:    r.FormValue("wday"),
		}

		// do something with v
		_ = v
		fmt.Println("v is \n", v)
		fmt.Println("vfileid is \n", v.VfileID)
		fmt.Println("tdays is \n", v.Tdays)
		fmt.Println("today is \n", v.Wday)

		//Read entire file
		dat, err := ioutil.ReadFile(v.VfileID + ".txt")
		check(err)
		fmt.Println("v.VfileID is \n", string(dat))

		//Read and save each verse by Topic
		verses := strings.Split(string(dat), "TOPIC:")

		//Show all verses - Topic, Reference and Verse
		versecount := len(verses)
		//		CharPerDay := len(dat) / 30
		//i, err := strconv.Atoi("-42")
		//		CharPerDay := len(dat) / int(v.Tdays)
		tdays := "10"
		var ntdays = int32(tdays)
		CharPerDay := len(dat) / ntdays

		var topic, reference, verse string

		//Show the total number of characters of all the verses
		//		fmt.Fprintln(w, "total char/30 is", CharPerDay)
		//		fmt.Fprintln(w, "Total verses in verses is ", versecount)

		//Show verses for current day or otherwise
		now := time.Now()
		//		DayOfMonth := now.Day()
		fmt.Println("Day of month is \n", now)
		DayOfMonth := v.Wday
		StartDayCount := 0
		DayCount := 1
		var TotalCharNow int

		for j := 1; j < versecount; j++ {
			//If (total verses - present) < (30 - current day of month) then show last read
			//If (versecount - j) < (30 - DayOfMonth) then show last read
			if (versecount - j) < (30 - DayOfMonth) {
				topic, reference, verse = TopRefVerse(verses[j])
				fmt.Fprint(w, "Topic: ", topic)
				fmt.Fprint(w, "Reference: ", reference)
				fmt.Fprint(w, "Verse: ", verse)
				break
			}
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
			}
		}

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":8082", nil)
}

//Parse Topic, Reference and Verse from whole verse text
func TopRefVerse(wholeverse string) (string, string, string) {
	TopicRefVerse := strings.Split(wholeverse, "REFERENCE:")
	RefVerse := strings.Split(TopicRefVerse[1], "VERSE:")
	return TopicRefVerse[0], RefVerse[0], RefVerse[1]
}
