package prediction

import (
	"fmt"

	"github.com/streamingfast/sparkle/entity"
	pbcodec "github.com/streamingfast/sparkle/pb/sf/ethereum/codec/v1"
)

func (s *Subgraph) HandlePredictionPauseEvent(trace *pbcodec.TransactionTrace, ev *PredictionPauseEvent) error {
	if s.StepBelow(2) {
		return nil
	}
	market := NewMarket("1")
	if err := s.Load(market); err != nil {
		return fmt.Errorf("loading market: %w", err)
	}
	// if it didn't exist, well.. it'll be created.
	market.Epoch = S(ev.Epoch.String())
	market.Paused = true
	if err := s.Save(market); err != nil {
		return err
	}

	round := NewRound(ev.Epoch.String())
	if err := s.Load(round); err != nil {
		return err
	}
	if !round.Exists() {
		return nil
	}
	round.Failed = entity.NewBool(true).Ptr()

	if err := s.Save(round); err != nil {
		return err
	}

	if round.Previous == nil {
		return nil
	}

	previousRound := NewRound(*round.Previous)
	if err := s.Load(previousRound); err != nil {
		return err
	}
	if !previousRound.Exists() {
		return nil
	}
	previousRound.Failed = entity.NewBool(true).Ptr()
	if err := s.Save(previousRound); err != nil {
		return err
	}
	return nil
}
