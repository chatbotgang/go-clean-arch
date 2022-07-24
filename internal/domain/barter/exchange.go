package barter

func ExchangeGoods(requestGood Good, targetGood Good) []Good {
	requestGood.OwnerID, targetGood.OwnerID = targetGood.OwnerID, requestGood.OwnerID
	return []Good{requestGood, targetGood}
}
