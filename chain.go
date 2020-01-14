package chain

import (
	"fmt"

	ru "github.com/jmcvetta/randutil"
)

type MarkovChain struct {
	States map[string][]ru.Choice
}

func (mc *MarkovChain) SetState(key, state string, prob int) {
	// If the key doesn't exist, append the key to the state transition graph
	if _, ok := mc.States[key]; ok == false {
		mc.States[key] = append(mc.States[key], ru.Choice{prob, state})
	} else {
		// Use struct to store if the int has been set or not
		// and the index where the value was found
		var sStatus struct {
			HasBeenSet bool
			Idx        int
		}

		ss := mc.States[key]
		for i := 0; i < len(ss); i++ {
			if ss[i].Item == state {
				sStatus.HasBeenSet = true
				sStatus.Idx = i
			}
		}

		if sStatus.HasBeenSet {
			mc.States[key][sStatus.Idx] = ru.Choice{prob, state}
		} else {
			mc.States[key] = append(mc.States[key], ru.Choice{prob, state})
		}
	}
}

func (mc *MarkovChain) IncrementState(key, state string) {
	// check for key and add key and key and state if they both don't exist
	if _, ok := mc.States[key]; ok == false {
		mc.States[key] = append(mc.States[key], ru.Choice{1, state})
	} else {
		var desiredState *ru.Choice
		ss := mc.States[key]
		for i := 0; i < len(ss); i++ {
			if ss[i].Item == state {
				desiredState = &ss[i]
			}
		}

		if desiredState == nil {
			mc.States[key] = append(mc.States[key], ru.Choice{1, state})
		} else {
			desiredState.Weight += 1
		}

	}
}

func (mc *MarkovChain) PredictState(key string) (string, error) {
	if _, ok := mc.States[key]; ok == false {
		return "", fmt.Errorf("You have entered a key that does not exist. Add that key to the chain and run again.")
	}

	c, err := ru.WeightedChoice(mc.States[key])
	if err != nil {
		return "", fmt.Errorf("The weighted random selection failed.  Check and ensure your MarkovChain has no negative values and at least one non-zero value.")
	}

	// casting interface as a string to force the use of string
	return c.Item.(string), nil
}

func (mc *MarkovChain) GenerateStates(key string, numStates int) ([]string, error) {
	ss := make([]string, 0, numStates)
	for i := 0; i < numStates; i++ {
		ns, err := mc.PredictState(key)
		if err != nil {
			return ss, err
		}

		ss = append(ss, ns)
		key = ns
	}

	return ss, nil
}

func NewChain() *MarkovChain {
	return &MarkovChain{States: make(map[string][]ru.Choice)}
}
