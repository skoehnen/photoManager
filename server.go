package main

import (
    "os"
    "crypto/md5"
    "strconv"
    "io"
    "fmt"
    "net/http"
    "time"
    "html/template"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
)

type Photo struct {
	name      string
}

// upload logic
func upload(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Upload called")
    fmt.Println("method:", r.Method)
    if r.Method == "GET" {
        crutime := time.Now().Unix()
        h := md5.New()
        io.WriteString(h, strconv.FormatInt(crutime, 10))
        token := fmt.Sprintf("%x", h.Sum(nil))
        t, err := template.ParseFiles("upload.gtpl")
	if( err != nil ){
	  fmt.Println(err)
	  return
	}
        t.Execute(w, token)
    } else {
        r.ParseMultipartForm(32 << 20)
        file, handler, err := r.FormFile("uploadfile")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
        fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile("./test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
    }
}

func view(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "This is where the pictures are shown")
    
    session, err := mgo.Dial("localhost")
    if err != nil {
       fmt.Println(err)
       return
    }
    c := session.DB("photos").C("photoList")

    var results []Photo

    err = c.Find(bson.M{"name": name}).All(&results)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
    fmt.Printf("Server started\n")
    http.HandleFunc("/upload", upload)
    http.HandleFunc("/view", view)
    http.ListenAndServe(":8080", nil)
}

