package chain

import (
	ru "github.com/jmcvetta/randutil"
)

type MarkovChain struct {
	States map[string][]ru.Choice
}

func (mc *MarkovChain) Insert(key, state string) {
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
			mc.States[key] = append(ss, ru.Choice{1, state})
		} else {
			desiredState.Weight += 1
		}

	}
}

func NewChain() MarkovChain {
	mc := MarkovChain{}
	mc.States = make(map[string][]ru.Choice)

	return mc
}
