package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

func main() {
	var c Config
	c.load("./config.yaml")

	h := Handler{config: c}

	s := fiber.New()
	s.Get("/api/time", h.timeHandler)

	s.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	log.Fatalln(
		s.Listen(fmt.Sprintf(":%d", c.ServiceAPort)),
	)
}

type Handler struct {
	config Config
}

func (h *Handler) timeHandler(c *fiber.Ctx) error {
	ep := fmt.Sprintf("%s/api/time", h.config.ServiceBUrl)
	r, err := http.Get(ep)
	if err != nil {
		return err
	}

	msg := struct{ DateTime string }{}

	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	c.JSON(msg)
	return nil
}

type Config struct {
	ServiceAPort int `yaml:"serviceAPort"`

	ServiceBUrl string `yaml:"serviceBUrl"`
}

// load .yaml config file
func (c *Config) load(path string) {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("cant read config file: %v", err)
	}

	err = yaml.Unmarshal(f, c)
	if err != nil {
		log.Fatalf("cant unmarshal config file: %v", err)
	}
}
