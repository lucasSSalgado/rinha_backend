package controller

import (
	"rinha/models"
	"rinha/service"
	"rinha/util"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type ClientController struct {
	serv *service.ClientService
}

func NewControl(serv *service.ClientService) *ClientController {
	return &ClientController{
		serv: serv,
	}
}

func (c *ClientController) InitRoutes() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	},
	)
	cli := app.Group("/clientes")

	cli.Post("/:id/transacoes", c.transacoes)
	cli.Get("/:id/extrato", c.extrato)

	app.Listen(":8080")
}

func (c *ClientController) transacoes(ctx *fiber.Ctx) error {
	imutableId := utils.CopyString(ctx.Params("id"))

	if err := c.serv.FindClientById(imutableId); err != nil {
		return ctx.SendStatus(404)
	}

	var transaction models.TransacaoRequDto
	if err := ctx.BodyParser(&transaction); err != nil {
		return ctx.SendStatus(422)
	}

	if err := util.CheckFields(&transaction); err != nil {
		return ctx.SendStatus(422)
	}

	limite, saldo, err := c.serv.LidarComTransacao(&transaction, imutableId)
	if err != nil {
		return ctx.SendStatus(422)
	}

	return ctx.JSON(
		fiber.Map{
			"limite": limite,
			"saldo":  saldo,
		},
	)
}

func (c *ClientController) extrato(ctx *fiber.Ctx) error {
	imutableId := utils.CopyString(ctx.Params("id"))
	if err := c.serv.FindClientById(imutableId); err != nil {
		return ctx.SendStatus(404)
	}

	historico, err := c.serv.GetHistorico(imutableId)
	if err != nil {
		return ctx.SendStatus(412)
	}

	return ctx.JSON(historico)
}
