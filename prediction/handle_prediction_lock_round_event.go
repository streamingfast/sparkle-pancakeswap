package prediction

import (
	"fmt"
	"math/big"

	"github.com/streamingfast/sparkle/entity"
	pbcodec "github.com/streamingfast/sparkle/pb/sf/ethereum/codec/v1"
)

func (s *Subgraph) HandlePredictionLockRoundEvent(trace *pbcodec.TransactionTrace, ev *PredictionLockRoundEvent) error {
	if s.StepBelow(2) {
		return nil
	}
	round := NewRound(fmt.Sprintf("%d", ev.Epoch))
	if err := s.Load(round); err != nil {
		return err
	}
	if !round.Exists() {
		return fmt.Errorf("Tried to lock round without an existing round (epoch: %s).", ev.Epoch)
	}

	round.LockAt = entity.NewInt(bi().SetInt64(s.Block().Timestamp().Unix())).Ptr()
	round.LockBlock = entity.NewIntFromLiteral(int64(s.Block().Number())).Ptr()
	round.LockHash = trace.Hash

	round.LockPrice = entity.NewFloat(bf().Quo(
		bf().SetInt(ev.Price),
		big.NewFloat(1e8), // 8 decimals
	)).Ptr()

	if err := s.Save(round); err != nil {
		return err
	}

	return nil
}
