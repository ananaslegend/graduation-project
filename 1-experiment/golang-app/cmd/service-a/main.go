package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber"
	"gopkg.in/yaml.v3"
)

func main() {
	var c Config
	c.load("../../config.yaml")

	h := Handler{config: c}

	s := fiber.New()
	s.Get("/api/time", h.timeHandler)
	log.Fatalln(
		s.Listen(c.ServiceAPort),
	)
}

type Handler struct{
	config Config
}

func (h *Handler) timeHandler(c *fiber.Ctx){
	ep := fmt.Sprintf("%s/api/time", h.config.ServiceBUrl)
	r, err := http.Get(ep)
	if err != nil {
		c.JSON(fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	
	msg := struct{DateTime string}{}

	err = json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		c.JSON(fmt.Sprintf("Error: %s", err.Error()))
		return
	}
	defer r.Body.Close()

	c.JSON(msg)
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