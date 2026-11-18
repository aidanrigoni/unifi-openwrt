// SPDX-License-Identifier: GPL-2.0-or-later
package main

import "testing"

func TestConsoleFromURL(t *testing.T) {
	for in, want := range map[string]string{
		"http://10.0.0.1:8080/inform":      "10.0.0.1",
		"http://192.168.1.1:8080/inform":   "192.168.1.1",
		"http://[fd00::1]:8080/inform":     "fd00::1",
		"http://unifi.example:8080/inform": "unifi.example",
		"":                                 "192.168.1.1",
		"not a url":                        "192.168.1.1",
	} {
		if got := consoleFromURL(in); got != want {
			t.Errorf("consoleFromURL(%q) = %q, want %q", in, got, want)
		}
	}
}
