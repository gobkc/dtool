package lib

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type RtTables struct {
	TableId   string
	TableName string
}

func GetRtTableAll() (map[string]string, map[string]string, error) {
	rtTablePath := "/etc/iproute2/rt_tables"
	var rtTableMap map[string]string
	rtTableMap = make(map[string]string)
	frtTableMap := make(map[string]string)
	var rtTableString string
	if rtTableBytes, err := ioutil.ReadFile(rtTablePath); err != nil {
		return nil, nil, err
	} else {
		rtTableString = string(rtTableBytes)
	}

	reg := regexp.MustCompile(`[0-9]+ [a-zA-Z0-9\@\.\_]+`)
	splitArr := reg.FindAllString(rtTableString, -1)
	if splitArr != nil {
		for _, v := range splitArr {
			splitRow := strings.Split(v, " ")
			if len(splitRow) < 2 {
				continue
			}
			rtTableMap[splitRow[0]] = splitRow[1]
			frtTableMap[splitRow[1]] = splitRow[0]
		}
	}
	return rtTableMap, frtTableMap, nil
}

func SetRtTable(rtName string) int {
	rt, frt, _ := GetRtTableAll()
	rt["0"] = "unspec"
	rt["253"] = "default"
	rt["254"] = "main"
	rt["255"] = "local"

	var rtId int

	/*如果已经存在设置过的就不再设置*/
	if id, ok := frt[rtName]; ok {
		rtId, _ := strconv.Atoi(id)
		return rtId
	}

	/*第一次遍历，写入MAP*/
	for i := 1; i <= 252; i++ {
		k := strconv.Itoa(i)
		if _, ok := rt[k]; !ok {
			rt[k] = rtName
			rtId = i
			break
		}
	}

	/*排序*/
	var waitWrite string
	for i := 0; i <= 255; i++ {
		k := strconv.Itoa(i)
		if v, ok := rt[k]; ok {
			waitWrite += fmt.Sprintf("%s %s\n", k, v)
		}
	}

	rtFile := "/etc/iproute2/rt_tables"
	if err := ioutil.WriteFile(rtFile, []byte(waitWrite), 0644); err != nil {
		fmt.Println(err)
		rtId = 0
	}
	return rtId
}
