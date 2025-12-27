package model

// AllCards 返回从 3 到大王的所有牌（不带花色）
func AllCards() []string {
	values := []string{"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"}
	cards := []string{}
	for _, v := range values {
		cards = append(cards, "CARD_"+v)
	}
	// 小王和大王
	cards = append(cards, "CARD_BLACK") // 小王
	cards = append(cards, "CARD_RED")   // 大王
	return cards
}

// ToRealCard 将牌字符串转换成整数值，CARD_3 = 1
func ToRealCard(card string) int {
	switch card {
	case "CARD_3":
		return 1
	case "CARD_4":
		return 2
	case "CARD_5":
		return 3
	case "CARD_6":
		return 4
	case "CARD_7":
		return 5
	case "CARD_8":
		return 6
	case "CARD_9":
		return 7
	case "CARD_10":
		return 8
	case "CARD_J":
		return 9
	case "CARD_Q":
		return 10
	case "CARD_K":
		return 11
	case "CARD_A":
		return 12
	case "CARD_2":
		return 13
	case "CARD_BLACK":
		return 14 // 小王
	case "CARD_RED":
		return 15 // 大王
	default:
		return 0
	}
}
