package task

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dellinger2023/net-flux/gen"
	"github.com/dellinger2023/net-flux/pkg/dataflow/core"
)

const (
	defaultProbeCount        = 5
	defaultProbeTimeout      = 5 * time.Second
	probeInterval            = 100 * time.Millisecond
	probeUserAgent           = "net-flux-prober/1.0"
	unreachableRTTMS         = int32(5000) // > 1000
	unreachableJitterMS      = int32(5001) // > 5000
	unreachablePacketLoss    = 1.0
)

type netProber struct {
	destAddr   string
	probeCount int
	timeout    time.Duration
	client     *http.Client

	nodeInfo *core.NodeInfo
	mu       sync.Mutex
	result   *gen.NetworkMetric
}

func newProber(destAddr string, nodeInfo *core.NodeInfo) *netProber {
	p := &netProber{
		destAddr:   destAddr,
		probeCount: defaultProbeCount,
		timeout:    defaultProbeTimeout,
		nodeInfo:   nodeInfo,
	}
	p.client = p.buildClient()
	return p
}

func (p *netProber) Result() *gen.NetworkMetric {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.result
}

func (p *netProber) Probe() error {
	targetURLs, err := probeTargetURLs(p.destAddr)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	rtts := make([]time.Duration, 0, p.probeCount)
	failures := 0
	lastStatus := 0
	var lastErr error

	for i := 0; i < p.probeCount; i++ {
		if ctx.Err() != nil {
			failures += p.probeCount - i
			lastErr = ctx.Err()
			break
		}
		if i > 0 {
			select {
			case <-ctx.Done():
				failures += p.probeCount - i
				lastErr = ctx.Err()
				i = p.probeCount
			case <-time.After(probeInterval):
			}
			if i >= p.probeCount {
				break
			}
		}

		start := time.Now()
		status, probeErr := p.doProbe(ctx, targetURLs)
		if probeErr != nil {
			failures++
			lastErr = probeErr
			continue
		}

		rtts = append(rtts, time.Since(start))
		lastStatus = status
	}

	targetURL := targetURLs[0]
	destHost, err := resolveDestinationIP(targetURL)
	if err != nil {
		destHost = destinationHost(targetURL)
	}

	result := p.buildResult(targetURL, destHost, rtts, failures, lastStatus, lastErr)

	p.mu.Lock()
	p.result = result
	p.mu.Unlock()
	return nil
}

func (p *netProber) buildResult(
	targetURL, destHost string,
	rtts []time.Duration,
	failures, lastStatus int,
	lastErr error,
) *gen.NetworkMetric {
	sourceIP := ""
	nodeID := ""
	machineID := ""
	if p.nodeInfo != nil {
		sourceIP = p.nodeInfo.IP
		nodeID = strconv.Itoa(p.nodeInfo.Node)
		machineID = p.nodeInfo.ID
	}

	packetLoss := float64(failures) / float64(p.probeCount)
	if len(rtts) == 0 {
		extra := map[string]string{
			"url":         targetURL,
			"status_code": "0",
			"success":     "0",
			"total":       strconv.Itoa(p.probeCount),
			"reachable":   "false",
		}
		if lastErr != nil {
			extra["error"] = lastErr.Error()
		}
		return &gen.NetworkMetric{
			MachineId:     machineID,
			NodeId:        nodeID,
			SourceIp:      sourceIP,
			DestinationIp: destHost,
			Rtt:           unreachableRTTMS,
			Jitter:        unreachableJitterMS,
			PacketLoss:    unreachablePacketLoss,
			StatusCode:    0,
			Timestamp:     time.Now().Unix(),
			Extra:         extra,
		}
	}

	rttMS := int32(averageDuration(rtts).Milliseconds())
	if rttMS == 0 {
		rttMS = 1
	}

	return &gen.NetworkMetric{
		MachineId:     machineID,
		NodeId:        nodeID,
		SourceIp:      sourceIP,
		DestinationIp: destHost,
		Rtt:           rttMS,
		Jitter:        int32(jitterDuration(rtts).Milliseconds()),
		PacketLoss:    packetLoss,
		StatusCode:    int32(lastStatus),
		Timestamp:     time.Now().Unix(),
		Extra: map[string]string{
			"url":         targetURL,
			"status_code": strconv.Itoa(lastStatus),
			"success":     strconv.Itoa(len(rtts)),
			"total":       strconv.Itoa(p.probeCount),
			"reachable":   "true",
		},
	}
}

func (p *netProber) buildClient() *http.Client {
	dialer := &net.Dialer{Timeout: p.timeout}
	return &http.Client{
		Timeout: p.timeout,
		Transport: &http.Transport{
			DisableKeepAlives:     true,
			DialContext:           dialer.DialContext,
			TLSHandshakeTimeout:   p.timeout,
			ResponseHeaderTimeout: p.timeout,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return errors.New("too many redirects")
			}
			return nil
		},
	}
}

func (p *netProber) doProbe(ctx context.Context, targetURLs []string) (int, error) {
	perURLTimeout := p.timeout / time.Duration(len(targetURLs))
	if perURLTimeout <= 0 {
		perURLTimeout = p.timeout
	}

	var lastErr error
	for _, targetURL := range targetURLs {
		if ctx.Err() != nil {
			return 0, ctx.Err()
		}

		urlCtx, cancel := context.WithTimeout(ctx, perURLTimeout)
		status, err := p.doRequest(urlCtx, http.MethodGet, targetURL)
		cancel()
		if err == nil {
			return status, nil
		}
		lastErr = err
	}
	return 0, lastErr
}

func (p *netProber) doRequest(ctx context.Context, method, rawURL string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", probeUserAgent)

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= http.StatusBadRequest {
		return resp.StatusCode, fmt.Errorf("http probe got status %d", resp.StatusCode)
	}
	return resp.StatusCode, nil
}

func probeTargetURLs(addr string) ([]string, error) {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, errors.New("destination address is empty")
	}

	if strings.Contains(addr, "://") {
		u, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		if u.Host == "" {
			return nil, fmt.Errorf("invalid http address: %s", addr)
		}
		return []string{u.String()}, nil
	}

	host := addr
	if strings.Contains(host, "/") {
		host = strings.Split(host, "/")[0]
	}

	path := "/"
	if idx := strings.Index(addr, "/"); idx >= 0 {
		path = addr[idx:]
		if path == "" {
			path = "/"
		}
	}

	return []string{
		"https://" + host + path,
		"http://" + host + path,
	}, nil
}

func destinationHost(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return u.Hostname()
	}
	return host
}

func resolveDestinationIP(rawURL string) (string, error) {
	host := destinationHost(rawURL)
	if ip := net.ParseIP(host); ip != nil {
		return ip.String(), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultProbeTimeout)
	defer cancel()

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no ip found for host %s", host)
	}
	return ips[0].IP.String(), nil
}

func averageDuration(values []time.Duration) time.Duration {
	if len(values) == 0 {
		return 0
	}

	var total time.Duration
	for _, v := range values {
		total += v
	}
	return total / time.Duration(len(values))
}

func jitterDuration(values []time.Duration) time.Duration {
	if len(values) < 2 {
		return 0
	}

	var total time.Duration
	for i := 1; i < len(values); i++ {
		diff := values[i] - values[i-1]
		if diff < 0 {
			diff = -diff
		}
		total += diff
	}
	return total / time.Duration(len(values)-1)
}
