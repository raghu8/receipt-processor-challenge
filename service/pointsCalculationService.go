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

	// Adding points for alphanumeric characters
	points += isAlphaNumeric(transaction.Retailer)

	// Adding points for multiple of 0.25 //should this not be considered if it has no decimal but is still a multiple of .25?
	total, err := strconv.ParseFloat(transaction.Total, 64)
	if err != nil {
		return nil, err
	}
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// Adding 5 points for every two items on the list
	if len(transaction.Items) != 0 && len(transaction.Items)%2 == 0 {
		points += (len(transaction.Items) / 2) * 5
	}
	if len(transaction.Items) != 0 && len(transaction.Items)%2 != 0 {
		points += ((len(transaction.Items) - 1) / 2) * 5
	}

	// Adding 50 points if the total is a round dollar amount with no cents
	if math.Mod(total, 1) == 0 {
		points += 50
	}

	//Discount if trimmed length of items title is a multiple of 3
	for _, item := range transaction.Items {
		trimmedItem := strings.TrimSpace(item.ShortDescription)
		if len(trimmedItem)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return nil, err
			}

			//round up to the nearest dollar
			discountedPoints := math.Ceil(price * 0.20)
			points += int(discountedPoints)
		}

	}

	//Add points if purchase day is odd
	purchaseDate := strings.Split(transaction.PurchaseDate, "-")
	day, err := strconv.Atoi(purchaseDate[2])
	if err != nil {
		return nil, err
	}
	if day%2 != 0 {
		points += 6
	}

	//Add points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime := strings.Split(transaction.PurchaseTime, ":")
	hour, errHr := strconv.Atoi(purchaseTime[0])
	minute, errMin := strconv.Atoi(purchaseTime[1])
	if errHr != nil {
		return nil, errHr
	}

	if errMin != nil {
		return nil, errMin
	}

	if (hour >= 14 && minute >= 0) && hour < 16 {
		points += 10
	}

	pointsModel.Points = points

	return &pointsModel, nil
}
