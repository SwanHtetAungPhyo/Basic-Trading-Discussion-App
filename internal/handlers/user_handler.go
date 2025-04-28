package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserHandlers interface {
	LoginHandler(ctx *fiber.Ctx) error
	RegisterHandler(ctx *fiber.Ctx) error
	LogoutHandler(ctx *fiber.Ctx) error
}

var _ UserHandlers = (*UserHanding)(nil)

type UserHanding struct {
	log *logrus.Logger
}

func (u UserHanding) LoginHandler(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHanding) RegisterHandler(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}

func (u UserHanding) LogoutHandler(ctx *fiber.Ctx) error {
	//TODO implement me
	panic("implement me")
}
