package dash

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type statsResponse struct {
	UptimeSeconds float64    `json:"uptime_seconds"`
	Session       statsBlock `json:"session"`
	AllTime       statsBlock `json:"all_time"`
}

type statsBlock struct {
	Requests          int64   `json:"requests"`
	RequestBytes      int64   `json:"request_bytes"`
	ResponseOrigBytes int64   `json:"response_original_bytes"`
	ResponseWireBytes int64   `json:"response_wire_bytes"`
	SavingsBytes      int64   `json:"savings_bytes"`
	SavingsPct        float64 `json:"savings_pct"`
}

func (d *Dashboard) handleStats(w http.ResponseWriter, r *http.Request) {
	s := d.source.Stats()
	if s == nil {
		http.Error(w, "stats not available", http.StatusServiceUnavailable)
		return
	}

	resp := statsResponse{
		UptimeSeconds: d.source.Uptime().Seconds(),
		Session:       buildStatsBlock(s.SessionRequests(), s.SessionRequestOriginalBytes(), s.SessionResponseOriginalBytes(), s.SessionResponseWireBytes()),
		AllTime:       buildStatsBlock(s.TotalRequests(), s.RequestOriginalBytes(), s.ResponseOriginalBytes(), s.ResponseWireBytes()),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func buildStatsBlock(reqs, reqBytes, respOrig, respWire int64) statsBlock {
	saved := respOrig - respWire
	var pct float64
	if respOrig > 0 {
		pct = float64(saved) / float64(respOrig) * 100
	}
	return statsBlock{
		Requests:          reqs,
		RequestBytes:      reqBytes,
		ResponseOrigBytes: respOrig,
		ResponseWireBytes: respWire,
		SavingsBytes:      saved,
		SavingsPct:        pct,
	}
}

func (d *Dashboard) handleRecentRequests(w http.ResponseWriter, r *http.Request) {
	rl := d.source.RequestLog()
	if rl == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("[]"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rl.Recent())
}

func (d *Dashboard) handleSSE(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Subscribe to request events
	var eventCh <-chan interface{}
	var unsub func()
	rl := d.source.RequestLog()
	if rl != nil {
		ch, u := rl.Subscribe()
		unsub = u
		// Wrap typed channel into interface channel
		wrappedCh := make(chan interface{}, 16)
		go func() {
			for ev := range ch {
				wrappedCh <- ev
			}
			close(wrappedCh)
		}()
		eventCh = wrappedCh
	} else {
		dummyCh := make(chan interface{})
		eventCh = dummyCh
		unsub = func() {}
	}
	defer unsub()

	// Send stats every 2 seconds
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Send initial stats immediately
	d.sendStatsEvent(w, flusher)

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			d.sendStatsEvent(w, flusher)
		case ev, ok := <-eventCh:
			if !ok {
				return
			}
			data, _ := json.Marshal(ev)
			fmt.Fprintf(w, "event: request\ndata: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func (d *Dashboard) sendStatsEvent(w http.ResponseWriter, flusher http.Flusher) {
	s := d.source.Stats()
	if s == nil {
		return
	}
	resp := statsResponse{
		UptimeSeconds: d.source.Uptime().Seconds(),
		Session:       buildStatsBlock(s.SessionRequests(), s.SessionRequestOriginalBytes(), s.SessionResponseOriginalBytes(), s.SessionResponseWireBytes()),
		AllTime:       buildStatsBlock(s.TotalRequests(), s.RequestOriginalBytes(), s.ResponseOriginalBytes(), s.ResponseWireBytes()),
	}
	data, _ := json.Marshal(resp)
	fmt.Fprintf(w, "event: stats\ndata: %s\n\n", data)
	flusher.Flush()
}

func (d *Dashboard) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}
