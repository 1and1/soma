package main

func (tk *treeKeeper) drain(s string) (j int) {
	switch s {
	case "action":
		j = len(tk.actionChan)
		for i := j; i > 0; i-- {
			<-tk.actionChan
		}
	case "error":
		j = len(tk.errChan)
		for i := j; i > 0; i-- {
			<-tk.errChan
		}
	default:
		panic("Requested drain for unknown")
	}
	return j
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
