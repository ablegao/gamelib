package gamelib

import (
	"strconv"
)

type MapConfig struct {
	MapPointString string
}

//转化地图坐标
func MapPointToArray(d string) [][]int {
	b := []byte(d)
	toData := [][]int{}
	i := 0
	tmp := []int{}
	for _, s := range b {
		ss := string(s)

		if ss != "," {
			i, _ := strconv.Atoi(ss)
			tmp = append(tmp, i)
			//to.Write([]byte(i))
			//binary.Write(to, binary.BigEndian, uint8(i))
		} else if ss == "," {
			i++
			if len(tmp) > 0 {
				toData = append(toData, tmp)
				tmp = []int{}
			}
		}
	}
	return toData
}
