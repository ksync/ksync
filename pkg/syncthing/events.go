package syncthing

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/events"
)

// TODO: figure out how to use events.EventType instead of string here.
func (s *Server) Events(types ...string) (
	<-chan *events.Event, error) {
	out := make(chan *events.Event)

	// TODO: this will replay all the events for new shares (potentially a ton
	// for long running watch processes). Should only get the latest.
	since := 0

	params := map[string]string{
		"since": strconv.Itoa(since),
	}

	// TODO: there appears to be a bug in syncthing that does not accept
	// types separated by escaped commas. resty does the right thing, and this
	// ends up breaking the filtering.
	// if len(types) > 0 {
	// 	params["events"] = strings.Join(types, ",")
	// }

	go func() {
		for {
			select {
			case <-s.stop:
				log.WithFields(s.Fields()).Debug("halting events polling")
				return
			default:
				params["since"] = strconv.Itoa(since)
				resp, err := s.client.NewRequest().
					SetQueryParams(params).
					SetResult([]events.Event{}).
					Get("events")

				if err != nil {
					log.Warn(err)
					continue
				}

				for _, event := range *resp.Result().(*[]events.Event) {
					since = event.SubscriptionID
					out <- &event
				}
			}
		}
	}()

	return out, nil
}
