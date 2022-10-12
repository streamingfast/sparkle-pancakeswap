package prediction

import (
	"fmt"

	pbcodec "github.com/streamingfast/sparkle/pb/sf/ethereum/type/v2"
)

func (s *Subgraph) HandlePredictionClaimEvent(trace *pbcodec.TransactionTrace, ev *PredictionClaimEvent) error {
	if s.StepBelow(2) {
		return nil
	}
	betID := fmt.Sprintf("%s%x", ev.Sender.Pretty(), ev.CurrentEpoch.Int64())

	bet := NewBet(betID)
	if err := s.Load(bet); err != nil {
		return err
	}
	if !bet.Exists() {
		return nil
	}

	bet.Claimed = true
	bet.ClaimedHash = trace.Hash
	return s.Save(bet)
}
