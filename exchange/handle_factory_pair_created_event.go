package exchange

import (
	"fmt"

	eth "github.com/streamingfast/eth-go"
	"github.com/streamingfast/sparkle/entity"

	"go.uber.org/zap"
)

func validateToken(tok *eth.Token) bool {
	return !tok.IsEmptyDecimal
}

func (s *Subgraph) HandleFactoryPairCreatedEvent(ev *FactoryPairCreatedEvent) error {
	factory := NewPancakeFactory(FactoryAddress)
	err := s.Load(factory)
	if err != nil {
		return err
	}

	if !factory.Exists() {
		bundle := NewBundle("1")
		if err := s.Save(bundle); err != nil {
			return err
		}
	}

	factory.TotalPairs = entity.IntAdd(factory.TotalPairs, IL(1))
	if err := s.Save(factory); err != nil {
		return err
	}

	token0 := NewToken(ev.Token0.Pretty())
	err = s.Load(token0)
	if err != nil {
		return err
	}

	if !token0.Exists() {
		tm, valid := s.GetTokenInfo(ev.Token0, validateToken)
		if !valid {
			s.Log.Info("token 0 is invalid, skipping creating pair",
				zap.String("token0", ev.Token0.Pretty()),
				zap.String("pair", factory.ID),
				zap.Uint64("block_number", s.Block().Number()),
				zap.String("block_id", s.Block().ID()),
				zap.String("transaction_hash", ev.Transaction.Hash.Pretty()),
			)
			return nil
		}

		token0 = NewToken(ev.Token0.Pretty())
		token0.Name = tm.Name
		token0.Symbol = tm.Symbol
		token0.Decimals = IL(int64(tm.Decimals))
		token0.DerivedBNB = FL(0).Ptr()
		token0.DerivedUSD = FL(0).Ptr()

		if err = s.Save(token0); err != nil {
			return fmt.Errorf("saving initialToken 0: %w", err)
		}
	}

	token1 := NewToken(ev.Token1.Pretty())
	err = s.Load(token1)
	if err != nil {
		return fmt.Errorf("loading initialToken 1: %w", err)
	}

	if !token1.Exists() {
		tm, valid := s.GetTokenInfo(ev.Token1, validateToken)
		if !valid {
			s.Log.Info("token 1 is invalid, skipping creating pair",
				zap.String("token1", ev.Token1.Pretty()),
				zap.String("pair", factory.ID),
				zap.Uint64("block_number", s.Block().Number()),
				zap.String("block_id", s.Block().ID()),
				zap.String("transaction_hash", ev.Transaction.Hash.Pretty()),
			)
			return nil
		}

		token1 = NewToken(ev.Token1.Pretty())
		token1.Symbol = tm.Symbol
		token1.Name = tm.Name
		token1.Decimals = IL(int64(tm.Decimals))
		token1.DerivedBNB = FL(0).Ptr()
		token1.DerivedUSD = FL(0).Ptr()

		if err = s.Save(token1); err != nil {
			return fmt.Errorf("saving initialToken 1: %w", err)
		}
	}

	pair := NewPair(ev.Pair.Pretty())
	pair.Token0 = token0.ID
	pair.Token1 = token1.ID
	pair.Block = entity.NewIntFromLiteralUnsigned(s.Block().Number())
	pair.Timestamp = entity.NewIntFromLiteral(s.Block().Timestamp().Unix())
	pair.Name = fmt.Sprintf("%s-%s", token0.Symbol, token1.Symbol)

	err = s.Save(pair)
	if err != nil {
		return fmt.Errorf("saving pair: %w", err)
	}

	err = s.CreatePairTemplateWithTokens(ev.Pair, ev.Token0, ev.Token1)
	if err != nil {
		return fmt.Errorf("creating pair template: %w", err)
	}

	return nil
}
