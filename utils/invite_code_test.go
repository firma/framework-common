package utils

import "testing"

func TestUserIdToInviteCode(t *testing.T) {
	number := int64(10000)
	s := EncodeInviteCode(number)
	inviteCode := UserIdToInviteCode(int(number), 8)
	num := InviteCodeToUserId(inviteCode)
	if decode, err := DecodeInviteCode(s); err == nil {
		if int64(decode) == number {
			t.Log("success", number, "encode", string(s), inviteCode, "-----", decode, num)
		} else {
			t.Fatal(s, decode)
		}

	}
}

func TestEncodeInviteCode(t *testing.T) {
	number := int64(17001000)
	encodeCode := EncodeInviteCode(number)

	userInviteCode := UserIdToInviteCode(int(number), 10)
	userInviteDecode := InviteCodeToUserId(userInviteCode)

	if decode, err := DecodeInviteCode(encodeCode); err == nil {
		if int64(decode) == number {
			t.Log("success", number, "encode", string(encodeCode), userInviteCode, ":", number, userInviteDecode)

		} else {
			t.Fatal(encodeCode, decode)
		}
	}
}
