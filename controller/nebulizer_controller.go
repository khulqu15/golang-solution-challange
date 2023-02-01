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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var nebulizerCollection *mongo.Collection = config.GetCollection(config.DB, "nebulizers")
var nebulizerDataCollection *mongo.Collection = config.GetCollection(config.DB, "nebulizer_data")
var validate = validator.New()

func GetAllNebulizers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var nebulizers []model.Nebulizer
	defer cancel()
	results, err := nebulizerCollection.Find(ctx, bson.M{})
	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var nebulizer model.Nebulizer
		var nebulizerData []model.NebulizerData
		if err = results.Decode(&nebulizer); err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}

		opts := options.Find().SetSort(bson.D{{"createdAt", -1}, {"_id", -1}}).SetLimit(1)
		cursor, err := nebulizerDataCollection.Find(ctx, bson.D{{"nebulizer", nebulizer.Id}}, opts)
		if err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}

		for cursor.Next(ctx) {
			var dataItem model.NebulizerData
			err := cursor.Decode(&dataItem)
			if err != nil {
				APIResponse(c, http.StatusInternalServerError, "error", err.Error())
			}
			nebulizerData = append(nebulizerData, dataItem)
		}

		nebulizer.Data = nebulizerData
		nebulizers = append(nebulizers, nebulizer)
	}
	return APIResponse(c, http.StatusOK, "success", nebulizers)
}

func CreateNebulizer(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var nebulizer model.Nebulizer
	defer cancel()
	if err := c.BodyParser(&nebulizer); err != nil {
		APIResponse(c, http.StatusBadRequest, "error", err.Error())
	}
	if validationErr := validate.Struct(&nebulizer); validationErr != nil {
		APIResponse(c, http.StatusBadRequest, "error", validationErr.Error())
	}
	newNebulizer := model.Nebulizer{
		Id:       primitive.NewObjectID(),
		Name:     nebulizer.Name,
		Location: nebulizer.Location,
	}
	result, err := nebulizerCollection.InsertOne(ctx, newNebulizer)
	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	return APIResponse(c, http.StatusOK, "success", result)
}

func CreateNebulizerData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	nebulizerId := c.Params("nebulizerID")
	var nebulizer model.Nebulizer
	var nebulizerData model.NebulizerData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(nebulizerId)

	err := nebulizerCollection.FindOne(ctx, bson.M{
		"id": objId,
	}).Decode(&nebulizer)

	if err != nil {
		APIResponse(c, http.StatusBadRequest, "error", err.Error())
	}

	if err := c.BodyParser(&nebulizerData); err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	if validationErr := validate.Struct(&nebulizerData); validationErr != nil {
		APIResponse(c, http.StatusBadRequest, "error", validationErr.Error())
	}

	newNebulizerData := model.NebulizerData{
		Id:        primitive.NewObjectID(),
		Nebulizer: objId,
		Power:     nebulizerData.Power,
		Smoke:     nebulizerData.Smoke,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := nebulizerDataCollection.InsertOne(ctx, newNebulizerData)

	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	return c.Status(http.StatusCreated).JSON(response.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": result},
	})
}

func GetANebulizer(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	nebulizerId := c.Params("nebulizerID")
	var nebulizer model.Nebulizer
	var nebulizerData []model.NebulizerData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(nebulizerId)

	err := nebulizerCollection.FindOne(ctx, bson.M{
		"id": objId,
	}).Decode(&nebulizer)

	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	opts := options.Find().SetSort(bson.D{{"createdAt", -1}, {"_id", -1}})
	cursor, err := nebulizerDataCollection.Find(ctx, bson.D{{"nebulizer", nebulizer.Id}}, opts)
	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	for cursor.Next(ctx) {
		var dataItem model.NebulizerData
		err := cursor.Decode(&dataItem)
		if err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}
		nebulizerData = append(nebulizerData, dataItem)
	}

	nebulizer.Data = nebulizerData

	return APIResponse(c, http.StatusOK, "success", nebulizer)
}

