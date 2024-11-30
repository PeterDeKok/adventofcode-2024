package bubbles

import (
	"peterdekok.nl/adventofcode/twentytwentyfour/src/manage/op/result/question"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles/answer"
	"peterdekok.nl/adventofcode/twentytwentyfour/src/tui/bubbles/list"
)

type List = list.List
type ListItem = list.Item

func NewList[V ListItem](component string, items []V, delegate list.Delegater) *List {
	return list.New[V](component, items, delegate)
}

type Answer = answer.Model

func NewAnswer(component string, a *question.Answer) *Answer {
	return answer.New(component, a)
}
