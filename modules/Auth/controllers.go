package auth

import (
	validation "server/modules/Validation"

	"github.com/gofiber/fiber/v2"
)

func RegisterHandler(c *fiber.Ctx) error {
	var RegisterRequest RegisterRequest
	if err := c.BodyParser(&RegisterRequest); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid request body"})
	}

	switch RegisterRequest.Role {
	case "user":
		return RegisterUserHandler(c)
	case "admin":
		return RegisterAdminHandler(c)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Role"})
	}
}

func LoginHandler(c *fiber.Ctx) error {
	var input LoginRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if err := validation.ValidateInputs(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := Login(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
	})

}

func GetUserHandler(c *fiber.Ctx) error {
	users, err := GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Could not fetch users"})
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

func UpdateUserHandler(c *fiber.Ctx) error {
	id := c.Query("id")
	var Updateuserdata map[string]interface{}

	err := c.BodyParser(&Updateuserdata)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse json"})
	}

	err = UpdateUser(id, Updateuserdata)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Faild to update user"})
	}
	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

func DeleteUserHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	err := DeleteUser(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete a user"})
	}

	return c.JSON(fiber.Map{"message": "User Deleted Successfully"})
}

func GetAdminHandler(c *fiber.Ctx) error {

	admins, err := GetAllAdmins()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to fetch admin data"})
	}
	return c.Status(fiber.StatusOK).JSON(admins)
}

func UpdateAdminHandler(c *fiber.Ctx) error {

	id := c.Query("id")
	var Updateadmindata map[string]interface{}

	if err := c.BodyParser(&Updateadmindata); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse JSON"})
	}

	if err := UpdateAdmin(id, Updateadmindata); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to update admin"})
	}
	return c.JSON(fiber.Map{"message": "Admin updated successfully"})
}

func DeleteAdminHandler(c *fiber.Ctx) error {
	id := c.Query("id")

	err := DeleteAdmin(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to deleted admin"})
	}
	return c.JSON(fiber.Map{"message": "Admin Deleted successfully"})
}
