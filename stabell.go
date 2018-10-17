package main

import (
	"encoding/json"
	"fmt"
	"github.com/fabioberger/chrome"
	"honnef.co/go/js/dom"
)

func main() {
	session := Session{
		chrome.NewChrome(),
	}

	el1 := dom.GetWindow().Document().GetElementByID("vm-1")
	el1.AddEventListener("click", true, session.SaveWindowSession)
}

type Session struct {
	chrome *chrome.Chrome
}

// Save current window tabs to storage
func (session Session) SaveWindowSession(event dom.Event) {
	fmt.Println("Triggered with event: ", event)
	// Get tabs in current window

	queryInfo := chrome.Object{
		"currentWindow": true,
	}

	session.chrome.Tabs.Query(queryInfo, session.saveTabs)
}

func (session Session) saveTabs(tabs []chrome.Tab) {
	/*urls*/_, err := generateJSON(tabs)
	if err != nil {
		//
	}

	/*fmt.Println(session.chrome.Storage.Sync.Get("urls"))
	session.chrome.Storage.Sync.Set("urls", string(urls))
	fmt.Println(session.chrome.Storage.Sync.Get("urls"))*/
}

func generateJSON(tabs []chrome.Tab) (mu []byte, err error) {
	var urls = make(map[int]string)

	for i := range tabs {
		urls[i] = tabs[i].Url
	}

	mu, err = json.Marshal(urls)

	return
}
