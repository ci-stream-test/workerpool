/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package workerpool

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWPool(t *testing.T) {
	maxIdle := 1
	maxIdleTime := time.Millisecond * 500
	p := New(maxIdle, maxIdleTime)
	var (
		sum  int32
		wg   sync.WaitGroup
		size int32 = 100
	)
	require.Zero(t, p.Size())
	for i := int32(0); i < size; i++ {
		wg.Add(1)
		p.Go(func() {
			defer wg.Done()
			atomic.AddInt32(&sum, 1)
		})
	}
	require.NotZero(t, p.Size())

	wg.Wait()
	require.Equal(t, size, atomic.LoadInt32(&sum))
	for p.Size() != int32(maxIdle) { // waiting for workers finished and idle workers left
		runtime.Gosched()
	}
	for p.Size() > 0 { // waiting for idle workers timeout
		time.Sleep(maxIdleTime)
	}
	require.Zero(t, p.Size())
}
