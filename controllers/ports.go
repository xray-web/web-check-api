package controllers

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type PortsController struct{}

var PORTS = []int{
	20, 21, 22, 23, 25, 53, 80, 67, 68, 69,
	110, 119, 123, 143, 156, 161, 162, 179, 194,
	389, 443, 587, 993, 995,
	3000, 3306, 3389, 5060, 5900, 8000, 8080, 8888,
}

func (ctrl *PortsController) GetPortsHandler(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	domain := strings.TrimPrefix(url, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	openPorts, failedPorts := checkPorts(domain)

	c.JSON(http.StatusOK, gin.H{
		"openPorts":   openPorts,
		"failedPorts": failedPorts,
	})
}

func checkPorts(domain string) (openPorts []int, failedPorts []int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	openPorts = make([]int, 0)
	failedPorts = make([]int, 0)

	timeout := 1500 * time.Millisecond
	overallTimeout := 9000 * time.Millisecond

	done := make(chan struct{})

	go func() {
		for _, port := range PORTS {
			wg.Add(1)
			go func(port int) {
				defer wg.Done()
				if checkPort(domain, port, timeout) {
					mu.Lock()
					openPorts = append(openPorts, port)
					mu.Unlock()
				} else {
					mu.Lock()
					failedPorts = append(failedPorts, port)
					mu.Unlock()
				}
			}(port)
		}
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(overallTimeout):
		mu.Lock()
		for _, port := range PORTS {
			if !containsInt(openPorts, port) && !containsInt(failedPorts, port) {
				failedPorts = append(failedPorts, port)
			}
		}
		mu.Unlock()
	}

	sort.Ints(openPorts)
	sort.Ints(failedPorts)
	return
}

func checkPort(domain string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", domain, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func containsInt(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func HandleGetPorts() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")
		if url == "" {
			JSONError(w, ErrMissingURLParameter, http.StatusBadRequest)
			return
		}
	})
}
