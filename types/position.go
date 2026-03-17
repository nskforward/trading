package types

import (
	"math"
	"time"
)

type Position struct {
	Updated time.Time
	Symbol  string
	Size    float64
	Price   float64
}

func (p *Position) Merge(tradePrice, tradeSize float64) {
	if tradeSize == 0 {
		return // ничего не делаем
	}

	// Если позиции нет, просто открываем новую
	if p.Size == 0 {
		p.Size = tradeSize
		p.Price = tradePrice
		return
	}

	// 1. Увеличение существующей позиции (знаки одинаковые)
	if (p.Size > 0 && tradeSize > 0) || (p.Size < 0 && tradeSize < 0) {
		totalSize := p.Size + tradeSize
		// Средневзвешенная цена
		avgPrice := (p.Size*p.Price + tradeSize*tradePrice) / totalSize
		p.Price = math.Round(avgPrice*100) / 100
		p.Size = totalSize
		return
	}

	// 2. Знаки противоположные — уменьшение позиции или переворот
	if math.Abs(tradeSize) <= math.Abs(p.Size) {
		// Частичное или полное закрытие: средняя цена не меняется
		p.Size += tradeSize // tradeSize противоположного знака
		if p.Size == 0 {
			p.Price = 0 // опционально обнуляем цену
		}
		return
	}

	// Переворот: полностью закрываем текущую позицию и открываем новую в противоположном направлении
	// Новая позиция = остаток от сделки (имеет знак tradeSize)
	newSize := tradeSize + p.Size // tradeSize > |p.Size| по модулю
	p.Size = newSize
	p.Price = tradePrice
}
