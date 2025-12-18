package postItems

import (
	"fmt"
	"strconv"
)

func priceToInt(price string) (int, error) {
	if price == "" {
		return 0, fmt.Errorf("price is required")
	}
	priceInt, err := strconv.Atoi(price)
	if err != nil {
		return 0, fmt.Errorf("invalid price format: %w", err)
	}
	if priceInt <= 0 {
		return 0, fmt.Errorf("price must be greater than 0")
	}
	return priceInt, nil
}
