package main

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"time"

	"golang.org/x/sync/errgroup"
)

func monitorGoroutines(ctx context.Context, prevGoroutines int) {
	prev := prevGoroutines

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Мониторинг завершен.")
			return
		default:
			current := runtime.NumGoroutine()
			fmt.Println("Текущее количество горутин:", current)
			if prev > 0 {
				diff := math.Abs(float64(current-prev)) / float64(prev) * 100
				if diff >= 20 {
					if current > prev {
						fmt.Println("⚠️ Предупреждение: Количество горутин увеличилось более чем на 20%!")
					} else {
						fmt.Println("⚠️ Предупреждение: Количество горутин уменьшилось более чем на 20%!")
					}
				}
			}
			prev = current
			time.Sleep(300 * time.Millisecond)
		}
	}
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	// Мониторинг горутин
	go func() error {
		monitorGoroutines(ctx, runtime.NumGoroutine())
		return nil
	}()

	// Имитация активной работы приложения с созданием горутин
	for i := 0; i < 64; i++ {
		idx := i
		g.Go(func() error {
			time.Sleep(5 * time.Second)
			if idx != 0 && idx%10 == 0 {
				return fmt.Errorf("ошибка в горутине #%d", idx)
			}
			return nil
		})
		time.Sleep(80 * time.Millisecond)
	}

	// Ожидание завершения всех горутин
	if err := g.Wait(); err != nil {
		fmt.Println("Ошибка:", err)
	}
}
