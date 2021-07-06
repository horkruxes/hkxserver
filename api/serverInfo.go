package api

import (
	"github.com/ewenquim/horkruxes/service"
	"github.com/gofiber/fiber/v2"
)

type PublicServerInfo struct {
}

func GetServerInfo(s service.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		s.ServerConfig.Port = -1
		return c.Status(200).JSON(s.ServerConfig)
	}
}
