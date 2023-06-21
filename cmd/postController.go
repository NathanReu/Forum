package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func deletePost(post Post) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	query := "DELETE FROM posts WHERE title = ?"
	_, err = db.Exec(query, post.Title)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Post delete successfully")
}

func createPostWithTitle(post Post) {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	postCreate := `INSERT INTO posts (title) VALUES (?)`
	_, errCreate := db.Exec(postCreate, post.Title)
	if errCreate != nil {
		fmt.Println(err)
	}
	fmt.Println("Post created successfully")
}

func likePost(user_ID string, post_ID string) bool {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?", user_ID, post_ID).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if count > 0 {
		return false
	}
	_, err = db.Exec("INSERT INTO likes (user_id, post_id) VALUES (?, ?)", user_ID, post_ID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var countLike int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", post_ID).Scan(&countLike)
	_, err = db.Exec("UPDATE posts SET like = ? WHERE id = ?", countLike, post_ID)
	if err != nil {
		fmt.Println(err)
	}
	return true
}
func dislikePost(user_ID string, post_ID string) bool {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE user_id = ? AND post_id = ?", user_ID, post_ID).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if count > 0 {
		return false
	}
	_, err = db.Exec("INSERT INTO dislike (user_id, post_id) VALUES (?, ?)", user_ID, post_ID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var countDislike int
	err = db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ?", post_ID).Scan(&countDislike)
	_, err = db.Exec("UPDATE posts SET dislike = ? WHERE id = ?", countDislike, post_ID)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

func createCommentController(comment Comment) bool {
	if comment.Text == "" || comment.IDPost == 0 || comment.IDCreator == 0 {
		return true
	}
	return false
}

func uploadImage(w http.ResponseWriter, r *http.Request) { // c'est un return string

	if testMethod(w, r, http.MethodPost) {
		http.Error(w, "Failed to load fonction (method wrong)", http.StatusBadRequest)
		return
		//return "Failed to load fonction (method wrong)"
	}

	err, file, handlers := getDataToFormUploadImage(w, r)
	if err {
		http.Error(w, "Failed to load data (data type is may be wrong)", http.StatusBadRequest)
		return
		//return "Failed to load data (data type is may be wrong)"
	}

	message, err := createImageToFolder(w, file, handlers)
	if err {
		fmt.Println("Err : ", message)
		return
		//return "Err : ", message

	}

	fmt.Println("success")
	fmt.Println(message)
	return
	//return message

}

func testMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		return true
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	return false
}

func getDataToFormUploadImage(w http.ResponseWriter, r *http.Request) (bool, multipart.File, *multipart.FileHeader) {
	file, handlers, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return true, nil, nil
	}

	if handlers.Header.Get("Content-Type")[0:5] != "image" {
		createErrorMessage("C'est pas une image", 400, w)
		return true, nil, nil
	}
	defer file.Close()

	return false, file, handlers
}

func createImageToFolder(w http.ResponseWriter, file multipart.File, handlers *multipart.FileHeader) (string, bool) {
	dst, err := os.Create("templates/assets/img/imagePost/" + handlers.Filename)
	if err != nil {
		http.Error(w, "Error to copy Image", http.StatusInternalServerError)
		return "Error to copy Image", true
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error to create Image", http.StatusInternalServerError)
		return "Error to create Image", true
	}
	createSuccessfulMessage("File uploaded successfully", 201, w)
	return "assets/img/imagePost/" + handlers.Filename, false
}

func getDataEditPost(r *http.Request) Post {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error : ", err)
		return Post{ID: -1}
	}
	var data Post

	err = json.Unmarshal(body, &data)

	fmt.Println("LA DATA DE TES MROTS C'EST : ", data)

	return data
}

func changedDataPost(post Post, postInput Post) Post {
	postInput.ID = post.ID

	if postInput.Title == "" {
		postInput.Title = post.Title
	}
	if postInput.Texte == "" {
		postInput.Texte = post.Texte
	}
	if postInput.Photo == "" {
		postInput.Photo = post.Photo
	}

	return postInput
}

func takeInfoPostId(id int) []map[string]interface{} {
	db, err := sql.Open("sqlite3", "./forum.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	query := "SELECT * FROM posts WHERE id = ?"
	rows, err := db.Query(query, id)
	var result []map[string]interface{}
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Photo, &post.Title, &post.Texte, &post.Hidden, &post.Like, &post.Dislike, &post.Signalement, &post.Categorie, &post.Ban, &post.Archived, &post.IDCreator, &post.NameCreator)
		if err != nil {
			fmt.Println(err)
		}
		userData := make(map[string]interface{})
		userData["title"] = post.Title
		userData["text"] = post.Texte
		userData["like"] = post.Like
		userData["dislike"] = post.Dislike
		result = append(result, userData)
	}
	return result
}
