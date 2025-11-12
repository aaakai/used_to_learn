package voice

import (
	"context"
	"fmt"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

type CallState string

const (
	StateIdle     CallState = "idle"
	StateCalling  CallState = "calling"
	StateTalking  CallState = "talking"
	StateRejected CallState = "rejected"
)

type VoiceManager struct {
	ctx      context.Context
	cancel   context.CancelFunc
	state    CallState
	callCh   chan struct{}
	answerCh chan bool
	streamer beep.StreamSeekCloser
	format   beep.Format
}

func NewVoiceManager() *VoiceManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &VoiceManager{
		ctx:      ctx,
		cancel:   cancel,
		state:    StateIdle,
		callCh:   make(chan struct{}),
		answerCh: make(chan bool),
	}
}

func (vm *VoiceManager) Init() error {
	sampleRate := beep.SampleRate(44100)
	speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	return nil
}

func (vm *VoiceManager) Call(number string, timeout time.Duration) error {
	if vm.state != StateIdle {
		return fmt.Errorf("cannot call in %s state", vm.state)
	}

	vm.state = StateCalling
	go func() {
		select {
		case <-time.After(timeout):
			vm.state = StateIdle
		case answered := <-vm.answerCh:
			if answered {
				vm.state = StateTalking
			} else {
				vm.state = StateRejected
				time.Sleep(2 * time.Second)
				vm.state = StateIdle
			}
		}
	}()

	// 简单蜂鸣声替代SineWave
	buzzer := beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for i := range samples {
			samples[i][0] = 0.5
			samples[i][1] = 0.5
		}
		return len(samples), true
	})
	sampleRate := beep.SampleRate(44100)
	speaker.Play(beep.Take(sampleRate.N(time.Second/2), buzzer))
	return nil
}

func (vm *VoiceManager) Answer() error {
	if vm.state != StateCalling {
		return fmt.Errorf("cannot answer in %s state", vm.state)
	}
	vm.answerCh <- true
	return nil
}

func (vm *VoiceManager) Reject() error {
	if vm.state != StateCalling {
		return fmt.Errorf("cannot reject in %s state", vm.state)
	}
	vm.answerCh <- false
	return nil
}

func (vm *VoiceManager) GetState() CallState {
	return vm.state
}

func (vm *VoiceManager) Close() {
	vm.cancel()
	if vm.streamer != nil {
		vm.streamer.Close()
	}
	speaker.Clear()
}
