// Copyright 2018 The Pouch Robot Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maintainer

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

var (
	once        sync.Once
	mutex       sync.RWMutex
	maintainers []string
	// FIXME: use config file instead of hard coding
	maintainerURL = "https://github.com/alibaba/pouch/blob/master/MAINTAINERS.md"
)

func Get() []string {
	once.Do(doSync)
	mutex.RLock()
	defer mutex.RUnlock()

	return maintainers
}

func Check(user string) bool {
	maintainerList := Get()

	for _, maintainerID := range maintainerList {
		if strings.ToLower(user) == strings.ToLower(maintainerID) {
			return true
		}
	}
	return false
}

func doSync() {
	// Request the HTML page.
	res, err := http.Get(maintainerURL)
	if err != nil {
		logrus.Errorf("failed to get maintainer page: %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logrus.Errorf("getting maintainer page status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	// Load the HTML document.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Errorf("failed to load the maintainer document: %v", err)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	var newestMaintainers []string
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		// First column in table's body.
		githubID := s.Find("td").First().Text()
		if githubID != "" {
			newestMaintainers = append(newestMaintainers, githubID)
		}
	})

	maintainers = newestMaintainers

	logrus.Infof("maintainers were successfully updated to %v", maintainers)
}

func Sync() {
	for {
		go doSync()

		// Sync maintainer list everyday.
		time.Sleep(24 * time.Hour)
	}
}
