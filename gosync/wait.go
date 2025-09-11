package gosync

import (
	"context"
	"sync"
)

// WaitGroup обертка над sync.WaitGroup, позволяющая вызвать Wait с контекстом
type WaitGroup struct {
	wg *sync.WaitGroup
}

func NewWaitGroup(wg *sync.WaitGroup) *WaitGroup {
	if wg == nil {
		wg = &sync.WaitGroup{}
	}

	return &WaitGroup{
		wg: wg,
	}
}

// Add добавляет delta к счетчику. Полностью идентичен методу из sync.WaitGroup
func (w *WaitGroup) Add(delta int) {
	w.wg.Add(delta)
}

// Done уменьшает счетчик на 1. Полностью идентичен методу из sync.WaitGroup
func (w *WaitGroup) Done() {
	w.wg.Done()
}

// Wait блокируется, пока счетчик не достигнет 0. Полностью идентичен методу из sync.WaitGroup
func (w *WaitGroup) Wait() {
	w.wg.Wait()
}

// WaitContext блокируется, пока счетчик не достигнет 0 или не завершится контекст
// Если контекст завершен до обнуления счетчика - вернет ошибку, с которой завершился сам контекст
func (w *WaitGroup) WaitContext(ctx context.Context) error {
	return WaitContext(ctx, NopWaitFn(w.Wait))
}
