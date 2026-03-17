package finam

import "github.com/nskforward/trading/types"

type quoteCache struct {
	quotes map[string]*types.Quote
}

func newQuoteCache() *quoteCache {
	return &quoteCache{
		quotes: make(map[string]*types.Quote),
	}
}

func (s *quoteCache) Get(in types.Quote) *types.Quote {
	cached, ok := s.quotes[in.Symbol]
	if !ok {
		cached = &in
		s.quotes[in.Symbol] = cached
		return cached
	}

	if in.Ask == 0 {
		in.Ask = cached.Ask
	}

	if in.Bid == 0 {
		in.Bid = cached.Bid
	}

	if in.Bid != cached.Bid || in.Ask != cached.Ask {
		cached.Ask = in.Ask
		cached.Bid = in.Bid
		return cached
	}

	return nil
}
