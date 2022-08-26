package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

type EmojiJson struct {
	Emojis []Emoji `json:"emoji"`
}
type Emoji struct {
	Name            string        `json:"name"`
	IsAlias         int           `json:"is_alias"`
	AliasFor        string        `json:"alias_for"`
	URL             string        `json:"url"`
	Created         int           `json:"created"`
	TeamID          string        `json:"team_id"`
	UserID          string        `json:"user_id"`
	UserDisplayName string        `json:"user_display_name"`
	AvatarHash      string        `json:"avatar_hash"`
	CanDelete       bool          `json:"can_delete"`
	IsBad           bool          `json:"is_bad"`
	Synonyms        []interface{} `json:"synonyms"`
}

func main() {
	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	emojiJson := new(EmojiJson)
	if err := json.Unmarshal(b, &emojiJson); err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	if err := os.Mkdir(`emoji`, 0750); err != nil && os.IsExist(err) {
		log.Fatal(err)
	}
	for _, emoji := range emojiJson.Emojis {
		wg.Add(1)
		go func(wg *sync.WaitGroup, e Emoji) {
			defer wg.Done()
			url, err := url.Parse(e.URL)
			if err != nil {
				log.Println(err)
				return
			}
			resp, err := http.Get(url.String())
			if err != nil {
				log.Println(err)
				return
			}
			ext := filepath.Ext(url.Path)
			fileName := fmt.Sprintf(`emoji/%s_%s%s`, e.Name, e.UserDisplayName, ext)
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Print(err)
				return
			}
			if err := os.WriteFile(fileName, b, 0666); err != nil {
				log.Print(err)
			}
			log.Println("saved", fileName)
		}(wg, emoji)
	}
	wg.Wait()
}
