package util

//StrPtrMapEqual checks two strptr maps for value equality
func StrPtrMapEqual(this map[string]*string, that map[string]*string) bool {
	if len(this) != len(that) {
		return false
	}
	for k, thisV := range this {
		thatV, ok := that[k]
		if !ok {
			return false // key not present
		}
		if thisV == nil && thatV == nil {
			continue // both nil ok
		}
		if thisV == nil || thatV == nil {
			return false // either but not both nil not ok
		}
		if *thisV != *thatV {
			// guaranteed no nil, so check for value equality
			return false
		}
	}

	return true
}
