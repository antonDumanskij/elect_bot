package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Sdd struct {
	Components []Component `json:"components"`
}

type Component struct {
	Template string                          `json:"template_name"`
	Schedule map[string]map[string][][]Group `json:"schedule"`
}

type Group struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Type  string `json:"type"`
}

const CHAT_ID = 427693118

func main() {
	req, err := http.NewRequest("GET", "https://api.yasno.com.ua/api/v1/pages/home/schedule-turn-off-electricity", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("cache-control", "no-cache")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	dat := Sdd{}
	if err := json.Unmarshal(bodyText, &dat); err != nil {
		panic(err)
	}
	var numberDay = int(time.Now().Weekday())

	scheduleToday := dat.Components[3].Schedule["dnipro"]["group_1"][numberDay-1]

	for _, v := range scheduleToday {
		if v.Start-1 == time.Now().Hour() {
			as, _ := json.Marshal(v)
			sendMessage(as)
			return
		}
	}

	if time.Now().Hour() == 8 {
		as, _ := json.Marshal(scheduleToday)
		sendMessage(as)
	}
}

func sendMessage(as []byte) {
	var qp map[string]string = make(map[string]string)
	qp["chat_id"] = strconv.Itoa(CHAT_ID)
	qp["text"] = bytes.NewBuffer(as).String()

	q, err := json.Marshal(qp)
	if err != nil {
		return
	}

	reqT, _ := http.NewRequest("POST", "https://api.telegram.org/bot7124172157:AAEBzrTkjLtJguN8jpFBh0OduVBx319j3hA/sendMessage", bytes.NewBuffer(q))
	reqT.Header.Add("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(reqT)
	if err != nil {
		log.Fatal(err)
	}
}
