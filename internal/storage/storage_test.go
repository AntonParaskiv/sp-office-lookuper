package storage

import (
	"math/rand"
	"sync"
	"testing"
)

func TestStorage_SetOfficeSortPoint(t *testing.T) {
	workerPool := 1000
	workerJobCount := 1000

	t.Run("concurent write", func(t *testing.T) {
		s := NewStorage()

		wg := &sync.WaitGroup{}
		wg.Add(workerPool)

		for w := 0; w < workerPool; w++ {
			go func() {
				defer wg.Done()
				for j := 0; j < workerJobCount; j++ {
					officeID := int64(rand.Int() % 1000)
					sortPointID := int64(rand.Int() % 1000)
					s.SetOfficeSortPoint(officeID, sortPointID)
				}
			}()
		}

		wg.Wait()
	})
}

func TestStorage_GetSortPoint(t *testing.T) {
	workerPool := 1000
	workerJobCount := 1000

	t.Run("concurent write/read", func(t *testing.T) {
		s := NewStorage()

		wg := &sync.WaitGroup{}

		// write
		wg.Add(workerPool)
		for w := 0; w < workerPool; w++ {
			go func() {
				defer wg.Done()
				for j := 0; j < workerJobCount; j++ {
					officeID := int64(rand.Int() % 1000)
					sortPointID := int64(rand.Int() % 1000)
					s.SetOfficeSortPoint(officeID, sortPointID)
				}
			}()
		}

		// read
		wg.Add(workerPool)
		for w := 0; w < workerPool; w++ {
			go func() {
				defer wg.Done()
				for j := 0; j < workerJobCount; j++ {
					officeID := int64(rand.Int() % 1000)
					_, _ = s.GetSortPoint(officeID)
				}
			}()
		}

		wg.Wait()

	})
}
