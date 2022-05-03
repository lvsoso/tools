package main

import "testing"

func TestPass(t *testing.T) {
	t.Log(rtest pass")
}

func TestFail(t *testing.T) {
	t.Fatal("test fail")
}

func TestH(t *testing.T) {

}
