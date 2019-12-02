package router

import (
	"regexp"
)

type Pattern struct {
	PatternString  string
	ControllerName string
	ActionName     string
	PatternLen     int
	PatternMap     []otherParam
}

type otherParam struct {
	Re   regexp.Regexp
	Name string
}
type freeFishUrl struct {
	ControllerName   string
	ControllerAction string
	OtherKeyMap      map[string]interface{}
}

func (p *Pattern) isMatch(pathSplitList []string) (bool, *freeFishUrl) {
	f := new(freeFishUrl)
	f.OtherKeyMap = map[string]interface{}{}
	if len(pathSplitList) == p.PatternLen {
		for i := 0; i < p.PatternLen; i++ {
			if !p.PatternMap[i].Re.MatchString(pathSplitList[i]) {
				return false, f
			}
		}
		f.ControllerName = p.ControllerName
		f.ControllerAction = p.ActionName
		return true, f
	}
	return false, f

}
