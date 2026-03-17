package types

import (
	"math"
)

type Position struct {
	Symbol string
	Size   float64
	Price  float64
}

func (p *Position) Merge(trade Trade) {
	if trade.Size == 0 {
		return
	}

	if p.Size == 0 {
		p.Size = trade.Size
		p.Price = trade.Price
		return
	}

	// 1. Увеличение существующей позиции (знаки одинаковые)
	if (p.Size > 0 && trade.Size > 0) || (p.Size < 0 && trade.Size < 0) {
		totalSize := p.Size + trade.Size
		// Средневзвешенная цена
		avgPrice := (p.Size*p.Price + trade.Size*trade.Price) / totalSize
		p.Price = math.Round(avgPrice*100) / 100
		p.Size = totalSize
		return
	}

	// 2. Знаки противоположные — уменьшение позиции или переворот
	if math.Abs(trade.Size) <= math.Abs(p.Size) {
		// Частичное или полное закрытие: средняя цена не меняется
		p.Size += trade.Size // tradeSize противоположного знака
		if p.Size == 0 {
			p.Price = 0 // опционально обнуляем цену
		}
		return
	}

	// Переворот: полностью закрываем текущую позицию и открываем новую в противоположном направлении
	// Новая позиция = остаток от сделки (имеет знак tradeSize)
	p.Size = p.Size + trade.Size
	p.Price = trade.Price
}
