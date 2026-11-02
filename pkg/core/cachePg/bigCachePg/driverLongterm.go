package bigCachePg

/*
 * Copyright 2025 The Go-Spring Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	StarterBigCache "go-spring.org/starter-bigcache"
	"go-spring.org/stdlib/errutil"
)

func init() {
	// A third Driver bean, selected per instance with
	// spring.bigcache.instances.<name>.driver=longterm.
	gs.Provide(func() StarterBigCache.Driver {
		return LongTermBigCacheDriver{}
	}).Name("longterm").Caller(1)
}

const (
	// A long-term instance must outlive any operational window the caller may
	// reasonably expect; anything shorter is a regular cache, not a long-term
	// one.
	longTermMinLifeWindow = 24 * time.Hour

	// Without age cleanup, the hard cap is the only memory bound. Requiring one
	// keeps a long-term instance from growing without limit.
	longTermMinHardMaxCacheSize = 1
)

// LongTermBigCacheDriver is a custom Driver for data that should stay in the
// cache for a long time.
//
// It deliberately does not add a per-entry deadline. Instead it rejects a
// life-window shorter than 24h, requires clean-window=0 so the background
// cleaner never removes an entry by age, and requires a hard cap so the
// process-local cache stays memory-bounded. Entries therefore leave only
// through capacity eviction (or when the process exits).
type LongTermBigCacheDriver struct{}

// CreateClient builds one long-term instance. The native cache is used as-is:
// long-term is a configuration policy here, not an extra value envelope, so
// values written through this driver are stored without a TTL header.
func (LongTermBigCacheDriver) CreateClient(ctx context.Context, name string,
	c StarterBigCache.Config) (*StarterBigCache.Cache, error) {

	log.Info(ctx, log.TagAppDef, "LongTermBigCacheDriver::CreateClient",
		log.String("name", name),
		log.String("life-window", c.LifeWindow.String()),
		log.Int("hard-max-cache-size", c.HardMaxCacheSize))

	if err := validateLongTermConfig(name, c); err != nil {
		return nil, err
	}

	client, err := bigcache.New(ctx, newBigCacheConfig(c))
	if err != nil {
		return nil, errutil.Explain(err, "bigcache: create instance %q", name)
	}

	cache, err := StarterBigCache.NewCache(client, name)
	if err != nil {
		_ = client.Close()
		return nil, err
	}
	return cache, nil
}

// validateLongTermConfig makes the long-term policy explicit, so a typo in an
// instance's properties fails at startup instead of silently producing a cache
// that expires entries or grows without bound.
func validateLongTermConfig(name string, c StarterBigCache.Config) error {
	if c.LifeWindow < longTermMinLifeWindow {
		return errutil.Explain(nil,
			"bigcache: create instance %q: long-term life-window %s is below %s",
			name, c.LifeWindow, longTermMinLifeWindow)
	}
	if c.CleanWindow != 0 {
		return errutil.Explain(nil,
			"bigcache: create instance %q: long-term requires clean-window=0, got %s",
			name, c.CleanWindow)
	}
	if c.HardMaxCacheSize < longTermMinHardMaxCacheSize {
		return errutil.Explain(nil,
			"bigcache: create instance %q: long-term requires hard-max-cache-size >= %d MB, got %d",
			name, longTermMinHardMaxCacheSize, c.HardMaxCacheSize)
	}
	return nil
}
