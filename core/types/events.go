// Copyright 2015, 2016 Eris Industries (UK) Ltd.
// This file is part of Eris-RT

// Eris-RT is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Eris-RT is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Eris-RT.  If not, see <http://www.gnu.org/licenses/>.

package types

import (
	evts "github.com/tendermint/go-events"
)

// TODO improve
// TODO: [ben] yes please ^^^
// [ben] To improve this we will switch out go-events with eris-db/event so
// that there is no need anymore for this poor wrapper.

// The events struct has methods for working with events.
type events struct {
	eventSwitch *evts.EventSwitch
}

func newEvents(eventSwitch *evts.EventSwitch) *events {
	return &events{eventSwitch}
}

// Subscribe to an event.
func (this *events) Subscribe(subId, event string, callback func(evts.EventData)) (bool, error) {
	this.eventSwitch.AddListenerForEvent(subId, event, callback)
	return true, nil
}

// Un-subscribe from an event.
func (this *events) Unsubscribe(subId string) (bool, error) {
	this.eventSwitch.RemoveListener(subId)
	return true, nil
}
