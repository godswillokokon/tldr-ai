package usage

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"tldr-ai-be/internal/errs"
)

// defaultPerCallUSD is used when USAGE_PER_CALL_USD is unset; 0 in env falls back to this.
const defaultPerCallUSD = 0.01

// Budget tracks monthly estimated USD and optional lifetime call caps, with
// reserve/commit for in-flight requests. State is in-memory (single process).
type Budget struct {
	mu sync.Mutex

	perCall  float64
	monthCap float64
	maxCalls int64
	month    time.Time

	monthSpent float64
	pending    map[uint64]float64
	nextID     uint64
	usedTotal  int64
}

// Reservation is an opaque token; create only via TryReserve.
type Reservation struct {
	id uint64
	b  *Budget
}

// NewFromEnv loads USAGE_BUDGET_USD (0 = no monthly $ cap), USAGE_PER_CALL_USD, USAGE_MAX_CALLS (0 = no call cap).
func NewFromEnv(get func(key string) string) *Budget {
	b := &Budget{pending: make(map[uint64]float64), month: monthStartUTC(time.Now().UTC())}
	b.perCall = defaultPerCallUSD
	if s := strings.TrimSpace(get("USAGE_PER_CALL_USD")); s != "" {
		if v, e := strconv.ParseFloat(s, 64); e == nil && v > 0 {
			b.perCall = v
		}
	}
	if s := strings.TrimSpace(get("USAGE_BUDGET_USD")); s != "" {
		if v, e := strconv.ParseFloat(s, 64); e == nil && v >= 0 {
			b.monthCap = v
		}
	}
	if s := strings.TrimSpace(get("USAGE_MAX_CALLS")); s != "" {
		if v, e := strconv.ParseInt(s, 10, 64); e == nil && v >= 0 {
			b.maxCalls = v
		}
	}
	return b
}

func monthStartUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func (b *Budget) now() time.Time { return time.Now().UTC() }

func (b *Budget) maybeRolloverLocked(now time.Time) {
	cur := monthStartUTC(now)
	if b.month.IsZero() {
		b.month = cur
		return
	}
	if cur.Equal(b.month) {
		return
	}
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

// TryReserve reserves the per-call estimate for the current UTC month.
func (b *Budget) TryReserve() (*Reservation, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.maybeRolloverLocked(b.now())
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

// Release rolls back a reservation (any failure after TryReserve, including decode).
func (b *Budget) Release(r *Reservation) {
	if b == nil || r == nil || r.b == nil || r.b != b {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.pending[r.id]; !ok {
		return
	}
	delete(b.pending, r.id)
}

// Commit records a successful call: applies reserved $ to the month and increments use count.
func (b *Budget) Commit(r *Reservation) {
	if b == nil || r == nil || r.b == nil || r.b != b {
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

// AdminReset clears all counters and pending state (this month, lifetime use).
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

// Snapshot is returned by GET /api/usage.
type Snapshot struct {
	Used         int64    `json:"used"`
	Cap          int64    `json:"cap"`
	Unlimited    bool     `json:"unlimited"`
	Remaining    *int64   `json:"remaining,omitempty"`
	SpentUSD     float64  `json:"spentUsd"`
	BudgetUSD    float64  `json:"budgetUsd"`
	ReservedUSD  float64  `json:"reservedUsd"`
	PerCallUSD   float64  `json:"perCallUsd"`
	UnlimitedUSD bool     `json:"unlimitedUsd"`
	RemainingUSD *float64 `json:"remainingUsd,omitempty"`
	MonthUTC     string   `json:"monthUtc"`
}

// Snapshot returns a consistent point-in-time view.
func (b *Budget) Snapshot() Snapshot {
	if b == nil {
		now := time.Now().UTC()
		return Snapshot{Unlimited: true, UnlimitedUSD: true, MonthUTC: now.Format("2006-01")}
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
		Used:         used,
		Cap:          cap,
		PerCallUSD:   b.perCall,
		SpentUSD:     spent,
		BudgetUSD:    budget,
		ReservedUSD:  resv,
		MonthUTC:     mon.Format("2006-01"),
		Unlimited:    cap == 0,
		UnlimitedUSD: budget <= 0,
	}
	if cap == 0 {
		s.Remaining = nil
	} else {
		rem := cap - used - int64(len(b.pending))
		if rem < 0 {
			rem = 0
		}
		s.Remaining = &rem
	}
	if budget > 0 {
		left := budget - spent - resv
		if left < 0 {
			left = 0
		}
		rounded := round2(left)
		s.RemainingUSD = &rounded
	}
	return s
}

func round2(f float64) float64 { return float64(int64(f*100+0.5)) / 100 }
