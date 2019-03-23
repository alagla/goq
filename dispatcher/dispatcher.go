package dispatcher

import (
	"fmt"
	. "github.com/iotaledger/iota.go/trinary"
)

// TODO size checks when join/affect. Can be with different sizes
// TODO dispose dispatcher
// TODO rename to 'supervisor' ??

func (disp *Dispatcher) incQuantCount() {
	disp.quantCountMutex.Lock()
	defer disp.quantCountMutex.Unlock()
	disp.quantCount++
}

func (disp *Dispatcher) getEnvironment_(name string) *environment {
	env, ok := disp.environments[name]
	if !ok {
		return nil
	}
	return env
}

func (disp *Dispatcher) getOrCreateEnvironment(name string) *environment {
	ret := disp.getEnvironment_(name)
	if ret != nil {
		return ret
	}
	disp.environments[name] = newEnvironment(disp, name, false)
	return disp.environments[name]
}

func (disp *Dispatcher) createEnvironment(name string, builtin bool) error {
	if disp.getEnvironment_(name) != nil {
		return fmt.Errorf("environment '%v' already exists", name)
	}
	disp.environments[name] = newEnvironment(disp, name, builtin)
	return nil
}

func (disp *Dispatcher) resetCallCounters() {
	for _, env := range disp.environments {
		for _, joinInfo := range env.joins {
			joinInfo.count = 0
		}
	}
}

func (disp *Dispatcher) quantStart(env *environment, effect Trits, waveMode bool, onQuantFinish func()) error {
	if disp.waveCoo.isWaveMode() {
		return fmt.Errorf("wave is already running")
	}
	var err error
	if effect, err = env.adjustEffect(effect); err != nil {
		return err
	}

	disp.resetCallCounters()

	disp.waveCoo.setWaveMode(waveMode)
	disp.quantWG.Add(1)

	env.effectChan <- effect

	go func() {
		env.dispatcher.quantWG.Wait()
		disp.waveCoo.setWaveMode(false)
		if onQuantFinish != nil {
			onQuantFinish()
		}
	}()
	return nil
}

// if in waveMode, continues to the next wave and stops

func (disp *Dispatcher) WaveNext() error {
	if !disp.waveCoo.isWaveMode() {
		return fmt.Errorf("not in wave mode")
	}
	disp.waveCoo.runWave()
	return nil
}

// if in waveMode, continues to the next mode and stops at the end of the quant

func (disp *Dispatcher) WaveRun() error {
	if !disp.waveCoo.isWaveMode() {
		return fmt.Errorf("not in wave mode")
	}
	disp.waveCoo.setWaveMode(false)
	disp.waveCoo.runWave()
	return nil
}

func (disp *Dispatcher) WaveValues() map[string]Trits {
	return disp.waveCoo.values()
}

func (disp *Dispatcher) IsWaveMode() bool {
	return disp.waveCoo.isWaveMode()
}
