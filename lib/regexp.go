package lib

import (
	"errors"
	"regexp"
)

/*正则查找字符串，返回最后一次匹配的内容*/
func ExpFindLast(rule string, srcString string) string {
	re, _ := regexp.CompilePOSIX(rule)
	find := re.FindAllStringSubmatch(srcString, -1)
	if findLen := len(find); findLen >= 1 {
		last := find[findLen-1]
		if lastLen := len(last); lastLen >= 1 {
			return last[lastLen-1]
		}
	}
	return ""
}

func GetDialLog(f string) (result map[string]string, err error) {
	result = make(map[string]string)
	result["ppp"] = ExpFindLast(`interface (.*)`, f)
	result["local"] = ExpFindLast(`local  IP address (.*)`, f)
	result["remote"] = ExpFindLast(`remote IP address (.*)`, f)
	result["dns1"] = ExpFindLast(`primary   DNS address (.*)`, f)
	result["dns2"] = ExpFindLast(`secondary DNS address (.*)*`, f)
	result["error"] = ExpFindLast(`Connect.*\n([^\f]*)Connection`, f)
	if result["error"] != "" {
		err = errors.New(result["error"])
	}
	return result, err
}
