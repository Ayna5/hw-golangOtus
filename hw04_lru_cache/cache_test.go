package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("capacity = 0 should be return error", func(t *testing.T) {
		c, err := NewCache(0)
		require.EqualError(t, err, ErrCache)
		require.Nil(t, c)
	})

	t.Run("empty cache", func(t *testing.T) {
		c, err := NewCache(10)
		require.NoError(t, err, "cannot get NewCache")

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c, err := NewCache(5)
		require.NoError(t, err, "cannot get NewCache")

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c, err := NewCache(3)
		require.NoError(t, err, "cannot get NewCache")

		wasInCache := c.Set("aaa", 1)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 2)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 3)
		require.False(t, wasInCache)
		wasInCache = c.Set("ddd", 4)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("pushing out rarely used elements", func(t *testing.T) {
		c, err := NewCache(3)
		require.NoError(t, err, "cannot get NewCache")

		wasInCache := c.Set("aaa", 1)
		require.False(t, wasInCache)
		wasInCache = c.Set("bbb", 2)
		require.False(t, wasInCache)
		wasInCache = c.Set("ccc", 3)
		require.False(t, wasInCache)
		c.Get("aaa")
		c.Set("ccc", 6)
		c.Set("aaa", 25)
		c.Get("aaa")
		c.Get("ccc")

		wasInCache = c.Set("ddd", 444)
		require.False(t, wasInCache)

		val, ok := c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	//t.Skip() // NeedRemove if task with asterisk completed

	c, err := NewCache(10)
	require.NoError(t, err, "cannot get NewCache")

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
