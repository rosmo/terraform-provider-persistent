package provider

import (
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	input := []string{}
	initial := 5
	last, res := assignKeys(input, nil, false, int64(initial), int64(initial))
	t.Logf("input keys: %v, output: %v", input, res)
	expectedLast := initial + len(input)
	if last != int64(expectedLast) {
		t.Errorf("Expecting last to be %d, got %d", expectedLast, last)
	}
	if len(res) != len(input) {
		t.Errorf("Result map length %d != key length %d", len(res), len(input))
	}
}

func TestInitial(t *testing.T) {
	input := []string{"a", "c", "b"}
	initial := 5
	last, res := assignKeys(input, nil, false, int64(initial), int64(initial-1))
	t.Logf("input keys: %v, output: %v", input, res)
	expectedLast := initial + len(input) - 1
	if last != int64(expectedLast) {
		t.Errorf("Expecting last to be %d, got %d", expectedLast, last)
	}
	if len(res) != len(input) {
		t.Errorf("Result map length %d != key length %d", len(res), len(input))
	}
	expected := map[string]int64{"a": 5, "b": 6, "c": 7}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Expected %v got %v", expected, res)
	}
}

func TestNop(t *testing.T) {
	input := []string{"a", "c", "b"}
	state := map[string]int64{"a": 5, "b": 9, "c": 11}
	initial := 5
	last := 11
	last2, res := assignKeys(input, state, false, int64(initial), int64(last))
	if !reflect.DeepEqual(res, state) {
		t.Errorf("Expected %v got %v", state, res)
	}
	if int64(last) != last2 {
		t.Errorf("Expected last value to stay at %d but got %d", last, last2)
	}
	last2, res = assignKeys(input, state, true, int64(initial), int64(last))
	if !reflect.DeepEqual(res, state) {
		t.Errorf("Expected %v got %v", state, res)
	}
	if int64(last) != last2 {
		t.Errorf("Expected last value to stay at %d but got %d", last, last2)
	}
}

func TestChange(t *testing.T) {
	input := []string{"d", "a", "c"}
	state := map[string]int64{"a": 5, "b": 6, "c": 7}
	initial := 5
	last := state["c"]
	last2, res := assignKeys(input, state, false, int64(initial), int64(last))
	expected := map[string]int64{"a": 5, "c": 7, "d": 8}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Expected %v got %v", state, res)
	}
	if int64(last2) != last+1 {
		t.Errorf("Expected last value to stay at %d but got %d", last, last2)
	}
	expected = map[string]int64{"a": 5, "c": 7, "d": 6}
	last2, res = assignKeys(input, state, true, int64(initial), int64(last))
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Expected %v got %v", state, res)
	}
	if int64(last2) != expected["d"] {
		t.Errorf("Expected last value to stay at %d but got %d", state["b"], last2)
	}
}
