// 详细原理请看这里 https://blog.frostmiku.com/archives/33/

package main

import (
	"strings"
)

const Base = uint(len(Chars))

func genCode(id uint) string {
	var b [CodeLength]uint
	var res string

	pid := id*Prime + Salt //扩大
	b[0] = pid
	for i := 0; i < CodeLength-1; i++ {
		b[i+1] = b[i] / Base
		b[i] = (b[i] + uint(i)*b[0]) % Base
	}

	// 校验位
	for i := 0; i < CodeLength-1; i++ {
		b[CodeLength-1] += b[i]
	}
	b[CodeLength-1] = b[CodeLength-1] * Prime % CodeLength

	for i := 0; i < CodeLength; i++ {
		res += string(Chars[b[(i*Prime2)%CodeLength]]) // 洗牌
	}

	//fmt.Println(res)

	return res
}

func decode(code string) int {
	if len(code) != CodeLength {
		return -1
	}
	var b [CodeLength]uint

	// 反洗牌
	for i := 0; i < CodeLength; i++ {
		b[(i*Prime2)%CodeLength] = uint(i)
	}

	// 转换回 Chars 下标
	for i := 0; i < CodeLength; i++ {
		j := strings.Index(Chars, string(code[b[i]]))
		if j == -1 {
			return -1 // 非法字符检查
		}
		b[i] = uint(j)
	}

	// 校验
	var expect uint
	for i := 0; i < CodeLength-1; i++ {
		expect += b[i]
	}
	expect = expect * Prime % CodeLength
	if b[5] != expect {
		return -1
	}

	// 反函数
	for i := CodeLength - 2; i >= 0; i-- {
		b[i] = (b[i] - uint(i)*(b[0]-Base)) % Base
	}
	var res uint = 0
	for i := CodeLength - 2; i > 0; i-- {
		res += b[i]
		res *= Base
	}

	// 反放大
	res = ((res + b[0]) - Salt) / Prime

	//fmt.Printf("%d\n",res)
	return int(res)
}
