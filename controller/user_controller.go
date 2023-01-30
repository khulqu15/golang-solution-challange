package controller

import (
	"context"
	"net/http"
	"solution-challange/config"
	"solution-challange/model"
	"solution-challange/response"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "users")
var validate = validator.New()

func GetAllUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []model.User
	defer cancel()

	results, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleUser model.User
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"data": err.Error()},
			})
		}
		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(response.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": users},
	})
}

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user model.User
	defer cancel()

	// Validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	// Use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(response.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validationErr.Error()},
		})
	}

	newUser := model.User{
		Id:       primitive.NewObjectID(),
		Name:     user.Name,
		Location: user.Location,
		Title:    user.Title,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusCreated).JSON(response.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": result},
	})
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userID")
	var user model.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{
		"id": objId,
	}).Decode(&user)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	return c.Status(http.StatusOK).JSON(response.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": user},
	})
}

func EditUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	var user model.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(response.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(response.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    &fiber.Map{"data": validationErr.Error()},
		})
	}

	update := bson.M{
		"name":     user.Name,
		"location": user.Location,
		"title":    user.Title,
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{
		"id": objId,
	}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	var updatedUser model.User
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    &fiber.Map{"data": err.Error()},
			})
		}
	}

	return c.Status(http.StatusOK).JSON(response.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": updatedUser},
	})
}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId := c.Params("userId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(response.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    &fiber.Map{"data": err.Error()},
		})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(response.UserResponse{
			Status:  http.StatusNotFound,
			Message: "error",
			Data:    &fiber.Map{"data": "User with specified ID not found!"},
		})
	}

	return c.Status(http.StatusOK).JSON(response.UserResponse{
		Status:  http.StatusOK,
		Message: "success",
		Data:    &fiber.Map{"data": "User successfully deleted!"},
	})
}
