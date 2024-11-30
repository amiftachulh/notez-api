package handler

import (
	"log"

	"github.com/amiftachulh/notez-api/model"
	"github.com/amiftachulh/notez-api/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UpdateNoteMemberRole(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	id := c.Locals("params").(*model.NoteMemberParams).ID

	isOwner, err := service.CheckIsNoteOwner(id, auth.ID)
	if err != nil {
		log.Println("Error checking note owner:", err)
		return fiber.ErrInternalServerError
	}
	if !isOwner {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	memberID := c.Params("memberID")
	mID, err := uuid.Parse(memberID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{
			Message: "Invalid member ID.",
		})
	}

	body := c.Locals("body").(*model.UpdateNoteMemberRole)

	result, err := service.UpdateNoteMemberRole(id, mID, body.Role)
	if err != nil {
		log.Println("Error updating note member role:", err)
		return fiber.ErrInternalServerError
	}
	if !result {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note or member not found.",
		})
	}

	return c.JSON(model.Response{
		Message: "Note member role updated.",
	})
}

func RemoveNoteMember(c *fiber.Ctx) error {
	auth := c.Locals("auth").(model.AuthUser)
	params := c.Locals("params").(*model.NoteMemberParams)

	isOwner, err := service.CheckIsNoteOwner(params.ID, auth.ID)
	if err != nil {
		log.Println("Error checking note owner:", err)
		return fiber.ErrInternalServerError
	}
	if !isOwner {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note not found.",
		})
	}

	result, err := service.RemoveNoteMember(params.ID, params.MemberID)
	if err != nil {
		log.Println("Error removing note member:", err)
		return fiber.ErrInternalServerError
	}
	if !result {
		return c.Status(fiber.StatusNotFound).JSON(model.Response{
			Message: "Note or member not found.",
		})
	}

	return c.JSON(model.Response{
		Message: "Note member removed.",
	})
}
