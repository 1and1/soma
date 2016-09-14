/*-
 * Copyright (c) 2016, 1&1 Internet SE
 * Copyright (c) 2016, Jörg Pernfuß
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package tree

import "sync"

//
// Interface: Finder
func (tec *Cluster) Find(f FindRequest, b bool) Attacher {
	if findRequestCheck(f, tec) {
		return tec
	}
	var (
		wg             sync.WaitGroup
		rawResult, res chan Attacher
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
	rawResult = make(chan Attacher, len(tec.Children))
	for child, _ := range tec.Children {
		wg.Add(1)
		c := child
		go func(fr FindRequest, bl bool) {
			defer wg.Done()
			rawResult <- tec.Children[c].(Finder).Find(fr, bl)
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
