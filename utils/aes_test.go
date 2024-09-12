package utils

import (
	"fmt"
	"testing"
)

func TestDesECBEncrypt(t *testing.T) {
	content := "gI3A7FQmFxKhqwkpwhjqEmb88LRQOynnUduintLjtSut1JVmN4MrRgH7Mb3w2SDXAwi/JhZsfABY2SXOcXAK3VQWEP4aHvVts5SKOXvk6jW11EpLjRwscLadjU1WYOBScfkcyH4TM1h0Tr6BjH031QKMnyrGPTBXTAxPK2jKxhA="
	sign := "dde6bc48933f8378f78dbc17835b3084"
	data, err := AesECBBase64Decrypt(content, []byte("6jVni23sES2zLOHq"))

	time := "1691628139606"
	md5 := Md5ToString(fmt.Sprintf("%s%s", time, "6jVni23sES2zLOHq"))
	fmt.Println(md5, sign)
	fmt.Println(data, err)
}
