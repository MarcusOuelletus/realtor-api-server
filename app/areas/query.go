package areas

import (
	"strings"

	"github.com/MarcusOuelletus/rets-server/global"
)

type cacheQuery struct{}

func QueryCache(search string, offset int) []string {
	search = strings.ToLower(search)

	var q = new(cacheQuery)

	switch len(search) {
	case 0:
		return q.getFirstN(offset)
	case 1:
		return q.getFirstNFromOne(search, offset)
	default:
		return q.getFirstNFromTwo(search, offset)
	}
}

func (q *cacheQuery) getFirstNFromOne(search string, offset int) []string {
	var firstOne = string(search[0])
	firstOne = strings.ToLower(firstOne)

	var matches = make([]string, 0, global.AREA_LIMIT)

	for i, location := range cacheObject.OneLetterMap[firstOne] {
		if i < offset {
			continue
		}

		if len(search) > len(*location) {
			continue
		}

		var lowerCaseLocation = strings.ToLower(*location)
		var lowerCaseSearch = strings.ToLower(search)

		if strings.Contains(lowerCaseLocation, lowerCaseSearch) {
			matches = append(matches, *location)
		}

		if len(matches) == global.AREA_LIMIT {
			break
		}
	}

	return matches
}

func (q *cacheQuery) getFirstNFromTwo(search string, offset int) []string {
	var firstTwo = string(search[0]) + string(search[1])
	firstTwo = strings.ToLower(firstTwo)

	var matches = make([]string, 0, global.AREA_LIMIT)

	for i, location := range cacheObject.TwoLetterMap[firstTwo] {
		if i < offset {
			continue
		}

		if len(search) > len(*location) {
			continue
		}

		var lowerCaseLocation = strings.ToLower(*location)
		var lowerCaseSearch = strings.ToLower(search)

		if strings.Contains(lowerCaseLocation, lowerCaseSearch) {
			matches = append(matches, *location)
		}

		if len(matches) == global.AREA_LIMIT {
			break
		}
	}

	return matches
}

func (q *cacheQuery) getFirstN(offset int) []string {
	if offset > len(cacheObject.Slice) {
		return []string{}
	}

	if offset+global.AREA_LIMIT > len(cacheObject.Slice) {
		return cacheObject.Slice[offset:]
	}

	return cacheObject.Slice[offset:global.AREA_LIMIT]
}
