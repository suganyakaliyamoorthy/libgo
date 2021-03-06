/*
** Copyright [2013-2017] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */
package bills

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/megamsys/libgo/api"
	"github.com/megamsys/libgo/events/alerts"
	"github.com/megamsys/libgo/pairs"
	constants "github.com/megamsys/libgo/utils"
	"strconv"
	"time"
)

const (
	EVENTSKEWS          = "/eventsskews"
	EVENTSKEWS_NEW      = "/eventsskews/content"
	EVENTEVENTSKEWSJSON = "Megam::Skews"
	HARDSKEWS           = "terminate"
	SOFTSKEWS           = "suspend"
	WARNING             = "warning"
	ACTIVE              = "active"
)

type ApiSkewsEvents struct {
	JsonClaz string        `json:"json_claz"`
	Results  []EventsSkews `json:"results"`
}
type EventsSkews struct {
	Id        string          `json:"id"`
	AccountId string          `json:"account_id"`
	CatId     string          `json:"cat_id"`
	Inputs    pairs.JsonPairs `json:"inputs"`
	Outputs   pairs.JsonPairs `json:"outputs"`
	Actions   pairs.JsonPairs `json:"actions"`
	JsonClaz  string          `json:"json_claz"`
	Status    string          `json:"status"`
	EventType string          `json:"event_type"`
}

func NewEventsSkews(email, cat_id string, mi map[string]string) ([]EventsSkews, error) {

	if email == "" {
		return nil, fmt.Errorf("account_id should not be empty")
	}

	args := api.NewArgs(mi)
	args.Email = email
	cl := api.NewClient(args, EVENTSKEWS+"/"+cat_id)
	response, err := cl.Get()
	if err != nil {
		return nil, err
	}

	ac := &ApiSkewsEvents{}
	err = json.Unmarshal(response, ac)
	if err != nil {
		return nil, err
	}
	return ac.Results, nil
}

func (s *EventsSkews) CreateEvent(o *BillOpts, ACTION string, mi map[string]string) error {
	var exp_at, gen_at time.Time
	var action, next string
	mm := make(map[string][]string, 0)
	if s.Inputs != nil {
		gen_at, _ = time.Parse(time.RFC3339, s.Inputs.Match(constants.ACTION_TRIGGERED_AT))
	} else {
		gen_at = time.Now()
	}

	softDue, err := time.ParseDuration(o.SoftGracePeriod)
	hardDue, err := time.ParseDuration(o.HardGracePeriod)
	if err != nil {
		return err
	}
	switch ACTION {
	case HARDSKEWS:
		exp_at = gen_at.Add(hardDue)
		action = HARDSKEWS
		next = "unrecoverable"
	case SOFTSKEWS:
		exp_at = gen_at.Add(hardDue)
		action = SOFTSKEWS
		next = HARDSKEWS
	case WARNING:
		mm[constants.ACTION_TRIGGERED_AT] = []string{gen_at.Format(time.RFC3339)}
		exp_at = gen_at.Add(softDue)
		action = WARNING
		next = SOFTSKEWS
	}
	mm[constants.NEXT_ACTION_DUE_AT] = []string{exp_at.Format(time.RFC3339)}
	mm[constants.ACTION] = []string{action}
	mm[constants.NEXT_ACTION] = []string{next}
	mm[constants.ASSEMBLIESID] = []string{o.AssembliesId}

	s.Inputs.NukeAndSet(mm)
	s.Status = ACTIVE
	return s.Create(mi, o)
}

func (s *EventsSkews) Create(mi map[string]string, o *BillOpts) error {
	args := api.NewArgs(mi)
	args.Email = s.AccountId
	cl := api.NewClient(args, EVENTSKEWS_NEW)
	_, err := cl.Post(s)
	if err != nil {
		return err
	}

	err = s.PushSkews(mi)
	if err != nil {
		return err
	}
	return s.skewsWarning(o)
}

func (sk *EventsSkews) skewsWarning(o *BillOpts) error {
	mm := make(map[string]string, 0)
	softDue, _ := time.ParseDuration(o.SoftGracePeriod)
	hardDue, _ := time.ParseDuration(o.HardGracePeriod)
	mm[constants.EMAIL] = sk.AccountId
	mm[constants.VERTNAME] = o.AssemblyName
	mm[constants.SOFT_ACTION] = SOFTSKEWS
	mm[constants.SOFT_GRACEPERIOD] = strconv.FormatInt(int64(softDue.Seconds()/3600/24), 10)
	mm[constants.SOFT_LIMIT] = o.SoftLimit
	mm[constants.HARD_GRACEPERIOD] = strconv.FormatInt(int64(hardDue.Seconds()/3600/24), 10)
	mm[constants.HARD_ACTION] = HARDSKEWS
	mm[constants.HARD_LIMIT] = o.HardLimit
	mm[constants.ACTION_TRIGGERED_AT] = sk.Inputs.Match(constants.ACTION_TRIGGERED_AT)
	mm[constants.NEXT_ACTION_DUE_AT] = sk.Inputs.Match(constants.NEXT_ACTION_DUE_AT)
	mm[constants.ACTION] = sk.Inputs.Match(constants.ACTION)
	mm[constants.NEXT_ACTION] = sk.Inputs.Match(constants.NEXT_ACTION)

	notifier := alerts.NewMailer(alerts.Mailer, alerts.Mailer)
	return notifier.Notify(alerts.SKEWS_WARNING, alerts.EventData{M: mm})
}

