package utils

import "testing"

func TestUserIdToInviteCode(t *testing.T) {
	number := int64(456789101111)
	s := EncodeInviteCode(number)
	info := UserIdToInviteCode(int(number), 6)
	num := InviteCodeToUserId(info)
	if decode, err := DecodeInviteCode(s); err == nil {
		if int64(decode) == number {
			t.Log("success", number, "encode", s, decode, "-----", number, info, num)

		} else {
			t.Fatal(s, decode)
		}

	}
}

func TestEncodeInviteCode(t *testing.T) {
	number := int64(17001000)
	s := EncodeInviteCode(number)
	info := UserIdToInviteCode(int(number), 6)
	num := InviteCodeToUserId(info)
	if decode, err := DecodeInviteCode(s); err == nil {
		if int64(decode) == number {
			t.Log("success", number, "encode", s, decode, ":", number, info, num)

		} else {
			t.Fatal(s, decode)
		}
	}
}
