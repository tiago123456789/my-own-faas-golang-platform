package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/proxy/services"
	"github.com/valyala/fasthttp"
)

type FunctionHandler struct {
	functionExecutor services.FunctionExecutorService
}

func NewFunctionHandler(
	functionExecutor services.FunctionExecutorService,
) *FunctionHandler {
	return &FunctionHandler{
		functionExecutor: functionExecutor,
	}
}

func (f *FunctionHandler) Execute(c *fiber.Ctx) error {
	function := c.Params("function")
	path := c.Params("*1")

	fmt.Println(function, path)
	err := f.functionExecutor.Run(function)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := proxy.DoRedirects(c, fmt.Sprintf("http://%s:3000/%s", function, path), 3, &fasthttp.Client{
		NoDefaultUserAgentHeader: false,
		DisablePathNormalizing:   true,
	}); err != nil {
		return err
	}

	return nil
}
