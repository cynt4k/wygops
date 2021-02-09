package sync

import (
	"fmt"
	"time"
)

type syncFunction func(s *Service)

func (s *Service) runTimer(cb syncFunction) {
	ticker := time.NewTicker(s.duration)

	defer ticker.Stop()

	go cb(s)

	for {
		select {
		case <-s.syncInterrupt:
			s.logger.Info(fmt.Sprintf("sync job for %s stopped", s.sourceType))
			ticker.Stop()
			return
		case <-ticker.C:
			s.logger.Debug(fmt.Sprintf("running scheduled sync for %s", s.sourceType))
			go cb(s)
		}
	}
}
