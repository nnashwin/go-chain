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

func TestInsert(t *testing.T) {
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

		mc.Insert(testArgs.Key, testArgs.State)

		var actual ru.Choice
		for i := 0; i < len(mc.States[testArgs.Key]); i++ {
			if mc.States[testArgs.Key][i].Item == test.Expected.Item {
				actual = mc.States[testArgs.Key][i]
			}
		}

		if actual.Weight != test.Expected.Weight {
			t.Errorf("%s failed with the following results:\nactual: %v, expected: %v", test.Name, actual, test.Expected)
		}
	}
}

func TestAddStateChoice(t *testing.T) {
	mc := NewChain()
	tests := []struct {
		Args     Args
		Expected ru.Choice
		Name     string
	}{
		{
			Args:     Args{"a", "b", 1},
			Expected: ru.Choice{Weight: 1, Item: "b"},
			Name:     "AddStateChoice adds a state when the chain doesn't contain the key or value",
		},
		{
			Args:     Args{"a", "c", 4},
			Expected: ru.Choice{Weight: 4, Item: "c"},
			Name:     "AddStateChoice adds a state when the chain has the key and not the value",
		},
		{
			Args:     Args{"a", "b", 12},
			Expected: ru.Choice{Weight: 12, Item: "b"},
			Name:     "AddStateChoice sets the Weight to a specific one even when the node was previously set.",
		},
	}
	for _, test := range tests {
		testArgs := test.Args
		mc.AddStateChoice(testArgs.Key, testArgs.State, testArgs.Probability)

		var actual ru.Choice
		for i := 0; i < len(mc.States[testArgs.Key]); i++ {
			if mc.States[testArgs.Key][i].Item == test.Expected.Item {
				actual = mc.States[testArgs.Key][i]
			}
		}

		if actual.Weight != test.Expected.Weight {
			t.Errorf("%s failed AddStateChoice with the following results:\nactual: %v, expected: %v", test.Name, actual.Weight, test.Expected.Weight)
		}
	}
}
