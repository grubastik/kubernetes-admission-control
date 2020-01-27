package root

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//Index provide information about service
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is sample of kubernetes admission webhook server.\n"))
	w.Write([]byte("The only thing it's do is add label.\n\n"))
	w.Write([]byte("Author is Oleksii Herman.\n"))
}
