package main

import (
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/ollama/ollama/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/joho/godotenv/autoload"
)


type Model struct {
	Name       string
	ParamCount string
	Storage    string
}

var (
	ollamaApi, _ = url.Parse(os.Getenv("OLLAMA_BASE_URL"))
	models       = []*Model{
		{Name: "codellama:latest", ParamCount: "7B", Storage: "3.8GB"},
		{Name: "llama2:latest", ParamCount: "7B", Storage: "3.8GB"},
		{Name: "llama3:latest", ParamCount: "8B", Storage: "4.7GB"},
		{Name: "tinyllama:latest", ParamCount: "1.1B", Storage: "640MB"},
	}
	questions = []string{
		"What is the capital of France?",
		"Who is the CEO of Apple?",
		"What is the meaning of life?",
		"What is the best way to get to work?",
		"How do I make a cup of coffee?",
		"What is the capital of Australia?",
		"Can you tell me about a famous painting by Leonardo da Vinci?",
		"Who is the CEO of Tesla?",
		"Can you give me a recipe for a delicious vegetarian dish?",
		"What is the capital of Canada?",
		"How do I properly format my phone number for international calls?",
		"What is the best way to stay healthy during the pandemic?",
	}
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Kitchen})
}

func main() {
	var (
		ollama         = api.NewClient(ollamaApi, http.DefaultClient)
		c              = New(ollama)
		runs           = make(Runs, len(models))
		done, shutdown = make(chan interface{}), make(chan os.Signal, 1)
	)

	log.Info().Timestamp().Msg("Start")

	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		<-shutdown
		log.Info().Timestamp().Msg("End")

		done <- nil
	}()

	color.New(color.BgBlack).Add(color.FgWhite).Println("--- Round 1 ---")

	for i, model := range models {
		run := c.TestModel(model)
		run.PrintResult()
		runs[i] = run
	}

	color.New(color.BgBlack).Add(color.FgWhite).Println("--- Round 2 ---")

	for i, model := range models {
		run := c.TestModel(model)
		run.PrintResult()
		// Average the two runs together
		runs[i].JoinAndAverage(run)
	}

	shutdown <- nil
	<-done
os.Exit(0)
}
