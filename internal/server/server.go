package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pixellini/go-audiobook/internal/app"
	"github.com/pixellini/go-audiobook/internal/flags"
	"github.com/pixellini/go-audiobook/internal/progress"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	flags *flags.Flags
}

func New() *Server {
	return &Server{
		flags: flags.New(),
	}
}

func (s *Server) StreamAudiobookCreation(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	app, err := app.NewWithFlags(s.flags)
	if err != nil {
		c.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": err.Error(),
		})
		return
	}

	progressChannel := make(chan progress.Progress, 64)
	app.SetProgressChannel(progressChannel)

	go func() {
		for prog := range progressChannel {
			_ = c.WriteJSON(map[string]interface{}{
				"type":       "progress",
				"step":       prog.Step,
				"percentage": prog.Percentage,
				"message":    prog.Message,
				"chapter":    prog.Chapter,
			})
		}
	}()

	err = app.Run()
	if err != nil {
		c.WriteJSON(map[string]interface{}{
			"type":    "error",
			"message": err.Error(),
		})
		return
	}

	c.WriteJSON(map[string]interface{}{
		"type":    "complete",
		"message": "Audiobook generation finished",
	})
}

func (s *Server) Start(port string) error {
	http.HandleFunc("/generate", s.StreamAudiobookCreation)
	// http.Handle("/", http.FileServer(http.Dir("./web/dist/")))

	log.Printf("Server started and listening on port %s", port)
	return http.ListenAndServe(port, nil)
}
