package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
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
	// Main menu
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Quit", func() { a.Quit() }),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			dialog.ShowCustom("About", "Close", container.NewVBox(
				widget.NewLabel("Welcome to Gopher, a simple Desktop app created in Go with Fyne."),
				widget.NewLabel("Version: v0.1"),
				widget.NewLabel("Author: Aurélie Vache"),
			), window)
		}))
	mainMenu := fyne.NewMainMenu(
		fileMenu,
		helpMenu,
	)
	window.Resize(fyne.Size{500, 500})
	window.SetMainMenu(mainMenu)
	hello = widget.NewLabel("hello beka")
	res, err := fyne.LoadResourceFromPath("images/logo.png")
	if err != nil {
		log.Fatalln(err)
	}
	theme, err := fyne.LoadResourceFromPath("images/theme.jpeg")
	themeImg := canvas.NewImageFromResource(theme)
	themeImg.SetMinSize(fyne.Size{500, 500})
	window.SetIcon(res)
	window.SetContent(container.NewVBox(
		hello,
		themeImg,
		widget.NewButton("get", func() {
			hello.SetText("Hi there")
			a.SendNotification(fyne.NewNotification("Обновление", "у вас новое сообщение"))
		}),
		widget.NewButton("Обновить данные?", func() {
			hello.SetText("")
		}),
	))
	hello.Move(fyne.NewPos(-200, 400))
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
	m := map[string]string{
		"status":    "ok",
		"error":     "",
		"deqystvie": "текст изменен",
	}

	data, _ := json.Marshal(m)
	w.Write(data)
	hello.SetText(t.Text)
}
