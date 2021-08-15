package api_helper

import (
	// "crawler/c_help"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Receive_data(w http.ResponseWriter, r *http.Request) {

}

func Send_data(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// for c := range data {
	// 	fmt.Fprint(w, "<p> %s </p>", data[c].About)
	// }
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h2> Could not find the page</h2>")
	// w.Write([]byte(`{'status': 'first_phase_done', 'message': 'please send your OTP'}`))

}

func Initialise_server() {
	r := mux.NewRouter()
	s := r.Methods("GET").Subrouter()
	s.NotFoundHandler = http.HandlerFunc(NotFound)
	// s.HandleFunc("/", handlerFunc).Methods(http.MethodGet)
	// s.HandleFunc("/name/{name}", post).Methods("GET").Schemes("http")
	// s.HandleFunc("/dog", handlerFunc2)
	r.HandleFunc("/data", Send_data).Methods("POST")
	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:3000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func main() {
	fmt.Println("Api Helper")
}
