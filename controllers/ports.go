package controllers

import (
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
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
	timeout := 9000 * time.Millisecond
	delay := 1500 * time.Millisecond

	openPorts = make([]int, 0)
	failedPorts = make([]int, 0)

	done := make(chan bool)

	go func() {
		for _, port := range PORTS {
			if checkPort(domain, port, delay) {
				openPorts = append(openPorts, port)
			} else {
				failedPorts = append(failedPorts, port)
			}
		}
		done <- true
	}()

	select {
	case <-done:
		// Completed within timeout
	case <-time.After(timeout):
		// Timeout reached
		remainingPorts := make([]int, 0)
		for _, p := range PORTS {
			if !containsInt(openPorts, p) && !containsInt(failedPorts, p) {
				remainingPorts = append(remainingPorts, p)
			}
		}
		failedPorts = append(failedPorts, remainingPorts...)
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
