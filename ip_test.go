package main

import "testing"

func TestGetLocalIPAddress(t *testing.T) {
	ip, _ := getLocalIPAddress()
	got := ip.String()
	want := "192.168.1.64"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
