package main

import "sync"

type Cash struct {
	m    sync.Mutex
	data map[string]Spending
}

func (a *Cash) Add(user string, s Spending) {
	a.m.Lock()
	defer a.m.Unlock()
	a.data[user] = s
}

func (a *Cash) Delete(user string) {
	a.m.Lock()
	defer a.m.Unlock()
	delete(a.data, user)
}

func (a *Cash) Get(user string) (Spending, bool) {
	a.m.Lock()
	defer a.m.Unlock()
	r, ok := a.data[user]
	return r, ok
}

var spendingCash Cash

//todo нужно будет сделать проверку на спец симфолы типа _:_
const deleteCategory = "deleteCategory:"
const newSpending = "newSpending:"
