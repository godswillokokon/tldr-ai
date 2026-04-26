package usage

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"tldr-ai-be/internal/errs"
)

// Defaults for env-backed limits (0 = unlimited for caps).
const (
	defaultPerCallUSD = 0.01
)

// Budget tracks monthly estimated USD and optional lifetime call caps, with
// reservation/ commit for in-flight requests. Counters are in-memory (single process).
type Budget struct {
	mu sync.Mutex

	perCall   float64
	monthCap  float64 // 0 = unlimited Usd
	maxCalls  int64   // 0 = unlimited
	month     time.Time
	monthSpent float64
	pending   map[uint64]float64
	nextID    uint64
	usedTotal int64
}

// Reservation is an opaque token for Release/Commit; do not create manually.
type Reservation struct {
	id uint64
	b  *Budget
}

// NewFromEnv loads USAGE_BUDGET_USD (monthly, 0=unlimited), USAGE_PER_CALL_USD, USAGE_MAX_CALLS (lifetime, 0=unlimited).
func NewFromEnv(get func(key string) string) *Budget {
	b := &Budget{pending: make(map[uint64]float64), month: monthStartUTC(time.Now().UTC())}
	b.perCall = defaultPerCallUSD
	if s := strings.TrimSpace(get("USAGE_PER_CALL_USD")); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil && v >= 0 {
			if v == 0 {
				b.perCall = defaultPerCallUSD
			} else {
				b.perCall = v
			}
		}
	}
	if s := strings.TrimSpace(get("USAGE_BUDGET_USD")); s != "" {
		if v, err := strconv.ParseFloat(s, 64); err == nil && v >= 0 {
			b.monthCap = v
		}
	}
	if s := strings.TrimSpace(get("USAGE_MAX_CALLS")); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil && v >= 0 {
			b.maxCalls = v
		}
	}
	// Allow tests to override "now" by optional hook (unused in production).
	_ = get("USAGE_TEST_TIME")
	_ = os.DevNull
	return b
}

func monthStartUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func (b *Budget) now() time.Time { return time.Now().UTC() }

// maybeRollover locked: new UTC month only resets committed spend for the month, not in-flight reserved USD.
func (b *Budget) maybeRolloverLocked(now time.Time) {
	cur := monthStartUTC(now)
	if b.month.IsZero() {
		b.month = cur
		return
	}
	if cur.Equal(b.month) {
		return
	}
	// Rollover: committed month spend resets; pending reservations count toward the new month’s budget
	// (reserved USD is still in-flight, so we keep b.pending and monthReserved is derived from it).
	if cur.After(b.month) {
		b.month = cur
		b.monthSpent = 0
	}
}

func monthReservedLocked(b *Budget) float64 {
	var s float64
	for _, v := range b.pending {
		s += v
	}
	return s
}

// TryReserve reserves one per-call estimate against the current UTC month. On failure, returns
// *errs.UsageCapExceeded.
func (b *Budget) TryReserve() (*Reservation, error) {
	if b == nil {
		return &Reservation{}, nil
	}
	now := b.now()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.maybeRolloverLocked(now)
	if b.maxCalls > 0 && b.usedTotal+int64(len(b.pending)) >= b.maxCalls {
		return nil, errs.UsageCapExceeded("Call limit reached for this service")
	}
	resv := monthReservedLocked(b)
	if b.monthCap > 0 && b.monthSpent+resv+b.perCall > b.monthCap+1e-9 {
		return nil, errs.UsageCapExceeded("Monthly usage budget reached")
	}
	b.nextID++
	id := b.nextID
	b.pending[id] = b.perCall
	return &Reservation{id: id, b: b}, nil
}

// Release rolls back a reservation (failure path after AI or decode errors).
func (b *Budget) Release(r *Reservation) {
	if b == nil || r == nil || r.b == nil {
		return
	}
	if r.b != b {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	amt, ok := b.pending[r.id]
	if !ok {
		return
	}
	delete(b.pending, r.id)
	_ = amt
}

// Commit finalizes a successful call: moves reserved amount into this month’s spent and increments lifetime use count.
func (b *Budget) Commit(r *Reservation) {
	if b == nil || r == nil || r.b == nil {
		return
	}
	if r.b != b {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	amt, ok := b.pending[r.id]
	if !ok {
		return
	}
	delete(b.pending, r.id)
	b.monthSpent += amt
	b.usedTotal++
}

// AdminReset zeroes in-memory counters (all-time calls, this month’s spend) and pending reservations.
func (b *Budget) AdminReset() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.month = monthStartUTC(b.now())
	b.monthSpent = 0
	b.pending = make(map[uint64]float64)
	b.usedTotal = 0
}

// Snapshot is a point-in-time view for JSON APIs (UTC).
type Snapshot struct {
	Used   int64   `json:"used"`
	Cap    int64   `json:"cap"`
	Unlimited   bool  `json:"unlimited"`
	Remaining   *int64 `json:"remaining,omitempty"`

	SpentUSD    float64  `json:"spentUsd"`
	BudgetUSD   float64  `json:"budgetUsd"`
	ReservedUSD float64  `json:"reservedUsd"`
	PerCallUSD  float64  `json:"perCallUsd"`
	UnlimitedUSD bool    `json:"unlimitedUsd"`
	RemainingUSD *float64 `json:"remainingUsd,omitempty"`
	MonthUTC    string   `json:"monthUtc"`
}

// Snapshot returns current usage; safe to call from handlers.
func (b *Budget) Snapshot() Snapshot {
	if b == nil {
		return Snapshot{Unlimited: true, UnlimitedUSD: true, MonthUTC: time.Now().UTC().Format("2006-01")}
	}
	now := b.now()
	b.mu.Lock()
	defer b.mu.Unlock()
	b.maybeRolloverLocked(now)
	resv := monthReservedLocked(b)
	used := b.usedTotal
	cap := b.maxCalls
	budget := b.monthCap
	spent := b.monthSpent
	mon := b.month
	s := Snapshot{
		Used:        used,
		Cap:         cap,
		PerCallUSD:  b.perCall,
		SpentUSD:    spent,
		BudgetUSD:   budget,
		ReservedUSD: resv,
		MonthUTC:    mon.Format("2006-01"),
	}
	s.Unlimited = cap == 0
	if cap == 0 {
		s.Remaining = nil
	} else {
		rem := cap - used - int64(len(b.pending))
		if rem < 0 {
			rem = 0
		}
		s.Remaining = &rem
	}
	s.UnlimitedUSD = budget <= 0
	if budget > 0 {
		left := budget - spent - resv
		if left < 0 {
			left = 0
		}
		rounded := round2(left)
		s.RemainingUSD = &rounded
	}
	// legacy-friendly alias for clients expecting "unlimited" for usd: mirror unlimitedUsd
	_ = fmt.Sprint(used) // no-op; avoid unused
	return s
}

func round2(f float64) float64 { return float64(int64(f*100+0.5)) / 100 }
