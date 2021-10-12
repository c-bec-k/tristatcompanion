package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c-bec-k/tristatcompanion/internal/data"
)

type config struct {
	port   int
	env    string
	token  string
	api    int
	pubkey string
}

type application struct {
	config   config
	logger   *log.Logger
	cmdCache map[data.Snowflake]func(http.ResponseWriter, map[string]interface{}, data.Interaction)
}

var (
	UserAgent = "DiscordBot (https://github.com/c-bec-k/tsc, v0.1.0)"
)

//var version = "1.0.0"

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "addr", 8080, "HTTP network address")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.token, "token", "", "your bot token")
	flag.IntVar(&cfg.api, "api", 9, "default API version")
	flag.StringVar(&cfg.pubkey, "pubkey", "", "Your bot's public key")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	bot := application{
		config: cfg,
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", bot.home)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", bot.config.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	if cfg.token == "" {
		log.Fatal("You need to provide a bot token to connect to the API!")
	}

	bot.cmdCache = map[data.Snowflake]func(http.ResponseWriter, map[string]interface{}, data.Interaction){
		895024426087747645: bot.ReplyRoll,
		897500010848067614: bot.ReplyAbout,
	}

	fmt.Printf("Running app on %v with version number %v\n", cfg.port, cfg.api)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error: %v\n", err)
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	verified := VerifyBot(r, app.config.pubkey)
	if !verified {
		http.Error(w, "signature mismatch", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "invalid request", http.StatusUnauthorized)
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer r.Body.Close()

	// fmt.Printf("Incoming Body: %s\n", b)
	apireq := data.Interaction{}

	err = json.Unmarshal(b, &apireq)
	if err != nil {
		fmt.Println("Verification failed")
		http.Error(w, err.Error(), 500)
		return
	}
	// fmt.Printf("APIreq struct: %+v\n", apireq)
	// fmt.Println("————————————————————")

	if apireq.Type == 1 {
		// fmt.Println("Verification successful! Sending JSON reply")
		jsonRes := `{"type": 1}`
		tokenHeader := fmt.Sprintf("bot %v", app.config.token)
		w.Header().Set("Authorization", tokenHeader)
		w.Header().Set("User-Agent", UserAgent)
		w.Write([]byte(jsonRes))
	}

	commandID := apireq.Data.ID

	opts := map[string]interface{}{}
	for _, v := range apireq.Data.Options {
		opts[v.Name] = v.Value
	}

	if fn, ok := app.cmdCache[commandID]; ok {
		fn(w, opts, apireq)
	}
}
