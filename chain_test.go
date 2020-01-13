package chain

import (
	ru "github.com/jmcvetta/randutil"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	expected := MarkovChain{}
	expected.States = make(map[string][]ru.Choice)

	actual := NewChain()
	if reflect.DeepEqual(expected, actual) == false {
		t.Error("The NewChain function did not correctly instantiate a MarkovChain")
	}
}

type Args struct {
	Key, State  string
	Probability int
}

func TestSetState(t *testing.T) {
	mc := NewChain()
	tests := []struct {
		Args                 Args
		ExpectedChoiceLength int
		Expected             ru.Choice
		Name                 string
	}{
		{
			Args:                 Args{"a", "b", 1},
			Expected:             ru.Choice{Weight: 1, Item: "b"},
			ExpectedChoiceLength: 1,
			Name:                 "SetState adds a state when the chain doesn't contain the key or value",
		},
		{
			Args:                 Args{"a", "c", 4},
			ExpectedChoiceLength: 2,
			Expected:             ru.Choice{Weight: 4, Item: "c"},
			Name:                 "SetState adds a state when the chain has the key and not the value",
		},
		{
			Args:                 Args{"a", "b", 12},
			ExpectedChoiceLength: 2,
			Expected:             ru.Choice{Weight: 12, Item: "b"},
			Name:                 "SetState sets the Weight to a specific one even when the node was previously set.",
		},
	}
	for _, test := range tests {
		testArgs := test.Args
		mc.SetState(testArgs.Key, testArgs.State, testArgs.Probability)

		var actual ru.Choice
		for i := 0; i < len(mc.States[testArgs.Key]); i++ {
			if mc.States[testArgs.Key][i].Item == test.Expected.Item {
				actual = mc.States[testArgs.Key][i]
			}
		}

		if actual.Weight != test.Expected.Weight {
			t.Errorf("%s failed SetState with the following results:\nactual: %v, expected: %v\n", test.Name, actual.Weight, test.Expected.Weight)
		}

		actualChoiceLength := len(mc.States[testArgs.Key])
		if actualChoiceLength != test.ExpectedChoiceLength {
			t.Errorf("The length of Choices is incorrect for the following test: %s", test.Name)
		}

	}
}

func TestIncrementState(t *testing.T) {
	mc := NewChain()
	tests := []struct {
		Args     Args
		Expected ru.Choice
		Name     string
	}{
		{
			Name:     "inserts correctly when the key and state both don't exist",
			Expected: ru.Choice{Weight: 1.0, Item: "b"},
			Args:     Args{"a", "b", 0},
		},
		{
			Name:     "increments correctly when the key and state both exist already",
			Expected: ru.Choice{Weight: 2.0, Item: "b"},
			Args:     Args{"a", "b", 0},
		},
		{
			Name:     "inserts state correctly when the key exists and the state does not",
			Expected: ru.Choice{Weight: 1.0, Item: "c"},
			Args:     Args{"a", "c", 0},
		},
	}

	for _, test := range tests {
		testArgs := test.Args

		mc.IncrementState(testArgs.Key, testArgs.State)

		var actual ru.Choice
		for i := 0; i < len(mc.States[testArgs.Key]); i++ {
			if mc.States[testArgs.Key][i].Item == test.Expected.Item {
				actual = mc.States[testArgs.Key][i]
			}
		}

		if actual.Weight != test.Expected.Weight {
			t.Errorf("%s failed with the following results:\nactual: %v, expected: %v\n", test.Name, actual, test.Expected)
		}
	}
}

func TestPredictState(t *testing.T) {
	mc := NewChain()
	mc.IncrementState("a", "b")
	mc.IncrementState("a", "b")
	mc.SetState("a", "c", 0)
	mc.IncrementState("b", "c")

	expected := "b"
	actual, err := mc.PredictState("a")
	if err != nil {
		t.Error(err.Error())
	}

	if actual != expected {
		t.Errorf("TestPredictState failed with the following results:\nactual: %s, expected: %s\n", actual, expected)
	}
}

func TestPredictStateFailures(t *testing.T) {
	mc := NewChain()
	actual, err := mc.PredictState("a")
	if err == nil {
		t.Errorf("TestPredictState should fail when trying to predict an empty key.  Returned the following instead:\nActual: %s, Error: %s", actual, err)
	}

	mc.SetState("a", "b", 0)
	mc.SetState("a", "c", 0)
	actual, err = mc.PredictState("a")
	if err == nil {
		t.Errorf("TestPredictState should fail when trying to predict a key with no non-zero states.  Returned the following instead:\nActual: %s, Error: %s", actual, err)
	}

	mc.SetState("a", "d", -1)
	mc.SetState("a", "e", 1)

	actual, err = mc.PredictState("a")
	if err == nil {
		t.Errorf("TestPredictState should fail when trying to predict a key that has a negative integer as a state.  Returned the following instead:\nActual: %s, Error: %s", actual, err)
	}
}
