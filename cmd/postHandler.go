package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var data GestionPost
	err = json.Unmarshal(body, &data)
	createPostWithTitle(Post{
		Title: data.CreatePost,
	})
}

func displayPostVisible(w http.ResponseWriter, r *http.Request) {
	unhiddenPost := takePostUnHidden()
	jsonData, err := json.Marshal(unhiddenPost)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(jsonData)
}
