package devops

import (
	"time"

	"github.com/timescale/tsbs/pkg/data"
	"github.com/timescale/tsbs/pkg/data/usecases/common"
)

// DevopsSimulator generates data similar to telemetry, with metrics from a variety of device systems.
// It fulfills the Simulator interface.
type DevopsSimulator struct {
	*commonDevopsSimulator
	simulatedMeasurementIndex int
	firstMeasurement          bool
}

// Next advances a Point to the next state in the generator.
// func (d *DevopsSimulator) Next(p *data.Point) bool {
// 	// switch to the next metric if needed
// 	if d.hostIndex == uint64(len(d.hosts)) {
// 		d.hostIndex = 0
// 		d.simulatedMeasurementIndex++
// 	}

// 	if d.simulatedMeasurementIndex == len(d.hosts[0].SimulatedMeasurements) {
// 		d.simulatedMeasurementIndex = 0

// 		for i := 0; i < len(d.hosts); i++ {
// 			d.hosts[i].TickAll(d.interval)
// 		}

// 		d.adjustNumHostsForEpoch()
// 	}

// 	return d.populatePoint(p, d.simulatedMeasurementIndex)
// }

// Next advances a Point to the next state in the generator.
func (d *DevopsSimulator) Next(p *data.Point) bool {
	// when all the measurements for a certain host have been generated the
	// simulator needs to adjust the clock for this host and move on to the next one
	if d.simulatedMeasurementIndex == len(d.hosts[0].SimulatedMeasurements)-1 {
		d.hosts[d.hostIndex].TickAll(d.interval)
		d.hostIndex++
		d.simulatedMeasurementIndex = 0
	} else {
		// generate the new measurement for the current host
		if d.firstMeasurement {
			d.firstMeasurement = false
		} else {
			d.simulatedMeasurementIndex++
		}
	}

	// if we have iterated through all the hosts, move on to the new epoch
	if d.hostIndex == uint64(len(d.hosts)) {
		d.hostIndex = 0
		d.adjustNumHostsForEpoch()
	}

	return d.populatePoint(p, d.simulatedMeasurementIndex)
}

func (s *DevopsSimulator) TagKeys() []string {
	tagKeysAsStr := make([]string, len(MachineTagKeys))
	for i, t := range MachineTagKeys {
		tagKeysAsStr[i] = string(t)
	}
	return tagKeysAsStr
}

func (s *DevopsSimulator) TagTypes() []string {
	types := make([]string, len(MachineTagKeys))
	for i := 0; i < len(MachineTagKeys); i++ {
		types[i] = machineTagType.String()
	}
	return types
}

func (d *DevopsSimulator) Headers() *common.GeneratedDataHeaders {
	return &common.GeneratedDataHeaders{
		TagTypes:  d.TagTypes(),
		TagKeys:   d.TagKeys(),
		FieldKeys: d.Fields(),
	}
}

// DevopsSimulatorConfig is used to create a DevopsSimulator.
type DevopsSimulatorConfig commonDevopsSimulatorConfig

// NewSimulator produces a Simulator that conforms to the given SimulatorConfig over the specified interval
func (d *DevopsSimulatorConfig) NewSimulator(interval time.Duration, limit uint64) common.Simulator {
	hostInfos := make([]Host, d.HostCount)
	for i := 0; i < len(hostInfos); i++ {
		hostInfos[i] = d.HostConstructor(NewHostCtx(i, d.Start))
	}

	epochs := calculateEpochs(commonDevopsSimulatorConfig(*d), interval)
	maxPoints := epochs * d.HostCount * uint64(len(hostInfos[0].SimulatedMeasurements))
	if limit > 0 && limit < maxPoints {
		// Set specified points number limit
		maxPoints = limit
	}
	dg := &DevopsSimulator{
		commonDevopsSimulator: &commonDevopsSimulator{
			madePoints: 0,
			maxPoints:  maxPoints,

			hostIndex: 0,
			hosts:     hostInfos,

			epoch:          0,
			epochs:         epochs,
			epochHosts:     d.InitHostCount,
			initHosts:      d.InitHostCount,
			timestampStart: d.Start,
			timestampEnd:   d.End,
			interval:       interval,
		},
		simulatedMeasurementIndex: 0,
		firstMeasurement:          true,
	}

	return dg
}
