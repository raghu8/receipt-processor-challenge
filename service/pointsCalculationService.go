package service

import (
	"fmt"
	"math"
	"recipt-processor/model"
	"recipt-processor/repository"
	"strconv"
	"strings"
	"unicode"
)

func CreateTransaction(transaction *model.Transaction) error {
	return repository.InsertTransaction(transaction)
}

func GetTransaction(uuid string) (*model.Points, error) {
	var points = 0
	var pointsModel model.Points
	transaction, err := repository.GetTransaction(uuid)
	if err != nil {
		return nil, err
	}

	points += calculateAlphaNumericPoints(transaction)

	quarterPoints, err := calculateMultipleOfQuarterPoints(transaction)
	if err != nil {
		return nil, err
	}
	points += quarterPoints

	points += calculateItemPoints(transaction)

	roundDollarPoints, err := calculateRoundDollarPoints(transaction)
	if err != nil {
		return nil, err
	}
	points += roundDollarPoints

	discountPoints, err := calculateDiscountPoints(transaction)
	if err != nil {
		return nil, err
	}
	points += discountPoints

	dayPoints, err := calculatePurchaseDayPoints(transaction)
	if err != nil {
		return nil, err
	}
	points += dayPoints

	timePoints, err := calculatePurchaseTimePoints(transaction)
	if err != nil {
		return nil, err
	}
	points += timePoints

	pointsModel.Points = points
	return &pointsModel, nil
}

func GetAllTransactions() ([]*model.Transaction, error) {
	return repository.GetAllTransactions()
}

func DeleteTransaction(uuid string) error {
	return repository.DeleteSpecificTransaction(uuid)
}

func DeleteSpecifiedTransaction(uuid string) error {
	return repository.DeleteSpecificTransaction(uuid)
}

func isAlphaNumeric(s string) int {
	var count = 0
	for _, r := range s {
		fmt.Println(r)
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count += 1
		}
	}
	return count
}

func calculateAlphaNumericPoints(transaction *model.Transaction) int {
	return isAlphaNumeric(transaction.Retailer)
}

func calculateMultipleOfQuarterPoints(transaction *model.Transaction) (int, error) {
	total, err := strconv.ParseFloat(transaction.Total, 64)
	if err != nil {
		return 0, err
	}
	if math.Mod(total, 0.25) == 0 {
		return 25, nil
	}
	return 0, nil
}

func calculateItemPoints(transaction *model.Transaction) int {
	if len(transaction.Items) == 0 {
		return 0
	}
	return (len(transaction.Items) / 2) * 5
}

func calculateRoundDollarPoints(transaction *model.Transaction) (int, error) {
	total, err := strconv.ParseFloat(transaction.Total, 64)
	if err != nil {
		return 0, err
	}
	if math.Mod(total, 1) == 0 {
		return 50, nil
	}
	return 0, nil
}

func calculateDiscountPoints(transaction *model.Transaction) (int, error) {
	points := 0
	for _, item := range transaction.Items {
		trimmedItem := strings.TrimSpace(item.ShortDescription)
		if len(trimmedItem)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, err
			}
			discountedPoints := math.Ceil(price * 0.20)
			points += int(discountedPoints)
		}
	}
	return points, nil
}

func calculatePurchaseDayPoints(transaction *model.Transaction) (int, error) {
	purchaseDate := strings.Split(transaction.PurchaseDate, "-")
	day, err := strconv.Atoi(purchaseDate[2])
	if err != nil {
		return 0, err
	}
	if day%2 != 0 {
		return 6, nil
	}
	return 0, nil
}

func calculatePurchaseTimePoints(transaction *model.Transaction) (int, error) {
	purchaseTime := strings.Split(transaction.PurchaseTime, ":")
	hour, errHr := strconv.Atoi(purchaseTime[0])
	minute, errMin := strconv.Atoi(purchaseTime[1])
	if errHr != nil {
		return 0, errHr
	}
	if errMin != nil {
		return 0, errMin
	}
	if (hour >= 14 && minute >= 0) && hour < 16 {
		return 10, nil
	}
	return 0, nil
}
