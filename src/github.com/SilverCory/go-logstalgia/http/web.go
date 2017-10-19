package http

import (
	"fmt"
	"github.com/SilverCory/go-logstalgia/config"
	"github.com/SilverCory/go-logstalgia/http/socket"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"sync"
)

type Server struct {
	rootRouter  *mux.Router
	templater   *template.Template
	config      *config.LogstalgiaConfig
	clientsLock *sync.Mutex
	socket      *socket.Handler
}

func New(conf *config.LogstalgiaConfig) (s *Server) {

	s = &Server{
		config:      conf,
		rootRouter:  mux.NewRouter(),
		templater:   template.Must(template.New("index.html").ParseFiles(conf.TemplateFileDirectory + "index.html")),
		clientsLock: &sync.Mutex{},
		socket:      socket.New(),
	}

	s.rootRouter.HandleFunc("/ws", s.socket.UpgradeWebsocket)
	s.rootRouter.HandleFunc("/", s.handleIndex)
	s.rootRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	return

}

type PageData struct {
	Speed     int
	Framerate int
	Colours   bool
	Time      bool
	Summarise bool
}

type LogEntry struct {
	Time   string `json:"time"`
	Path   string `json:"path"`
	Size   int    `json:"size"`
	IP     string `json:"ip"`
	Method string `json:"method"`
	Result int    `json:"result"`
}

func (s *Server) Broadcast(p *LogEntry) {
	s.socket.BroadcastJSON(p)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {

	err := s.templater.Execute(w, s.config.PageConfig)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) Open() {
	http.Handle("/", s.rootRouter)
	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Println("Fatal err:", err)
	}
}
