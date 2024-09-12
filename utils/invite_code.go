package utils

import (
	"container/list"
	"errors"
	"fmt"
	"math/rand"
)

const (
	// 生成激活码的所有字符
	ActivateCodeBase = "hve8s2dzx9c7p5ik3mjufr4wy1tn6bgq"
	// 激活码第二位：分隔用字符
	ActivateCodeSuffixChar = "a"
	// 生成邀请码的所有字符
	InviteCodeBase = "HVE8S2DZX9C7P5IK3MJUFR4WY1TN6BGQ"
	// 邀请码第二位：分隔用字符
	InviteCodeSuffixChar = "A"
	// 生成优惠吗的所有字符
	PromoCodeBase = "HVE8S2DZX9C7P5IK3MJUFR4WY1TN6BGQ"
	// 优惠码第二位：分隔用字符
	PromoCodeSuffixChar = "A"
	// 进制
	BinLen = 32
)

// 生成一定长度的激活码
func IdToActivateCode(id int, size int) string {
	var activateCode string
	for id > 0 {
		index := id % BinLen
		activateCode = string(ActivateCodeBase[index]) + activateCode
		id = id / BinLen
	}

	length := LengthString(activateCode)
	if length < size-1 {
		activateCode += ActivateCodeSuffixChar
		length++
		for length < size {
			activateCode += string(ActivateCodeBase[rand.Intn(len(ActivateCodeBase))])
			length++
		}
	} else if length < size {
		activateCode += ActivateCodeSuffixChar
	} else {
		activateCode = Substring(activateCode, 0, size)
	}

	return activateCode
}

// 生成一定长度的优惠吗
func IdToPromoCode(id int, size int) string {
	var promoCode string
	for id > 0 {
		index := id % BinLen
		promoCode = string(PromoCodeBase[index]) + promoCode
		id = id / BinLen
	}

	length := LengthString(promoCode)
	if length < size-1 {
		promoCode += PromoCodeSuffixChar
		length++
		for length < size {
			promoCode += string(PromoCodeBase[rand.Intn(len(PromoCodeBase))])
			length++
		}
	} else if length < size {
		promoCode += PromoCodeSuffixChar
	} else {
		promoCode = Substring(promoCode, 0, size)
	}

	return promoCode
}

// 根据 userId 得到用户的邀请码
func UserIdToInviteCode(userId int, size int) string {
	var inviteCode string
	for userId > 0 {
		index := userId % BinLen
		inviteCode = string(InviteCodeBase[index]) + inviteCode
		userId = userId / BinLen
	}

	length := len(inviteCode)
	if length < size-1 {
		inviteCode += InviteCodeSuffixChar
		length++
		for length < size {
			inviteCode += string(InviteCodeBase[rand.Intn(len(InviteCodeBase))])
			length++
		}
	} else if length < size {
		inviteCode += InviteCodeSuffixChar
	} else {
		inviteCode = Substring(inviteCode, 0, size)
	}

	return inviteCode
}

var baseStr string = "NPRSF5G6QW1C2D3E4H7J8K9AZBLMTUVXY0"
var base []byte = []byte(baseStr)
var baseMap map[byte]int

func InitBaseMap() {
	baseMap = make(map[byte]int)
	for i, v := range base {
		baseMap[v] = i
	}
}

// 将用户id转化为6位固定邀请码
func EncodeInviteCode(n uint64) []byte {
	quotient := n
	mod := uint64(0)
	l := list.New()
	for quotient != 0 {
		//fmt.Println("---quotient:", quotient)
		mod = quotient % 34
		quotient = quotient / 34
		l.PushFront(base[int(mod)])
		//res = append(res, base[int(mod)])
		//fmt.Printf("---mod:%d, base:%s\n", mod, string(base[int(mod)]))
	}
	listLen := l.Len()

	if listLen >= 6 {
		res := make([]byte, 0, listLen)
		for i := l.Front(); i != nil; i = i.Next() {

			res = append(res, i.Value.(byte))
		}
		return res
	} else {
		res := make([]byte, 0, 6)
		for i := 0; i < 6; i++ {
			if i < 6-listLen {
				res = append(res, base[0])
			} else {
				res = append(res, l.Front().Value.(byte))
				l.Remove(l.Front())
			}

		}
		return res
	}

}

// 邀请码转化为用户id
func DecodeInviteCode(str []byte) (uint64, error) {
	InitBaseMap()
	if baseMap == nil {
		return 0, errors.New("no init base map")
	}
	if str == nil || len(str) == 0 {
		return 0, errors.New("parameter is nil or empty")
	}
	var res uint64 = 0
	var r uint64 = 0
	for i := len(str) - 1; i >= 0; i-- {
		v, ok := baseMap[str[i]]
		if !ok {
			fmt.Printf("")
			return 0, errors.New("character is not base")
		}
		var b uint64 = 1
		for j := uint64(0); j < r; j++ {
			b *= 34
		}
		res += b * uint64(v)
		r++
	}
	return res, nil
}
