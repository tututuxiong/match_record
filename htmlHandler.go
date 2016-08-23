package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
	//	"io/ioutil"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	//	"github.com/PuerkitoBio/goquery"
)

func root_begin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
	url, _ := url.QueryUnescape(r.URL.String())

	if len(url) == 0 || url == "/" {
		t, _ := template.ParseFiles("./front_web/" + "main_page.html")
		t.Execute(w, nil)
		return
	} else {
		//fmt.Println("parse file " + url)
		if strings.Contains(url, "html") {
			if strings.Contains(url, "_team.html") {
				urlByte := []byte(url)
				name := urlByte[1:strings.Index(url, "_team.html")]
				t, _ := template.ParseFiles("./front_web/" + "team.html")
				p := PersonHtmlTemplate{Name: string(name)}

				t.Execute(w, p)
				return
			} else if strings.Contains(url, "_person.html") {
				urlByte := []byte(url)
				name := urlByte[1:strings.Index(url, "_person.html")]

				t, _ := template.ParseFiles("./front_web/" + "person.html")
				p := PersonHtmlTemplate{Name: string(name)}

				t.Execute(w, p)
				return
			} else if strings.Contains(url, "updateScore.html") {
				var s UpdateScorePageTemlate
				s.Date = r.FormValue("Date")
				s.MajorRound = r.FormValue("MajorRound")
				s.TypeId = r.FormValue("TypeId")
				s.SmallRound = r.FormValue("SmallRound")
				s.OurTeamId = r.FormValue("OurTeamId")
				s.EnemyTeamId = r.FormValue("EnemyTeamId")
				s.OurScore = r.FormValue("OurScore")
				s.EnemyScore = r.FormValue("EnemyScore")
				s.OurTeamName = r.FormValue("OurTeamId")
				s.EnemyTeamName = r.FormValue("EnemyTeamId")
				s.OurTeamName = getTeamNameById(s.OurTeamId)
				s.EnemyTeamName = getTeamNameById(s.EnemyTeamId)
				s.OurPlayer0 = r.FormValue("OurPlayer0")
				s.OurPlayer1 = r.FormValue("OurPlayer1")
				s.EnemyPlayer0 = r.FormValue("EnemyPlayer0")
				s.EnemyPlayer1 = r.FormValue("EnemyPlayer1")

				t, _ := template.ParseFiles("./front_web/" + "updateScore.html")
				t.Execute(w, s)
				return
			} else {
				fmt.Fprintln(w, nil)
			}
		}
	}
}

func getPersonData(name string) []personData_json {
	personDatas := getPersonRecordDatas(name)
	return personDatas
}

func handle_ajax(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	method := r.Form.Get("method")
	fmt.Println(method)

	if method == "getPersonData" {
		name := r.Form.Get("name")
		p := getPersonData(name)
		b, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	} else if method == "getPersonList" {
		p := getPersonList()
		b, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	} else if method == "getTeamInfo" {
		p := getTeamInfoList()
		b, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	} else if method == "getTeamMatchInfo" {
		name := r.Form.Get("name")
		p := getTeamRrcordInfos(name)
		b, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	} else if method == "getLatestResult" {
		p := getLatestRoundInfo()
		b, err := json.Marshal(p)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	}
}

func handleUpdateScore(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	Date := r.Form.Get("Date")
	MajorRound := r.Form.Get("MajorRound")
	SmallRound := r.Form.Get("SmallRound")
	EnemyTeamId := r.Form.Get("EnemyTeamId")
	OurTeamId := r.Form.Get("OurTeamId")
	OurScore := r.Form.Get("OurScore")
	EnemyScore := r.Form.Get("EnemyScore")

	handleUpdateScoreSql(OurTeamId, Date, MajorRound, SmallRound, OurScore)
	handleUpdateScoreSql(EnemyTeamId, Date, MajorRound, SmallRound, EnemyScore)

	http.Redirect(w, r, "./"+OurTeamId+"_team.html", http.StatusMovedPermanently)
}

func router() {
	http.HandleFunc("/", root_begin)
	//http.HandleFunc("/admin", root_begin)
	http.HandleFunc("/ajax", handle_ajax)
	http.HandleFunc("/updateScore", handleUpdateScore)
	http.Handle("/bm-all-utf8-0807-1815.s3db", http.FileServer(http.Dir("./")))
	http.Handle("/js/", http.FileServer(http.Dir("./front_web/")))
	http.Handle("/css/", http.FileServer(http.Dir("./front_web/")))
	http.Handle("/fonts/", http.FileServer(http.Dir("./front_web/")))
}

func main() {
	fmt.Println("Server begin ")

	Sql_open()

	port := os.Getenv("PORT")

	if port == "" {
		port = "5050"
	}
	fmt.Println("port: " + port)
	router()

	err := http.ListenAndServe(":"+port, nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
