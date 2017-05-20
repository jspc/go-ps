// +build linux

package ps

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
)

// Refresh reloads all the data associated with this process.
func (p *UnixProcess) Refresh() error {
	procPath := filepath.Join("/proc", strconv.Itoa(p.pid))

	argsPath := filepath.Join(procPath, "cmdline")
	argsBytes, err := ioutil.ReadFile(argsPath)
	if err != nil {
		return err
	}
	p.args = strings.Split(string(argsBytes), " ")

	statPath := filepath.Join(procPath, "stat")
	dataBytes, err := ioutil.ReadFile(statPath)
	if err != nil {
		return err
	}

	// First, parse out the image name
	data := string(dataBytes)
	binStart := strings.IndexRune(data, '(') + 1
	binEnd := strings.IndexRune(data[binStart:], ')')
	p.binary = data[binStart : binStart+binEnd]

	// Move past the image name and start parsing the rest
	data = data[binStart+binEnd+2:]
	_, err = fmt.Sscanf(data,
		"%c %d %d %d",
		&p.state,
		&p.ppid,
		&p.pgrp,
		&p.sid)

	return err
}
