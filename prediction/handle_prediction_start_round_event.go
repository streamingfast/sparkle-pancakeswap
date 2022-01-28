package prediction

import (
	"fmt"
	"math/big"

	"github.com/streamingfast/sparkle/entity"
	pbcodec "github.com/streamingfast/sparkle/pb/sf/ethereum/codec/v1"
)

func (s *Subgraph) HandlePredictionStartRoundEvent(trace *pbcodec.TransactionTrace, ev *PredictionStartRoundEvent) error {
	market := NewMarket("1")
	if err := s.Load(market); err != nil {
		return err
	}

	round := NewRound(fmt.Sprintf("%d", ev.Epoch))
	if err := s.Load(round); err != nil {
		return err
	}

	if !round.Exists() {
		round.Epoch = entity.NewInt(ev.Epoch)
		// FIXME: Previous is `string`, but code is trying to access `Int`, will need to check subgraph graph
		if ev.Epoch.Uint64() != 0 {
			round.Previous = S(bi().Sub(ev.Epoch, big.NewInt(1)).String())
		}

		round.StartAt = entity.NewIntFromLiteral(s.Block().Timestamp().Unix())
		round.StartBlock = entity.NewIntFromLiteral(int64(s.Block().Number()))
		round.StartHash = trace.Hash

		if err := s.Save(round); err != nil {
			return err
		}
	}

	market.Epoch = S(round.ID)
	market.Paused = false
	if err := s.Save(market); err != nil {
		return err
	}
	return nil
}
