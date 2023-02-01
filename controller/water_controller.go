package controller

import (
	"context"
	"net/http"
	"solution-challange/config"
	"solution-challange/model"
	"solution-challange/response"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var waterCollection *mongo.Collection = config.GetCollection(config.DB, "waters")
var waterDataCollection *mongo.Collection = config.GetCollection(config.DB, "water_data")

func GetAllWaters(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var waters []model.Water
	defer cancel()
	results, err := waterCollection.Find(ctx, bson.M{})
	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	defer results.Close(ctx)
	for results.Next(ctx) {
		var water model.Water
		var waterData []model.WaterData
		if err = results.Decode(&water); err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}

		opts := options.Find().SetSort(bson.D{{"createdAt", -1}, {"_id", -1}}).SetLimit(1)
		cursor, err := waterDataCollection.Find(ctx, bson.D{{"water", water.Id}}, opts)
		if err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}

		for cursor.Next(ctx) {
			var dataItem model.WaterData
			err := cursor.Decode(&dataItem)
			if err != nil {
				APIResponse(c, http.StatusInternalServerError, "error", err.Error())
			}
			waterData = append(waterData, dataItem)
		}

		water.Data = waterData
		waters = append(waters, water)
	}
	return APIResponse(c, http.StatusOK, "success", waters)
}

func CreateWater(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var water model.Water
	defer cancel()
	if err := c.BodyParser(&water); err != nil {
		APIResponse(c, http.StatusBadRequest, "error", err.Error())
	}
	if validationErr := validate.Struct(&water); validationErr != nil {
		APIResponse(c, http.StatusBadRequest, "error", validationErr.Error())
	}
	newWater := model.Water{
		Id:       primitive.NewObjectID(),
		Name:     water.Name,
		Location: water.Location,
	}
	result, err := waterCollection.InsertOne(ctx, newWater)
	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	return APIResponse(c, http.StatusOK, "success", result)
}

func CreateWaterData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	waterId := c.Params("waterID")
	var water model.Water
	var waterData model.WaterData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(waterId)

	err := waterCollection.FindOne(ctx, bson.M{
		"id": objId,
	}).Decode(&water)

	if err != nil {
		APIResponse(c, http.StatusBadRequest, "error", err.Error())
	}

	if err := c.BodyParser(&waterData); err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	if validationErr := validate.Struct(&waterData); validationErr != nil {
		APIResponse(c, http.StatusBadRequest, "error", validationErr.Error())
	}

	newWaterData := model.WaterData{
		Id:          primitive.NewObjectID(),
		Water:       objId,
		Acidity:     waterData.Acidity,
		Salt:        waterData.Salt,
		Oxygen:      waterData.Oxygen,
		Temperature: waterData.Temperature,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := waterDataCollection.InsertOne(ctx, newWaterData)

	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	return c.Status(http.StatusCreated).JSON(response.UserResponse{
		Status:  http.StatusCreated,
		Message: "success",
		Data:    &fiber.Map{"data": result},
	})
}

func GetAWater(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	waterId := c.Params("waterID")
	var water model.Water
	var waterData []model.WaterData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(waterId)

	err := waterCollection.FindOne(ctx, bson.M{
		"id": objId,
	}).Decode(&water)

	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	opts := options.Find().SetSort(bson.D{{"createdAt", -1}, {"_id", -1}})
	cursor, err := waterDataCollection.Find(ctx, bson.D{{"water", water.Id}}, opts)
	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	for cursor.Next(ctx) {
		var dataItem model.WaterData
		err := cursor.Decode(&dataItem)
		if err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}
		waterData = append(waterData, dataItem)
	}

	water.Data = waterData

	return APIResponse(c, http.StatusOK, "success", water)
}

func EditWater(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	waterId := c.Params("waterId")
	var water model.Water
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(waterId)
	if err := c.BodyParser(&water); err != nil {
		APIResponse(c, http.StatusBadRequest, "error", err.Error())
	}

	if validationErr := validate.Struct(&water); validationErr != nil {
		APIResponse(c, http.StatusBadRequest, "error", validationErr.Error())
	}

	update := bson.M{
		"name":     water.Name,
		"Location": water.Location,
	}

	result, err := waterCollection.UpdateOne(ctx, bson.M{
		"id": objId,
	}, bson.M{"$set": update})

	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	var updateWater model.Water
	if result.MatchedCount == 1 {
		err := waterCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updateWater)
		if err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}
	}

	return APIResponse(c, http.StatusOK, "success", updateWater)
}

func EditWaterData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	waterId := c.Params("waterId")
	dataId := c.Params("dataId")
	var water model.Water
	var waterData model.WaterData
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(waterId)
	err := waterCollection.FindOne(ctx, bson.M{
		"id": objId,
	}).Decode(&water)
	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	objDataId, _ := primitive.ObjectIDFromHex(dataId)
	if err := c.BodyParser(&waterData); err != nil {
		APIResponse(c, http.StatusBadRequest, "error", err.Error())
	}
	if validationErr := validate.Struct(&waterData); validationErr != nil {
		APIResponse(c, http.StatusBadRequest, "error", validationErr.Error())
	}

	update := bson.M{
		"acidity":     waterData.Acidity,
		"salt":        waterData.Salt,
		"oxygen":      waterData.Oxygen,
		"temperature": waterData.Temperature,
	}

	result, err := waterDataCollection.UpdateOne(ctx, bson.M{
		"id": objDataId,
	}, bson.M{"$set": update})

	if err != nil {
		APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}
	var updateData model.WaterData
	if result.MatchedCount == 1 {
		err := waterDataCollection.FindOne(ctx, bson.M{"id": objDataId}).Decode(&updateData)
		if err != nil {
			APIResponse(c, http.StatusInternalServerError, "error", err.Error())
		}
	}
	return APIResponse(c, http.StatusOK, "success", updateData)
}

func DeleteWater(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	waterId := c.Params("waterId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(waterId)

	waterResult := waterCollection.FindOne(ctx, bson.M{"id": objId})
	if waterResult.Err() != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", waterResult.Err().Error())
	}

	var water model.Water
	if err := waterResult.Decode(&water); err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	if _, err := waterDataCollection.DeleteMany(ctx, bson.M{"water": water.Id}); err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	result, err := waterCollection.DeleteOne(ctx, bson.M{"id": objId})

	if err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	if result.DeletedCount < 1 {
		return APIResponse(c, http.StatusNotFound, "not found", "Water with specified ID not found!")
	}

	return APIResponse(c, http.StatusOK, "success", "Water successfully deleted!")
}

func DeleteWaterData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	waterId := c.Params("waterId")
	dataId := c.Params("dataId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(waterId)

	waterResult := waterCollection.FindOne(ctx, bson.M{"id": objId})
	if waterResult.Err() != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", waterResult.Err().Error())
	}

	var water model.Water
	if err := waterResult.Decode(&water); err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	objDataId, _ := primitive.ObjectIDFromHex(dataId)
	dataResult := waterDataCollection.FindOne(ctx, bson.M{"id": objDataId})
	if dataResult.Err() != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", dataResult.Err().Error())
	}

	deleteResult, err := waterDataCollection.DeleteOne(ctx, bson.M{"id": objDataId})

	if err != nil {
		return APIResponse(c, http.StatusInternalServerError, "error", err.Error())
	}

	if deleteResult.DeletedCount < 1 {
		return APIResponse(c, http.StatusNotFound, "not found", "Water data with specified ID not found!")
	}

	return APIResponse(c, http.StatusOK, "success", "Water Data successfully deleted!")

}
