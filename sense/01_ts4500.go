package sense

import "maps"

type TS4500 struct{ Standard }

func (x *TS4500) Init() {
	x.Standard.Init()
	maps.Copy(x.standardSenseCodes, map[string]string{
		"0000": "OVERWRITTEN",
	})
}
