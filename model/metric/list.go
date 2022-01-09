package metric

import "sort"

var _ sort.Interface = (ByString)(nil)

type List []*Metric

func (l List) ToMap() map[Type]List {
	m := make(map[Type]List)
	for _, v := range l {
		if _, ok := m[v.Type()]; !ok {
			m[v.Type()] = make(List, 0)
		}
		m[v.Type()] = append(m[v.Type()], v)
	}
	return m
}

type ByString List

func (m ByString) Len() int           { return len(m) }
func (m ByString) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m ByString) Less(i, j int) bool { return m[i].String() < m[j].String() }
