package tree

func receiveRequestCheck(r ReceiveRequest, b Builder) bool {
	if r.ParentType == b.GetType() && (r.ParentId == b.GetID() || r.ParentName == b.GetName()) {
		return true
	}
	return false
}

func unlinkRequestCheck(u UnlinkRequest, b Builder) bool {
	if u.ParentType == b.GetType() && (u.ParentId == b.GetID() || u.ParentName == b.GetName()) {
		return true
	}
	return false
}

func findRequestCheck(f FindRequest, b Builder) bool {
	if f.ElementId == b.GetID() || (f.ElementType == b.GetType() && f.ElementName == b.GetName()) {
		return true
	}
	return false
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
