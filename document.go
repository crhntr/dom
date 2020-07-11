// +build js,wasm

package dom

import (
	"syscall/js"
)

type document byte
type window byte

const (
	Document document = 0
	Window   window   = 0
)

var (
	doc = js.Global().Get("document")
	win = js.Global()
)

func (_ document) JSValue() js.Value                           { return doc }
func (_ document) Get(key string) js.Value                     { return doc.Get(key) }
func (_ document) Set(key string, value interface{})           { doc.Set(key, value) }
func (_ document) Call(m string, args ...interface{}) js.Value { return doc.Call(m, args...) }
func (_ document) Type() js.Type                               { return doc.Type() }
func (_ document) Delete(p string)                             { doc.Delete(p) }

func (el document) AddEventListener(eventName string, listener EventListener) {
	el.Call("addEventListener", eventName, js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		listener.HandleEvent(Event(args[0]))
		return nil
	}))
}

func (el document) AddEventListenerFunc(eventName string, listener EventListenerFunc) {
	el.AddEventListener(eventName, listener)
}

func (el document) AddEventListenerChannel(eventName string, c chan Event) {
	el.AddEventListenerFunc(eventName, func(event Event) {
		c <- event
	})
}

func (_ window) JSValue() js.Value                           { return win }
func (_ window) Get(key string) js.Value                     { return win.Get(key) }
func (_ window) Set(key string, value interface{})           { win.Set(key, value) }
func (_ window) Call(m string, args ...interface{}) js.Value { return win.Call(m, args...) }
func (_ window) Type() js.Type                               { return win.Type() }
func (_ window) Delete(p string)                             { win.Delete(p) }

func (el window) AddEventListener(eventName string, listener EventListener) {
	el.Call("addEventListener", eventName, js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		listener.HandleEvent(Event(args[0]))
		return nil
	}))
}

func (el window) AddEventListenerFunc(eventName string, listener EventListenerFunc) {
	el.AddEventListener(eventName, listener)
}

func (el window) AddEventListenerChannel(eventName string, c chan Event) {
	el.AddEventListenerFunc(eventName, func(event Event) {
		c <- event
	})
}
