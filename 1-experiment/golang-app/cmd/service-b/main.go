package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

func main() {
	var c Config
	c.load("../../config.yaml")

	h := Handler{config: c}

	s := fiber.New()
	s.Get("/api/time", h.TimeHandler)

	s.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	log.Fatalln(
		s.Listen(fmt.Sprintf(":%d", c.ServiceBPort)),
	)
}

type Handler struct {
	config Config
}

func (h *Handler) TimeHandler(c *fiber.Ctx) error {
	c.JSON(struct{ DateTime string }{DateTime: time.Now().UTC().Format("2006-01-02 15:04:05")})
	return nil
}

type Config struct {
	ServiceBPort int `yaml:"serviceBPort"`
}

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
