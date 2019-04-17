package core

type vnd struct {
	objective     objective
	local2Opt     local2Opt
	localShifting localshifting
}

func (v vnd) process(x *Solution) {
	x2 := x

	for {
		v.localShifting.process(x2)
		v.local2Opt.process(x2)

		if v.objective.get(x2) < v.objective.get(x) {
			x = x2
		} else {
			break
		}
	}

	return
}
