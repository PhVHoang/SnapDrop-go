package backend

import (
	"flag"
	"log"
)

type RecordComponent struct {
	Id string `json:"id"`
}

func main() {
	addr := flag.String("address", ":80", "Address to host the HTTP server on.")
	flag.Parse()

	log.Println("Listening on", *addr)
	err := serve(*addr)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getRecordComponent() (*RecordComponent, error) {

}
func serve(addr string) error {
	return nil
}