func (s *EventsSkews) PushSkews(mi map[string]string) error {
	req := api.NewRequest(s.AccountId)
	req.CatId = s.Inputs.Match(constants.ASSEMBLIESID)
	skew_action := s.Inputs.Match(constants.ACTION)
	switch skew_action {
	case HARDSKEWS:
		req.Action = constants.DESTROY
		req.Category = constants.STATE
		req.CatType = "torpedo"
	case SOFTSKEWS:
		req.Action = constants.SUSPEND
		req.Category = constants.CONTROL
		req.CatType = "torpedo"
	case WARNING:
		return nil
	}
	return req.PushRequest(mi)
}

func (s *EventsSkews) ActionEvents(o *BillOpts, currentBal string, mi map[string]string) error {
	log.Debugf("checks skews actions for ondemand")
	sk := make(map[string]*EventsSkews, 0)
	// to get skews events for that particular cat_id/ asm_id
	evts, err := NewEventsSkews(o.AccountId, o.AssemblyId, mi)
	if err != nil {
		return err
	}
	for _, v := range evts {
		if v.Status == ACTIVE {
			sk[v.Inputs.Match(constants.ACTION)] = &v
		}
	}
	ACTION := s.action(o, currentBal)

	if len(sk) > 0 {
		if sk[ACTION] != nil {
			switch true {
			case ACTION == HARDSKEWS && sk[HARDSKEWS].isExpired():
				return sk[HARDSKEWS].CreateEvent(o, HARDSKEWS, mi)
			case ACTION == SOFTSKEWS && sk[SOFTSKEWS].isExpired():
				return sk[SOFTSKEWS].CreateEvent(o, HARDSKEWS, mi)
			case ACTION == WARNING && sk[WARNING].isExpired():
				return sk[WARNING].CreateEvent(o, SOFTSKEWS, mi)
			}
			return nil
		}
	}

	return s.CreateEvent(o, ACTION, mi)
}

func (s *EventsSkews) SkewsQuotaUnpaid(o *BillOpts, mi map[string]string) error {
	log.Debugf("checks skews actions for ondemand")
	actions := make(map[string]string, 0)
	sk := make(map[string]*EventsSkews, 0)
	// to get skews events for that particular cat_id/ asm_id
	evts, err := NewEventsSkews(o.AccountId, o.AssemblyId, mi)
	if err != nil {
		return err
	}
	for _, v := range evts {
		if v.Status == ACTIVE {
			sk[v.Inputs.Match(constants.ACTION)] = &v
			actions[v.Inputs.Match(constants.ACTION)] = ACTIVE
		}
	}
	if len(sk) > 0 {
		switch true {
		case actions[HARDSKEWS] == ACTIVE && sk[HARDSKEWS].isExpired():
			return sk[HARDSKEWS].CreateEvent(o, HARDSKEWS, mi)
		case actions[SOFTSKEWS] == ACTIVE && sk[SOFTSKEWS].isExpired():
			return sk[SOFTSKEWS].CreateEvent(o, HARDSKEWS, mi)
		case actions[WARNING] == ACTIVE && sk[WARNING].isExpired():
			return sk[WARNING].CreateEvent(o, SOFTSKEWS, mi)
		}
		return nil
	}

	return s.CreateEvent(o, WARNING, mi)
}

func (s *EventsSkews) action(o *BillOpts, currentBal string) string {
	cb, _ := strconv.ParseFloat(currentBal, 64)
	slimit, _ := strconv.ParseFloat(o.SoftLimit, 64)
	hlimit, _ := strconv.ParseFloat(o.HardLimit, 64)
	if cb <= hlimit {
		return HARDSKEWS
	} else if cb <= slimit {
		return SOFTSKEWS
	}
	return WARNING
}

func (s *EventsSkews) isExpired() bool {
	t1, _ := time.Parse(time.RFC3339, s.Inputs.Match(constants.ACTION_TRIGGERED_AT))
	t2, _ := time.Parse(time.RFC3339, s.Inputs.Match(constants.NEXT_ACTION_DUE_AT))
	duration := t2.Sub(t1)
	return t1.Add(duration).Sub(time.Now()) < time.Minute
}
