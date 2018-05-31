package goblet

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/extrame/go-random"
	"github.com/extrame/goblet/config"
)

func (s *Server) wrapError(w http.ResponseWriter, err interface{}, withStack bool) {
	w.WriteHeader(500)
	if s.Env() == config.ProductEnv {
		errKey := gorandom.RandomNumeric(10)
		log.Printf("%T,%v(ERROR key %s)\n", err, err, errKey)
		if withStack {
			log.Print(string(debug.Stack()))
		}
		html := fmt.Sprintf(`<body><h4>Internal Error(%s)</h4><br/>The Random Key is %s</body>`, errKey, errKey)
		w.Write([]byte(html))
	} else {
		w.Write([]byte("<body><h4>"))
		w.Write([]byte(fmt.Sprintf("%T,%v", err, err)))
		w.Write([]byte("</h4>"))
		if withStack {
			w.Write([]byte("<pre>"))
			w.Write([]byte(debug.Stack()))
			w.Write([]byte("</pre>"))
		}
		w.Write([]byte("</body>"))
	}
}

var Interrupted = errors.New("interrupted error")
