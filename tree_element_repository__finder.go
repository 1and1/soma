package somatree

import "sync"

//
// Interface: SomaTreeFinder
func (ter *SomaTreeElemRepository) Find(f FindRequest, b bool) SomaTreeAttacher {
	if findRequestCheck(f, ter) {
		return ter
	}
	var (
		wg             sync.WaitGroup
		rawResult, res chan SomaTreeAttacher
	)
	if len(ter.Children) == 0 {
		goto skip
	}
	rawResult = make(chan SomaTreeAttacher, len(ter.Children))
	for child, _ := range ter.Children {
		wg.Add(1)
		c := child
		go func(fr FindRequest, bl bool) {
			defer wg.Done()
			rawResult <- ter.Children[c].(SomaTreeFinder).Find(fr, bl)
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
			return ter.Fault
		} else {
			return nil
		}
	case len(res) > 1:
		return ter.Fault
	}
	return <-res
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
