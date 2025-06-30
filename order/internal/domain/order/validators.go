package order

func validateAddress(address string) bool {
	return address != ""
}

func validateItem(item Item) bool {
	if item.Price.Sign() <= 0 {
		return false
	}
	if item.Count <= 0 {
		return false
	}
	return true
}

func validateItems(items []Item) bool {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if !validateItem(item) {
			return false
		}
	}
	return true
}
