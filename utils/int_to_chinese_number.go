package utils

import (
	"strconv"
	"strings"
)

func IntToChinessNumber(input int) string {
	digits := []string{"零", "一", "二", "三", "四", "五", "六", "七", "八", "九"}
	positions := []string{"", "十", "百", "千", "万", "十万", "百万", "千万", "亿", "十亿", "百亿", "千亿"}
	intStrArray := strings.Split(strconv.Itoa(input), "")
	result := ""
	prevIsZero := false

	if input == 0 {
		return "零"
	}

	//处理0  deal zero
	for i := 0; i < len(intStrArray); i++ {
		sn := intStrArray[i]
		if sn != "0" && !prevIsZero {
			p, _ := strconv.Atoi(sn)
			result += digits[p] + positions[len(intStrArray)-i-1]
		} else if sn == "0" {
			prevIsZero = true
		} else if sn != "0" && prevIsZero {
			p, _ := strconv.Atoi(sn)
			result += "零" + digits[p] + positions[len(intStrArray)-i-1]
		}
	}
	//处理十 deal ten
	if input < 100 {
		result = strings.ReplaceAll(result, "一十", "十")
	}
	return result
}
