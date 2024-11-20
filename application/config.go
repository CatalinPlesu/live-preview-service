package application

import (
      "os"
      "strconv"
)

type Config struct {
      ServerPort uint16
      RabitMQURL string
}

func LoadConfig() Config {
      cfg := Config{
	      ServerPort: 3333,
	      RabitMQURL: "amqp://guest:guest@localhost:5672/",
      }

      if rabitMQURL, exists := os.LookupEnv("RABITMQ_URL"); exists {
	      cfg.RabitMQURL = rabitMQURL
      }

      if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
	      if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
		      cfg.ServerPort = uint16(port)
	      }
      }

      return cfg
}
