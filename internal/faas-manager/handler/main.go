package handler

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/services"
	"github.com/tiago123456789/my-own-faas-golang-platform/internal/faas-manager/types"
)

var validate *validator.Validate

type FunctionHandler struct {
	functionService services.FunctionService
}

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func NewFunctionHandler(
	functionService services.FunctionService,
) *FunctionHandler {
	return &FunctionHandler{
		functionService: functionService,
	}
}

func (f *FunctionHandler) FindAll(c *fiber.Ctx) error {
	return c.JSON(f.functionService.FindAll())
}

func (f *FunctionHandler) Deploy(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("No file found")
	}

	savePath := filepath.Join("uploads", fmt.Sprintf("%s_%s", uuid.NewString(), file.Filename))
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to save file")
	}

	var newFunction types.NewFunction

	if err := c.BodyParser(&newFunction); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if err := validate.Struct(newFunction); err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("field %s: wanted %s %s, got `%s`", err.Field(), err.Tag(), err.Param(), err.Value()))
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
	}

	err = f.functionService.Deploy(newFunction, savePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": []string{
				"Failed to publish message on queue",
			},
		})
	}

	return c.SendStatus(201)
}
