package domain

// Request / body size and field limits for the process API.
const (
	MaxRequestBodyBytes   = 1 << 20 // 1 MiB
	MaxInputTextBytes     = 256_000
	MinInputTextRunes     = 20
	MaxSummaryRunes       = 4_000
	MaxActionItemRunes    = 500
	MaxModelResponseBytes = 512 << 10 // 512 KiB
	MaxModelJSONBytes     = 512 << 10 // 512 KiB
)
