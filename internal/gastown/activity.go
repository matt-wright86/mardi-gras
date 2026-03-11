package gastown

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Event represents a single entry from the Gas Town event log (.events.jsonl).
type Event struct {
	Timestamp  string          `json:"ts"`
	Source     string          `json:"source"`
	Type       string          `json:"type"`
	Actor      string          `json:"actor"`
	Payload    json.RawMessage `json:"payload"`
	Visibility string          `json:"visibility"`
}

type recentEventCache struct {
	mu      sync.Mutex
	path    string
	limit   int
	size    int64
	modTime time.Time
	ring    []Event
	count   int
}

var cachedRecentEvents recentEventCache

// EventsPath returns the path to the Gas Town events log.
// Uses GT_HOME if set, otherwise defaults to ~/gt/.events.jsonl.
func EventsPath() string {
	home := os.Getenv("GT_HOME")
	if home == "" {
		userHome, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		home = filepath.Join(userHome, "gt")
	}
	return filepath.Join(home, ".events.jsonl")
}

// LoadRecentEvents reads the event log and returns the last `limit` events
// in reverse chronological order (newest first).
// Uses a ring buffer so memory usage is O(limit), not O(file_size).
// Returns nil, nil if the file does not exist or is empty.
func LoadRecentEvents(path string, limit int) ([]Event, error) {
	if limit <= 0 || path == "" {
		return nil, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	cachedRecentEvents.mu.Lock()
	defer cachedRecentEvents.mu.Unlock()

	if cachedRecentEvents.canReuse(path, limit, info) {
		return materializeRecentEvents(cachedRecentEvents.ring, cachedRecentEvents.count, limit), nil
	}

	if cachedRecentEvents.canAppend(path, limit, info) {
		if err := cachedRecentEvents.append(path, info); err == nil {
			return materializeRecentEvents(cachedRecentEvents.ring, cachedRecentEvents.count, limit), nil
		}
	}

	ring, count, err := loadRecentEventRing(path, limit)
	if err != nil {
		return nil, err
	}
	cachedRecentEvents.store(path, limit, info, ring, count)
	return materializeRecentEvents(ring, count, limit), nil
}

func (c *recentEventCache) canReuse(path string, limit int, info os.FileInfo) bool {
	return c.path == path &&
		c.limit == limit &&
		c.size == info.Size() &&
		c.modTime.Equal(info.ModTime())
}

func (c *recentEventCache) canAppend(path string, limit int, info os.FileInfo) bool {
	return c.path == path &&
		c.limit == limit &&
		c.size >= 0 &&
		info.Size() > c.size &&
		len(c.ring) == limit
}

func (c *recentEventCache) append(path string, info os.FileInfo) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Seek(c.size, io.SeekStart); err != nil {
		return err
	}

	count, err := scanRecentEvents(f, c.ring, c.count)
	if err != nil {
		return err
	}
	c.count = count
	c.size = info.Size()
	c.modTime = info.ModTime()
	return nil
}

func (c *recentEventCache) store(path string, limit int, info os.FileInfo, ring []Event, count int) {
	c.path = path
	c.limit = limit
	c.size = info.Size()
	c.modTime = info.ModTime()
	c.ring = ring
	c.count = count
}

func loadRecentEventRing(path string, limit int) ([]Event, int, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	defer f.Close()

	ring := make([]Event, limit)
	count, err := scanRecentEvents(f, ring, 0)
	if err != nil {
		return nil, 0, err
	}
	return ring, count, nil
}

func scanRecentEvents(r io.Reader, ring []Event, count int) (int, error) {
	scanner := bufio.NewScanner(r)
	// Allow long lines (some payloads can be large)
	scanner.Buffer(make([]byte, 0, 64*1024), 256*1024)
	limit := len(ring)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			continue
		}
		ring[count%limit] = ev
		count++
	}
	if err := scanner.Err(); err != nil {
		return count, err
	}
	return count, nil
}

func materializeRecentEvents(ring []Event, count, limit int) []Event {
	if count == 0 {
		return nil
	}

	n := min(count, limit)
	result := make([]Event, n)
	start := count - n
	for i := 0; i < n; i++ {
		result[n-1-i] = ring[(start+i)%limit]
	}
	return result
}

// AgentActivityHistogram builds a per-agent event count histogram from recent events.
// Returns a map of agent name → time-bucketed event counts (oldest to newest).
// `buckets` is the number of time windows; `window` is the total time span.
func AgentActivityHistogram(events []Event, buckets int, window time.Duration) map[string][]int {
	if len(events) == 0 || buckets <= 0 {
		return nil
	}

	now := time.Now()
	bucketDur := window / time.Duration(buckets)
	result := make(map[string][]int)

	for _, ev := range events {
		t, err := time.Parse(time.RFC3339, ev.Timestamp)
		if err != nil {
			continue
		}
		age := now.Sub(t)
		if age < 0 || age > window {
			continue
		}
		// Bucket index: 0 = oldest, buckets-1 = most recent
		idx := buckets - 1 - int(age/bucketDur)
		if idx < 0 {
			idx = 0
		}
		if idx >= buckets {
			idx = buckets - 1
		}

		actor := ev.Actor
		if _, ok := result[actor]; !ok {
			result[actor] = make([]int, buckets)
		}
		result[actor][idx]++
	}

	return result
}

// AgentEventCount returns the total event count per agent from a set of events.
func AgentEventCount(events []Event) map[string]int {
	counts := make(map[string]int)
	for _, ev := range events {
		counts[ev.Actor]++
	}
	return counts
}

// EventPayloadString extracts a string field from the event payload.
func EventPayloadString(ev Event, key string) string {
	if len(ev.Payload) == 0 {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(ev.Payload, &m); err != nil {
		return ""
	}
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
