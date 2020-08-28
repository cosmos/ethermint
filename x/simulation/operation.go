package simulation

import (
	"math/rand"
	"sort"

	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// queueOperations adds all future operations into the operation queue.
func queueOperations(queuedOps simulation.OperationQueue,
	queuedTimeOps []simulation.FutureOperation, futureOps []simulation.FutureOperation) {

	if futureOps == nil {
		return
	}

	for _, futureOp := range futureOps {
		futureOp := futureOp
		if futureOp.BlockHeight != 0 {
			if val, ok := queuedOps[futureOp.BlockHeight]; ok {
				queuedOps[futureOp.BlockHeight] = append(val, futureOp.Op)
			} else {
				queuedOps[futureOp.BlockHeight] = []simulation.Operation{futureOp.Op}
			}
			continue
		}

		// TODO: Replace with proper sorted data structure, so don't have the
		// copy entire slice
		index := sort.Search(
			len(queuedTimeOps),
			func(i int) bool {
				return queuedTimeOps[i].BlockTime.After(futureOp.BlockTime)
			},
		)
		queuedTimeOps = append(queuedTimeOps, simulation.FutureOperation{})
		copy(queuedTimeOps[index+1:], queuedTimeOps[index:])
		queuedTimeOps[index] = futureOp
	}
}

//________________________________________________________________________

// WeightedOperation is an operation with associated weight.
// This is used to bias the selection operation within the simulator.
type WeightedOperation struct {
	Weight int
	Op     simulation.Operation
}

// NewWeightedOperation creates a new WeightedOperation instance
func NewWeightedOperation(weight int, op simulation.Operation) WeightedOperation {
	return WeightedOperation{
		Weight: weight,
		Op:     op,
	}
}

// WeightedOperations is the group of all weighted operations to simulate.
type WeightedOperations []WeightedOperation

func (ops WeightedOperations) totalWeight() int {
	totalOpWeight := 0
	for _, op := range ops {
		totalOpWeight += op.Weight
	}
	return totalOpWeight
}

type selectOpFn func(r *rand.Rand) simulation.Operation

func (ops WeightedOperations) getSelectOpFn() selectOpFn {
	totalOpWeight := ops.totalWeight()
	return func(r *rand.Rand) simulation.Operation {
		x := r.Intn(totalOpWeight)
		for i := 0; i < len(ops); i++ {
			if x <= ops[i].Weight {
				return ops[i].Op
			}
			x -= ops[i].Weight
		}
		// shouldn't happen
		return ops[0].Op
	}
}
