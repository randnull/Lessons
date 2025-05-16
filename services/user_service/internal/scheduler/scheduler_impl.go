package scheduler

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/randnull/Lessons/internal/config"
	"github.com/randnull/Lessons/internal/models"
	"github.com/randnull/Lessons/internal/repository"
	lg "github.com/randnull/Lessons/pkg/logger"
	"github.com/randnull/Lessons/pkg/rabbitmq"
	"time"
)

type Scheduler struct {
	cfg            *config.SchedulerConfig
	userRepository repository.UserRepository
	ProducerBroker rabbitmq.RabbitMQInterface
	scheduler      *gocron.Scheduler
}

func NewScheduler(cfg *config.SchedulerConfig, userRepo repository.UserRepository, producerBroker rabbitmq.RabbitMQInterface) *Scheduler {
	scheduler := gocron.NewScheduler(time.UTC)
	return &Scheduler{
		cfg:            cfg,
		userRepository: userRepo,
		ProducerBroker: producerBroker,
		scheduler:      scheduler,
	}
}

func (s *Scheduler) RunResponseChecker(ctx context.Context) {
	lg.Info(fmt.Sprintf("[Scheduler] start initing at %v", time.Now()))

	_, err := s.scheduler.Every(1).Week().Sunday().At("20:30").Do(func() {
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
				lg.Error("[Scheduler] cannot push to broker, error: " + err.Error())
			}
		}
	})

	if err != nil {
		lg.Error(fmt.Sprintf("[Scheduler] error with scheduler: %v", err.Error()))
		return
	}

	s.scheduler.StartAsync()
	lg.Info("[Scheduler] scheduler started. Waiting Sunday 20:30 UTC...")

	<-ctx.Done()
	lg.Info("[Scheduler] stopped")
	s.scheduler.Stop()
}
