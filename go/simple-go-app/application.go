package main

import (
    "bytes"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

func main() {
    var buf bytes.Buffer

    port := os.Getenv("PORT")
    if port == "" {
        port = "5000"
    }

    logfile := "/var/log/golang/golang-server.log"
    f, err := os.Create(logfile)
    if err != nil {
        fmt.Println("Error: cannot create log file:", logfile, err)
        os.Exit(1)
    }
    defer f.Close()
    log.SetOutput(f)

    const indexPage = "index.html"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            if buf, err := ioutil.ReadAll(r.Body); err == nil {
                log.Printf("Received message: %s\n", string(buf))
                body, err := json.Marshal(map[string]interface{}{
                    "application": "simple Go (go1.x) application",
                    "message": "application executed successfully!",
                    "method": "POST",
                    "environment": os.Getenv("ENVIRONMENT"),
                    "version":     "0.0.0",
                })
                if err != nil {
                    return Response{StatusCode: 404}, err
                }
                json.HTMLEscape(&buf, body)
                resp := Response{
                    StatusCode: 200,
                    IsBase64Encoded: false,
                    Body: buf.String(),
                    Headers: map[string]string{
                        "Content-Type": "application/json",
                    },
                }
                return resp, nil
            }
        } else {
            log.Printf("Serving %s to %s...\n", indexPage, r.RemoteAddr)
            http.ServeFile(w, r, indexPage)
        }
    })

    http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            log.Printf("Received 'foo' GET\n")
            body, err := json.Marshal(map[string]interface{}{
                "message": "foo GET executed successfully!",
                "method": "GET",
                "path": "foo",
            })
            if err != nil {
                return Response{StatusCode: 404}, err
            }
            json.HTMLEscape(&buf, body)
            resp := Response{
                StatusCode: 200,
                IsBase64Encoded: false,
                Body: buf.String(),
                Headers: map[string]string{
                    "Content-Type": "application/json",
                },
            }
            return resp, nil
        } else {
            log.Printf("Received 'foo' unsupported method: %s\n", r.Method)
            return Response{StatusCode: 404}
        }
    })

    http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "POST" {
            if buf, err := ioutil.ReadAll(r.Body); err == nil {
		log.Printf("Received 'bar' POST: %s\n", string(buf))
                body, err := json.Marshal(map[string]interface{}{
                    "message": "bar POST executed successfully!",
                    "method": "POST",
                    "path": "bar",
                    "body": string(buf),
                })
                if err != nil {
                    return Response{StatusCode: 404}, err
                }
                json.HTMLEscape(&buf, body)
                resp := Response{
                    StatusCode: 200,
                    IsBase64Encoded: false,
                    Body: buf.String(),
                    Headers: map[string]string{
                        "Content-Type": "application/json",
                    },
                }
                return resp, nil
            }
        } else {
            log.Printf("Received 'bar' unsupported method: %s\n", r.Method)
            return Response{StatusCode: 404}
        }
    })

    log.Printf("Listening on port %s\n\n", port)
    http.ListenAndServe(":"+port, nil)
}
