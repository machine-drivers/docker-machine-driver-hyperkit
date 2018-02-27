/*
Copyright 2016 The Kubernetes Authors All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hyperkit

import (
	"time"
	"errors"
	"strings"
	"os/exec"
	"os"
	"github.com/docker/machine/libmachine/log"
	"bufio"
	"fmt"
)

type RetriableError struct {
	Err error
}

func (r RetriableError) Error() string {
	return "Temporary Error: " + r.Err.Error()
}

type MultiError struct {
	Errors []error
}

func (m *MultiError) Collect(err error) {
	if err != nil {
		m.Errors = append(m.Errors, err)
	}
}

func (m MultiError) ToError() error {
	if len(m.Errors) == 0 {
		return nil
	}

	errStrings := []string{}
	for _, err := range m.Errors {
		errStrings = append(errStrings, err.Error())
	}
	return errors.New(strings.Join(errStrings, "\n"))
}

func RetryAfter(attempts int, callback func() error, d time.Duration) (err error) {
	m := MultiError{}
	for i := 0; i < attempts; i++ {
		err = callback()
		if err == nil {
			return nil
		}
		m.Collect(err)
		if _, ok := err.(*RetriableError); !ok {
			return m.ToError()
		}
		time.Sleep(d)
	}
	return m.ToError()
}

func hdiutil(args ...string) error {
	cmd := exec.Command("hdiutil", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debugf("executing: %v %v", cmd, strings.Join(args, " "))

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func readLine(path string) (string, error) {
	inFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer inFile.Close()

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		if kernelOptionRegexp.Match(scanner.Bytes()) {
			m := kernelOptionRegexp.FindSubmatch(scanner.Bytes())
			return string(m[1]), nil
		}
	}
	return "", fmt.Errorf("couldn't find kernel option from %s image", path)
}