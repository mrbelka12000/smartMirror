package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"io/ioutil"
	"log"
	"net/http"
)

var hello *widget.Label

type test struct {
	Text string `json:"text"`
}

func main() {
	r := Routers()
	a := app.New()
	window := a.NewWindow("Smart Mirror")

	hello = widget.NewLabel("hello beka")
	res, err := fyne.LoadResourceFromPath("../Загрузки/logo.png")
	if err != nil {
		log.Fatalln(err)
	}
	window.SetIcon(res)
	window.SetContent(container.NewVBox(
		hello,
		widget.NewButton("get", func() {
			hello.SetText("Hi there")
		}),
	))
	go func() {
		log.Println("server started")
		log.Fatalln(http.ListenAndServe(":8080", r))
	}()
	window.ShowAndRun()

}

func Routers() *http.ServeMux {
	r := http.DefaultServeMux
	r.Handle("/set", http.HandlerFunc(Main))
	return r
}

func Main(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	t := &test{}
	if err := json.Unmarshal(body, &t); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println(t.Text)
	w.Write([]byte("vse ok"))
	hello.SetText(t.Text)
}