func EditNebulizer(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	nebulizerId := c.Params("nebulizerId")
	var nebulizer model.Nebulizer
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(nebulizerId)
	if err := c.BodyParser(&nebulizer); err != nil {
		APIResponse(c, http.StatusBadRequest, "error", err.Error())
	}

	if validationErr := validate.Struct(&nebulizer); validationErr != nil {
		APIResponse(c, http.StatusBadRequest, "error", validationErr.Error())
	}

	update := bson.M{
		"name":     nebulizer.Name,
		"Location": nebulizer.Location,
	}

	result, err := nebulizerCollection.UpdateOne(ctx, bson.M{
		"id": objId,
	}, bson.M{"$set": update})

	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	var updateNebulizer model.Nebulizer
	if result.MatchedCount == 1 {
		err := nebulizerCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updateNebulizer)
		if err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}
	}

	return APIResponse(c, http.StatusOK, "success", updateNebulizer)
}

func EditNebulizerData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	nebulizerId := c.Params("nebulizerId")
	dataId := c.Params("dataId")
	var nebulizer model.Nebulizer
	var nebulizerData model.NebulizerData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(nebulizerId)
	err := nebulizerCollection.FindOne(ctx, bson.M{
		"id": objId,
	}).Decode(&nebulizer)
	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	objDataId, _ := primitive.ObjectIDFromHex(dataId)
	if err := c.BodyParser(&nebulizerData); err != nil {
		APIResponse(c, http.StatusBadRequest, "error", err.Error())
	}
	if validationErr := validate.Struct(&nebulizerData); validationErr != nil {
		APIResponse(c, http.StatusBadRequest, "error", validationErr.Error())
	}

	update := bson.M{
		"power": nebulizerData.Power,
		"smoke": nebulizerData.Smoke,
	}

	result, err := nebulizerDataCollection.UpdateOne(ctx, bson.M{
		"id": objDataId,
	}, bson.M{"$set": update})

	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	var updateData model.NebulizerData
	if result.MatchedCount == 1 {
		err := nebulizerDataCollection.FindOne(ctx, bson.M{"id": objDataId}).Decode(&updateData)
		if err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}
	}
	return APIResponse(c, http.StatusOK, "success", updateData)
}

func DeleteNebulizer(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	nebulizerId := c.Params("nebulizerId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(nebulizerId)

	nebulizerResult := nebulizerCollection.FindOne(ctx, bson.M{"id": objId})
	if nebulizerResult.Err() != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", nebulizerResult.Err().Error())
	}

	var nebulizer model.Nebulizer
	if err := nebulizerResult.Decode(&nebulizer); err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	if _, err := nebulizerDataCollection.DeleteMany(ctx, bson.M{"nebulizer": nebulizer.Id}); err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	result, err := nebulizerCollection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	if result.DeletedCount < 1 {
		return APIResponse(c, http.StatusNotFound, "not found", "Nebulizer with specified ID not found!")
	}

	return APIResponse(c, http.StatusOK, "success", "Nebulizer successfully deleted!")
}

func DeleteNebulizerData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	nebulizerId := c.Params("nebulizerId")
	dataId := c.Params("dataId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(nebulizerId)

	nebulizerResult := nebulizerCollection.FindOne(ctx, bson.M{"id": objId})
	if nebulizerResult.Err() != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", nebulizerResult.Err().Error())
	}

	var nebulizer model.Nebulizer
	if err := nebulizerResult.Decode(&nebulizer); err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	objDataId, _ := primitive.ObjectIDFromHex(dataId)
	dataResult := nebulizerDataCollection.FindOne(ctx, bson.M{"id": objDataId})
	if dataResult.Err() != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", dataResult.Err().Error())
	}

	deleteResult, err := nebulizerDataCollection.DeleteOne(ctx, bson.M{"id": objDataId})

	if err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	if deleteResult.DeletedCount < 1 {
		return APIResponse(c, http.StatusNotFound, "not found", "Nebulizer data with specified ID not found!")
	}

	return APIResponse(c, http.StatusOK, "success", "Nebulizer Data successfully deleted!")

}
