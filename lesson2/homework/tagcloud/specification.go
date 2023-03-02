package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	memory []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() TagCloud {
	return TagCloud{}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (cloud *TagCloud) AddTag(tag string) {
	exist := false
	for idx, value := range cloud.memory {
		if value.Tag == tag {
			cloud.memory[idx].OccurrenceCount += 1
			exist = true
			break
		}
	}

	if !exist {
		cloud.memory = append(cloud.memory, TagStat{tag, 1})
	}
	sort.Slice(cloud.memory, func(i, j int) bool {
		return cloud.memory[i].OccurrenceCount > cloud.memory[j].OccurrenceCount
	})

}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (cloud TagCloud) TopN(n int) []TagStat {
	var ans []TagStat
	for i := 0; i < n && i < len(cloud.memory); i++ {
		ans = append(ans, cloud.memory[i])
	}
	return ans
}
