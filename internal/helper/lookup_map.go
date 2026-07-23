package helper

func SliceToLookupMap[TType comparable](in []TType) map[TType]bool {
	result := make(map[TType]bool, len(in))
	if in == nil {
		return result
	}

	for _, item := range in {
		result[item] = true
	}

	return result
}
