package main

import (
	"fmt"
	"sync"
)

type FSMState string            //状态
type FSMEvent string            //事件
type FSMHandler func() FSMState //处理方法，返回新的状态

//有限状态机
type FSM struct {
	mu       sync.Mutex                           //排他锁
	state    FSMState                             //当前状态
	handlers map[FSMState]map[FSMEvent]FSMHandler //处理地图集，每个状态都可以触发有限个状态
}

//获取当前状态
func (f *FSM) getState() FSMState {
	return f.state
}

//设置状态
func (f *FSM) setState(newState FSMState) {
	f.state = newState
}

//某状态添加事件处理方法
func (f *FSM) AddHandler(state FSMState, event FSMEvent, handler FSMHandler) *FSM {
	if _, ok := f.handlers[state]; !ok {
		f.handlers[state] = make(map[FSMEvent]FSMHandler)
	}
	if _, ok := f.handlers[state][event]; ok {
		f.handlers[state][event] = handler
	}
	f.handlers[state][event] = handler
	return f

}

//事件处理
func (f *FSM) Call(event FSMEvent) (FSMState,error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	events, ok := f.handlers[f.getState()]
	if !ok {
		return f.getState(),fmt.Errorf("无[%s]状态",f.getState())
	}
	fn, ok := events[event]
	if !ok {
		return f.getState(),fmt.Errorf("[%s]状态下无法触发[%s]事件",f.getState(),event)
	}
	old := f.getState()
	f.setState(fn())
	newstate := f.getState()
	fmt.Printf("状态从[%v]->[%v] \n", old, newstate)
	return f.getState(),nil
}

//实例化FSM
func NewFSM(initState FSMState) *FSM {
	return &FSM{
		state:    initState,
		handlers: make(map[FSMState]map[FSMEvent]FSMHandler),
	}
}

var (
	Poweroff       = FSMState("关闭电源")
	FirstGear      = FSMState("一档")
	SecondGear     = FSMState("二档")
	ThreadGear     = FSMState("三档")
	PowerOffDown   = FSMEvent("按下关闭")
	FirstGearDown  = FSMEvent("按下一档")
	SecondGearDown = FSMEvent("按下二档")
	ThreadGearDown = FSMEvent("按下三档")

	PoweroffHandle = FSMHandler(func() FSMState {
		fmt.Println("电源关闭")
		return Poweroff
	})
	FisrtGearHandle = FSMHandler(func() FSMState {
		fmt.Println("一档开启")
		return FirstGear
	})
	SecondGearHandle = FSMHandler(func() FSMState {
		fmt.Println("二档开启")
		return SecondGear
	})
	ThreadGearhandle = FSMHandler(func() FSMState {
		fmt.Println("三档开启")
		return ThreadGear
	})
)

type ElectrciFan struct {
	FSM *FSM
}

func NewElectrciFan(initSate FSMState) *ElectrciFan {
	return &ElectrciFan{
		FSM: NewFSM(initSate),
	}
}

func main() {
	e := NewElectrciFan(Poweroff) //初始是关闭的

	//关闭状态下,添加
	e.FSM.AddHandler(Poweroff, PowerOffDown, PoweroffHandle)   //关闭状态按下关闭会关闭
	e.FSM.AddHandler(Poweroff, FirstGearDown, FisrtGearHandle) //关-按下1档-1档开启
	e.FSM.AddHandler(Poweroff, SecondGearDown, SecondGearHandle)
	e.FSM.AddHandler(Poweroff, ThreadGearDown, ThreadGearhandle)
	//1档下添加

	e.FSM.AddHandler(FirstGear, PowerOffDown, PoweroffHandle)

	//2档下添加
	e.FSM.AddHandler(SecondGear, PowerOffDown, PoweroffHandle)

	//3档下添加
	e.FSM.AddHandler(ThreadGear, PowerOffDown, PoweroffHandle)

	_,err:=e.FSM.Call(FirstGearDown)
	if err!=nil {
		fmt.Println(err)
	}
	_,err=e.FSM.Call(SecondGearDown)
	if err!=nil {
		fmt.Println(err)
	}
	_,err=e.FSM.Call(PowerOffDown)
	if err!=nil {
		fmt.Println(err)
	}

}
