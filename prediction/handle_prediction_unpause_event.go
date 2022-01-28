package prediction

import (
	pbcodec "github.com/streamingfast/sparkle/pb/sf/ethereum/codec/v1"
)

func (s *Subgraph) HandlePredictionUnpauseEvent(trace *pbcodec.TransactionTrace, ev *PredictionUnpauseEvent) error {
	if s.StepBelow(2) {
		return nil
	}
	market := NewMarket("1")
	if err := s.Load(market); err != nil {
		return err
	}
	// if it didn't exist, well.. it'll be created.
	market.Epoch = S(ev.Epoch.String())
	market.Paused = false
	if err := s.Save(market); err != nil {
		return err
	}
	return nil
}
