package main

type openCollection struct {
	collection [][]card // 桌上牌
}

type closeCollection struct {
	collection []card // 未发牌队列
}

type doneCollection struct {
	collection []card //已消除队列
}
