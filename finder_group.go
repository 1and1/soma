package somatree

import "sync"

//
// Interface: Finder
func (teg *Group) Find(f FindRequest, b bool) Attacher {
	if findRequestCheck(f, teg) {
		return teg
	}
	var (
		wg             sync.WaitGroup
		rawResult, res chan Attacher
	)
	if len(teg.Children) == 0 {
		goto skip
	}
	rawResult = make(chan Attacher, len(teg.Children))
	for child, _ := range teg.Children {
		wg.Add(1)
		c := child
		go func(fr FindRequest, bl bool) {
			defer wg.Done()
			rawResult <- teg.Children[c].(Finder).Find(fr, bl)
		}(f, false)
	}
	wg.Wait()
	close(rawResult)

	res = make(chan Attacher, len(rawResult))
	for sta := range rawResult {
		if sta != nil {
			res <- sta
		}
	}
	close(res)
skip:
	switch {
	case len(res) == 0:
		if b {
			return teg.Fault
		} else {
			return nil
		}
	case len(res) > 1:
		return teg.Fault
	}
	return <-res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
