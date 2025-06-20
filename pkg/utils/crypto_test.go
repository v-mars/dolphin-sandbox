package utils

import (
	"fmt"
	"testing"
)

func TestDecryptDES(t *testing.T) {
	//TestAll()
	//AesDemo2()
	//aa := DeTxtByAes("15bbddabdf3113087b72336ac2b3063c", "xxx")
	//fmt.Println("aa:", aa)

	bb, err := EnTxtByAesWithErr("dsafkljasdlk;fj;lkadsjfl;kasdjfl;kadskfjl;kadsfkjl;asdfj;lasdfjl;askdfa", "WedolcMd1xxx*WSP")
	if err != nil {
		fmt.Println("err1:", err)
		return
	}
	fmt.Println("bb:", bb)
	cc, err := DeTxtByAesWithErr(bb, "WedolcMd1xxx*WSP")
	if err != nil {
		fmt.Println("err2:", err)
		return
	}
	fmt.Println("cc:", cc)
}
