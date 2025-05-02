package scheduler

import (
	"context"
	"fmt"
	"github.com/randnull/Lessons/internal/config"
	lg "github.com/randnull/Lessons/internal/logger"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/rabbitmq"
	"github.com/randnull/Lessons/internal/repository"
	"time"
)

type Scheduler struct {
	cfg            *config.SchedulerConfig
	userRepository repository.UserRepository
	ProducerBroker rabbitmq.RabbitMQInterface
}

func NewScheduler(cfg *config.SchedulerConfig, userRepo repository.UserRepository, producerBroker rabbitmq.RabbitMQInterface) *Scheduler {
	return &Scheduler{
		cfg:            cfg,
		userRepository: userRepo,
		ProducerBroker: producerBroker,
	}
}

func (s *Scheduler) RunResponseChecker(ctx context.Context) {
	lg.Info("[Scheduler] start working")

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	func() {
		for {
			select {
			case <-ctx.Done():
				lg.Info("[Scheduler] stopped")
				ticker.Stop()
				return
			case <-ticker.C:
				lg.Info("[Scheduler] starting new scheduled job")
				tutors, err := s.userRepository.GetAllTutorsResponseCondition(5)
				if err != nil {
					lg.Error("[Scheduler] cannot get tutors, error: " + err.Error())
				}
				lg.Info(fmt.Sprintf("[Scheduler] get %v tutors to update response count", len(tutors)))

				for _, tutor := range tutors {
					totalResponses, err := s.userRepository.AddResponses(tutor.TelegramID, int(5-tutor.ResponseCount))
					if err != nil {
						lg.Error("[Scheduler] cannot get tutors, error: " + err.Error())
					}
					NotifyModel := models.AddResponsesToTutor{
						TutorTelegramID: tutor.TelegramID,
						ResponseCount:   totalResponses,
					}
					err = s.ProducerBroker.Publish("add_responses", &NotifyModel)
					if err != nil {
						lg.Error("[Scheduler] cannot pushed to broker, error: " + err.Error())
					}
				}
			}
		}
	}()
}
