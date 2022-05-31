package routed

type Routed struct {
	NextRouted any
}

func (r *Routed) Next() (next any, flag bool) {
	next = r.NextRouted
	flag = r.NextRouted != nil
	return
}
