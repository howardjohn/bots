// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package userdata

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ghodss/yaml"

	"istio.io/bots/policybot/pkg/storage"
)

type Time struct {
	time.Time
}

// A company a user is associated with
type Affiliation struct {
	Organization string `json:"organization"`
	Start        Time   `json:"start"`
	End          Time   `json:"end"`
}

// Additional info about an Istio contributor
type UserInfo struct {
	GitHubLogin    string        `json:"github_login"`
	Affiliations   []Affiliation `json:"affiliations"`
	EmailAddresses []string      `json:"email_addresses"`
}

type UserData struct {
	Users []UserInfo `json:"users"`
}

func Load(file string) (UserData, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return UserData{}, fmt.Errorf("unable to read user data file %s: %v", file, err)
	}

	var ud UserData
	if err = yaml.Unmarshal(b, &ud); err != nil {
		return UserData{}, fmt.Errorf("unable to parse user data file %s: %v", file, err)
	}

	return ud, nil
}

func (ud UserData) Store(store storage.Store) error {
	var a []*storage.UserAffiliation

	for _, user := range ud.Users {
		for _, affiliation := range user.Affiliations {
			a = append(a, &storage.UserAffiliation{
				UserLogin:    user.GitHubLogin,
				Organization: affiliation.Organization,
				StartTime:    affiliation.Start.Time,
				EndTime:      affiliation.End.Time,
			})
		}
	}

	return store.WriteAllUserAffiliations(context.Background(), a)
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format("2006-01-02") + `"`), nil
}

func (t *Time) UnmarshalJSON(buf []byte) error {
	tt, err := time.Parse("2006-01-02", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}