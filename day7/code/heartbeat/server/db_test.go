package main

import "sync/atomic"

type idGenModelAtomicTest struct {
	counter         uint32
	cGen, cCur, cAl uint32
}

type modelStub interface {
	GetCounters() (uint32, uint32, uint32)
}

func (g *idGenModelAtomicTest) Generate() uint32 {
	atomic.AddUint32(&g.cGen, 1)
	return atomic.AddUint32(&g.counter, 1)
}

func (g *idGenModelAtomicTest) Current() uint32 {
	atomic.AddUint32(&g.cCur, 1)
	return atomic.LoadUint32(&g.counter)
}

func (g *idGenModelAtomicTest) Alive() bool {
	atomic.AddUint32(&g.cAl, 1)
	return true
}

func (g *idGenModelAtomicTest) GetCounters() (uint32, uint32, uint32) {
	return g.cGen, g.cCur, g.cAl
}

func NewIdGenModelAtomicTest() *idGenModelAtomicTest {
	return new(idGenModelAtomicTest)
}
