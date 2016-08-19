package main

import (
	"fmt"
	"html/template"
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
		fmt.Println("parse file " + url)
		if strings.Contains(url, "html") {
			if strings.Contains(url, "_team.html") {
				urlByte := []byte(url)
				name := urlByte[1:strings.Index(url, "_team.html")]
				t, _ := template.ParseFiles("./front_web/" + "team.html")
				p := PersonHtmlTemplate{Name: string(name)}
				fmt.Println(string(name))
				t.Execute(w, p)
				return
			}

			if strings.Contains(url, "_person.html") {
				urlByte := []byte(url)
				name := urlByte[1:strings.Index(url, "_person.html")]
				fmt.Println(string(name))
				t, _ := template.ParseFiles("./front_web/" + "person.html")
				p := PersonHtmlTemplate{Name: string(name)}
				fmt.Println(p)
				t.Execute(w, p)
				fmt.Println()
				return
			}

			t, _ := template.ParseFiles("./front_web/" + "person.html")
			t.Execute(w, t)

		}
	}
}

func getPersonData(name string) []personData_json {

	fmt.Println("person name is " + name)

	personDatas := getPersonRecordDatas(name)
	return personDatas
}

func handle_ajax(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.String())
	r.ParseForm()
	method := r.Form.Get("method")
	name := r.Form.Get("name")

	if method == "getPersonData" {
		p := getPersonData(name)
		b, err := json.Marshal(p)
		fmt.Println(string(b))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	} else if method == "getPersonList" {
		p := getPersonList()
		b, err := json.Marshal(p)
		fmt.Println(string(b))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	} else if method == "getTeamInfo" {
		p := getTeamInfoList()
		b, err := json.Marshal(p)
		fmt.Println(string(b))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	} else if method == "getTeamMatchInfo" {
		p := getTeamRrcordInfos(name)
		b, err := json.Marshal(p)
		fmt.Println(string(b))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintf(w, string(b))
	}

}

func router() {
	http.HandleFunc("/", root_begin)
	http.HandleFunc("/ajax", handle_ajax)
	http.Handle("/js/", http.FileServer(http.Dir("./front_web/")))
	http.Handle("/css/", http.FileServer(http.Dir("./front_web/")))
	http.Handle("/fonts/", http.FileServer(http.Dir("./front_web/")))
}

func main() {
	fmt.Println("Server begin ")
	router()
	err := http.ListenAndServe(":10101", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
