package somatree

import "sync"

//
// Interface: SomaTreeFinder
func (tec *Cluster) Find(f FindRequest, b bool) SomaTreeAttacher {
	if findRequestCheck(f, tec) {
		return tec
	}
	var (
		wg             sync.WaitGroup
		rawResult, res chan SomaTreeAttacher
	)
	if len(tec.Children) == 0 {
		goto skip
	}
	if f.ElementId != "" {
		if _, ok := tec.Children[f.ElementId]; ok {
			return tec.Children[f.ElementId]
		} else {
			// f.ElementId is not a child of ours
			goto skip
		}
	}
	if f.ElementType != "" && f.ElementType != "node" {
		// searched element can't be a child of a cluster
		goto skip
	}
	rawResult = make(chan SomaTreeAttacher, len(tec.Children))
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(fr FindRequest, bl bool) {
			defer wg.Done()
			rawResult <- tec.Children[c].(SomaTreeFinder).Find(fr, bl)
		}(f, false)
	}
	wg.Wait()
	close(rawResult)

	res = make(chan SomaTreeAttacher, len(rawResult))
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
			return tec.Fault
		} else {
			return nil
		}
	case len(res) > 1:
		return tec.Fault
	}
	return <-res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
