package tree

/*
 * Fault Handler Special Sauce
 *
 * Elemnts return pointers to the Fault Handler instead of nil pointers
 * when asked for something they do not have.
 *
 * This makes these chains safe:
 *		<foo>.Parent.(Receiver).GetBucket().Unlink()
 *
 * Instead of nil, the parent returns the Fault handler which implements
 * Receiver and Unlinker. Due to the information in the
 * Receive-/UnlinkRequest, it can log what went wrong.
 *
 */

//
// Interface: Receiver
func (tef *Fault) Receive(r ReceiveRequest) {
	panic(`Fault.Receive`)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
