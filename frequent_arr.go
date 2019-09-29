package localcache

type FreqArr []*LFUItem

func (f *FreqArr) Len() int {
	return len(*f)
}

func (f *FreqArr) Less(i, j int) bool {
	return (*f)[i].Weight < (*f)[j].Weight
}

func (f *FreqArr) Swap(i, j int) {
	// 位置变换的时候置换 index
	(*f)[i], (*f)[j] = (*f)[j], (*f)[i]
	(*f)[i].Index = i
	(*f)[j].Index = j
}

func (f *FreqArr) Push(x interface{}) {
	value := x.(*LFUItem)
	*f = append(*f, value)
	value.Index = (*f).Len() // 放入的时候设置 index
}

func (f *FreqArr) Pop() interface{} {
	n := len(*f)
	ret := (*f)[n-1]
	*f = (*f)[:n-1]
	return ret
}
