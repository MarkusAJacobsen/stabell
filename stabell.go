package main

import (
	"encoding/json"
	"fmt"
	"github.com/fabioberger/chrome"
	"github.com/gopherjs/gopherjs/js"
	"github.com/pkg/errors"
	"honnef.co/go/js/dom"
	"time"
)

const (
	StorageKey  string = "stabell-storage"
	JSUndefined string = "undefined"
)

var (
	cStorage = js.Global.Get("localStorage")
)

func main() {
	ctx := Ctx{
		chrome.NewChrome(),
		dom.GetWindow().Document().GetElementByID("status"),
	}

	el1 := dom.GetWindow().Document().GetElementByID("vm-1")
	el1.AddEventListener("click", true, ctx.SaveWindowSession)

	el2 := dom.GetWindow().Document().GetElementByID("vm-2")
	el2.AddEventListener("click", true, GetSavedSessionsHandler)
}

type Ctx struct {
	chrome *chrome.Chrome
	status dom.Element
}

type Storage []Session

type Session struct {
	Ts  int64    `json:"ts"`
	Url []string `json:"url"`
}

func (ctx Ctx) SetStatus(status string) {
	now := time.Now().Format("01/02 15:04:05")
	ctx.status.SetInnerHTML(status + " " + now)
}

// Save current window tabs to cStorage
func (ctx Ctx) SaveWindowSession(event dom.Event) {
	queryInfo := chrome.Object{
		"currentWindow": true,
	}

	ctx.chrome.Tabs.Query(queryInfo, ctx.saveTabs)
}

func (ctx Ctx) saveTabs(tabs []chrome.Tab) {
	urls := generateJSON(tabs)

	session := Session{
		time.Now().Unix(),
		mapToStringArray(urls),
	}

	restore, err := getFromChromeStorage()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if isSame(restore, session) {
		ctx.SetStatus("Session already saved")
		return
	}

	storage := append(restore, session)

	ctx.saveToChromeStorage(storage)
}

func (ctx Ctx) saveToChromeStorage(storage Storage) (err error) {
	bytes, err := json.Marshal(storage)
	if err != nil {
		return
	}

	cStorage.Set(StorageKey, string(bytes))

	ctx.SetStatus("Session saved")
	return
}

func getFromChromeStorage() (storage Storage, err error) {
	res := cStorage.Get(StorageKey).String()
	storage = Storage{}

	if res != JSUndefined {
		err = json.Unmarshal([]byte(res), &storage)
		if err != nil {
			errors.Wrap(err, "Thrown in getFromChromeStorage")
			return
		}
	}
	return
}

func generateJSON(tabs []chrome.Tab) (urls map[int]string) {
	urls = make(map[int]string)

	for i := range tabs {
		urls[i] = tabs[i].Url
	}

	return
}

func mapToStringArray(val map[int]string) (con []string) {
	con = make([]string, 0, len(val))

	for _, value := range val {
		con = append(con, value)
	}

	return
}

func GetSavedSessionsHandler(event dom.Event) {
	res, err := getFromChromeStorage()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(res)
}

func isSame(storage Storage, session Session) (same bool) {
	/*for i := range storage {
		for j := range storage[i].Url {
			if !strings.EqualFold(storage[i].Url[j], session.Url[j]) {
				continue
			}
		}
	}*/
	return false
}
