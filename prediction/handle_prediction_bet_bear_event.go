package prediction

import (
	"encoding/binary"
	"fmt"

	"github.com/streamingfast/sparkle/entity"
	pbcodec "github.com/streamingfast/sparkle/pb/sf/ethereum/codec/v1"
)

func (s *Subgraph) HandlePredictionBetBearEvent(trace *pbcodec.TransactionTrace, ev *PredictionBetBearEvent) error {
	if s.StepBelow(2) {
		return nil
	}
	round := NewRound(ev.CurrentEpoch.String())
	if err := s.Load(round); err != nil {
		return err
	}
	if !round.Exists() {
		return fmt.Errorf("Tried to bet (bear) without an existing round (epoch: %s).", ev.CurrentEpoch)
	}

	round.TotalBets = entity.IntAdd(round.TotalBets, IL(1))
	round.TotalAmount = entity.FloatAdd(round.TotalAmount, F(bf().Quo(
		bf().SetInt(ev.Amount),
		EIGHTEEN_BI,
	)))

	round.BearBets = entity.IntAdd(round.BearBets, IL(1))
	round.BearAmount = entity.FloatAdd(round.BearAmount, F(bf().Quo(
		bf().SetInt(ev.Amount),
		EIGHTEEN_BI,
	)))

	if err := s.Save(round); err != nil {
		return err
	}

	user := NewUser(ev.Sender.Pretty())
	if err := s.Load(user); err != nil {
		return err
	}
	if !user.Exists() {
		user.Address = entity.Bytes(ev.Sender)
		user.CreatedAt = entity.NewIntFromLiteral(s.Block().Timestamp().Unix()) ///uint64(s.Block().Timestamp().Unix())
		user.Block = entity.NewIntFromLiteral(int64(s.Block().Number()))        ///s.Block().Number()
	}
	user.UpdatedAt = entity.NewIntFromLiteral(s.Block().Timestamp().Unix()) ///uint64(s.Block().Timestamp().Unix())

	user.TotalBets = entity.IntAdd(user.TotalBets, IL(1))
	user.TotalBNB = entity.FloatAdd(user.TotalBNB, F(bf().Quo(
		bf().SetInt(ev.Amount),
		EIGHTEEN_BI,
	)))

	// user.totalBNB = user.totalBNB.plus(event.params.amount.divDecimal(EIGHTEEN_BD));

	if err := s.Save(user); err != nil {
		return err
	}

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(ev.CurrentEpoch.Uint64()))
	betID := fmt.Sprintf("%s%x", ev.Sender.Pretty(), buf)
	bet := NewBet(betID)
	bet.Round = round.ID
	bet.User = user.ID
	bet.Hash = trace.Hash
	bet.Amount = F(bf().Quo(
		bf().SetInt(ev.Amount),
		EIGHTEEN_BI,
	))

	bet.Position = PositionBear
	bet.Claimed = false

	if err := s.Save(bet); err != nil {
		return err
	}

	return nil
}
