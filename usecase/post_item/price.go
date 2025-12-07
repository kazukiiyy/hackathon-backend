package post_item

import (
	"fmt"
	"strconv"
)

func priceToInt(price string) int {
	priceInt, err := strconv.Atoi(price)
	if err != nil {
		fmt.Errorf("invalid price format: %w", err)
	}
	return priceInt
}
