package somatree

/*
 * Fault Handler Special Sauce
 *
 * Elemnts return pointers to the Fault Handler instead of nil pointers
 * when asked for something they do not have.
 *
 * This makes these chains safe:
 *		<foo>.Parent.(SomaTreeReceiver).GetBucket().Unlink()
 *
 * Instead of nil, the parent returns the Fault handler which implements
 * SomaTreeReceiver and SomaTreeUnlinker. Due to the information in the
 * Receive-/UnlinkRequest, it can log what went wrong.
 *
 */

//
// Interface: SomaTreeUnlinker
func (tef *Fault) Unlink(u UnlinkRequest) {
	panic(`Fault.Unlink`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
