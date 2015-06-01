package timediff

import (
	"testing"
	"time"
)

func TestGetDifference(t *testing.T) {

	future := time.Now().Add(time.Second * 3)
	d := GetDifference(future, time.Now())
	if d != "3 Seconds ago" {
		t.Fatalf("Expected 3 seconds ago, got %s.", d)
	}

	future = time.Now().Add(time.Minute * 3)
	d = GetDifference(future, time.Now())
	if d != "3 Minutes ago" {
		t.Fatalf("Expected 3 minutes ago, got %s.", d)
	}

	future = time.Now().Add(time.Hour * 3)
	d = GetDifference(future, time.Now())
	if d != "3 Hours ago" {
		t.Fatalf("Expected 3 hours ago, got %s.", d)
	}

	future = time.Now().Add(time.Hour * 24 * 2)
	d = GetDifference(future, time.Now())
	if d != "2 Days ago" {
		t.Fatalf("Expected 2 days ago, got %s.", d)
	}

	future = time.Now().Add(time.Hour * 24 * 8)
	d = GetDifference(future, time.Now())
	if d != "1 Week ago" {
		t.Fatalf("Expected 1 week ago, got %s.", d)
	}

	future = time.Now().Add(time.Hour * 24 * 35)
	d = GetDifference(future, time.Now())
	if d != "1 Month ago" {
		t.Fatalf("Expected 1 month ago, got %s.", d)
	}

	future = time.Now().Add(time.Hour * 24 * 365 * 2)
	d = GetDifference(future, time.Now())
	if d != "2 Years ago" {
		t.Fatalf("Expected 2 years ago, got %s.", d)
	}
}
