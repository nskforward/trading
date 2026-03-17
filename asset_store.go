package trading

import "github.com/nskforward/trading/types"

type AssetStore struct {
	broker types.Broker
	items  map[string]types.Asset
}

func NewAssetStore(broker types.Broker) *AssetStore {
	return &AssetStore{
		broker: broker,
		items:  make(map[string]types.Asset),
	}
}

func (store *AssetStore) Get(symbol string) (types.Asset, error) {
	v, ok := store.items[symbol]
	if ok {
		return v, nil
	}

	v, err := store.broker.GetAsset(symbol)
	if err != nil {
		return v, err
	}

	store.items[symbol] = v

	return v, nil
}
