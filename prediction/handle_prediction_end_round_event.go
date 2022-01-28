package prediction

import (
	"fmt"
	"math/big"

	"github.com/streamingfast/sparkle/entity"
	pbcodec "github.com/streamingfast/sparkle/pb/sf/ethereum/codec/v1"
)

func (s *Subgraph) HandlePredictionEndRoundEvent(trace *pbcodec.TransactionTrace, ev *PredictionEndRoundEvent) error {
	if s.StepBelow(3) {
		// round.LockPrice is computed on Step2 thus we need to compute
		// this on step 3
		return nil
	}
	round := NewRound(fmt.Sprintf("%d", ev.Epoch))
	if err := s.Load(round); err != nil {
		return err
	}
	if !round.Exists() {
		return fmt.Errorf("Tried to end round without an existing round (epoch: %s).", ev.Epoch)
	}

	round.EndAt = entity.NewIntFromLiteral(s.Block().Timestamp().Unix()).Ptr()
	round.EndBlock = entity.NewIntFromLiteral(int64(s.Block().Number())).Ptr()
	round.EndHash = trace.Hash
	round.ClosePrice = F(bf().Quo(
		bf().SetInt(ev.Price),
		big.NewFloat(1e8), // 8 decimals
	)).Ptr()

	cmp := round.ClosePrice.Float().Cmp(round.LockPrice.Float())
	if cmp == 0 {
		round.Position = &PositionHouse
	} else if cmp > 0 {
		round.Position = &PositionBull
	} else {
		round.Position = &PositionBear
	}

	round.Failed = entity.NewBool(false).Ptr()

	if err := s.Save(round); err != nil {
		return err
	}

	return nil
}
