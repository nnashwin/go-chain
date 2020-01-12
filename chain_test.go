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
	Key, State string
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
			Args:     Args{"a", "b"},
		},
		{
			Name:     "increments correctly when the key and state both exist already",
			Expected: ru.Choice{Weight: 2.0, Item: "b"},
			Args:     Args{"a", "b"},
		},
		{
			Name:     "inserts state correctly when the key exists and the state does not",
			Expected: ru.Choice{Weight: 1.0, Item: "c"},
			Args:     Args{"a", "c"},
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
