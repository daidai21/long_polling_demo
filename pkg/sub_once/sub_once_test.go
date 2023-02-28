package sub_once

import (
	"testing"
	"time"
)

func Test_Example(t *testing.T) {

	so := New()
	go func() {
		<-so.Sub()
		t.Log("go 1")
	}()
	go func() {
		<-so.Sub()
		t.Log("go 2")
		//time.Sleep(time.Second)
		<-so.Sub()
		t.Log("go 2 re")

	}()

	time.Sleep(time.Second)
	so.Pub()
	time.Sleep(2 * time.Second)
	so.Pub()

	time.Sleep(1 * time.Second)

	so.Close()

}
