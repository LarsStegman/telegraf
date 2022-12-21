package fit

import (
	"bytes"
	"github.com/influxdata/telegraf/metric"
	"github.com/influxdata/telegraf/plugins/parsers"
	"math"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/tormoder/fit"
)

type Parser struct {
	MetricName  string `toml:"metric_name"`
	DefaultTags map[string]string
}

func (p *Parser) Parse(buf []byte) ([]telegraf.Metric, error) {
	fitFile, err := fit.Decode(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	activity, err := fitFile.Activity()
	if err != nil {
		return nil, err
	}

	results := make([]telegraf.Metric, 0)
	tags := p.DefaultTags
	tags["activity_time"] = activity.Activity.Timestamp.Format(time.RFC3339)
	for _, record := range activity.Records {
		m := metric.New("record", tags, p.ActivityRecordToMetric(record), record.Timestamp)
		results = append(results, m)
	}

	for _, event := range activity.Events {
		m := metric.New("event", tags, p.ActivityEventToMetric(event), event.Timestamp)
		results = append(results, m)
	}

	for _, lap := range activity.Laps {
		m := metric.New("lap", tags, p.ActivityLapToMetric(lap), lap.Timestamp)
		results = append(results, m)
	}

	return results, nil
}

func CreateFieldIfNotNan(field string, val interface{}, dict map[string]interface{}) {
	if val != math.NaN() {
		dict[field] = val
	}
}

func (p *Parser) ActivityRecordToMetric(r *fit.RecordMsg) map[string]interface{} {
	invalidMsg := fit.NewRecordMsg()
	fields := make(map[string]interface{})
	CreateFieldIfNotNan("latitude__deg", r.PositionLat.Degrees(), fields)
	CreateFieldIfNotNan("longitude__deg", r.PositionLong.Degrees(), fields)
	CreateFieldIfNotNan("altitude__m", r.GetEnhancedAltitudeScaled(), fields)
	if r.HeartRate != invalidMsg.HeartRate {
		fields["heart_rate__bpm"] = r.HeartRate
	}

	if r.Cadence != invalidMsg.Cadence {
		fields["cadence__rpm"] = r.Cadence
	}
	CreateFieldIfNotNan("distance__m", r.GetDistanceScaled(), fields)
	CreateFieldIfNotNan("speed__m/s", r.GetEnhancedSpeedScaled(), fields)
	if r.Power != invalidMsg.Power {
		fields["power_W"] = r.Power
	}

	CreateFieldIfNotNan("grade__%", r.GetGradeScaled(), fields)
	if r.Resistance != invalidMsg.Resistance {
		fields["resistance_level"] = r.Resistance
	}

	CreateFieldIfNotNan("time_from_course__s", r.GetTimeFromCourseScaled(), fields)
	CreateFieldIfNotNan("cycle_length__m", r.GetCycleLengthScaled(), fields)
	if r.Temperature != 0x7F {
		fields["temperature__degC"] = r.Temperature
	}

	if r.Cycles != invalidMsg.Cycles {
		fields["cycles"] = r.Cycles
	}

	if r.TotalCycles != invalidMsg.TotalCycles {
		fields["cycles_total"] = r.TotalCycles
	}

	if r.CompressedAccumulatedPower != invalidMsg.CompressedAccumulatedPower {
		fields["power_accumulated_compressed__W"] = r.CompressedAccumulatedPower
	}

	if r.AccumulatedPower != invalidMsg.AccumulatedPower {
		fields["power_accumulated__W"] = r.AccumulatedPower
	}

	CreateFieldIfNotNan("pedal_left_contribution__%", LeftContribution(r.LeftRightBalance), fields)
	if r.GpsAccuracy != invalidMsg.GpsAccuracy {
		fields["gps_accuracy__m"] = r.GpsAccuracy
	}

	CreateFieldIfNotNan("speed_vertical__m/s", r.GetVerticalSpeedScaled(), fields)
	if r.Calories != invalidMsg.Calories {
		fields["calories__kcal"] = r.Calories
	}

	CreateFieldIfNotNan("oscillation_vertical__mm", r.GetVerticalOscillationScaled(), fields)
	CreateFieldIfNotNan("stance_time__%", r.GetStanceTimePercentScaled(), fields)
	CreateFieldIfNotNan("stance_time__ms", r.GetStanceTimeScaled(), fields)
	CreateFieldIfNotNan("torque_effectiveness_left__%", r.GetLeftTorqueEffectivenessScaled(), fields)
	CreateFieldIfNotNan("torque_effectiveness_right__%", r.GetRightTorqueEffectivenessScaled(), fields)
	CreateFieldIfNotNan("pedal_smoothness_left__%", r.GetLeftPedalSmoothnessScaled(), fields)
	CreateFieldIfNotNan("pedal_smoothness_right__%", r.GetRightPedalSmoothnessScaled(), fields)
	CreateFieldIfNotNan("pedal_smoothness__%", r.GetCombinedPedalSmoothnessScaled(), fields)
	CreateFieldIfNotNan("time__s", r.GetTime128Scaled(), fields)
	if r.StrokeType != invalidMsg.StrokeType {
		fields["stroke"] = r.StrokeType.String()
	}

	// TODO: Zone   Not sure what Zone means. What kind of zone? Heart rate? Speed? Cadence?
	CreateFieldIfNotNan("ball_speed__m/s", r.GetBallSpeedScaled(), fields)
	// TODO: Cadence256?   Deprecated I think, not sure what it means
	// TODO: Fractional cadence?   Deprecated I think, fractional part of Cadence 256, maybe?
	CreateFieldIfNotNan("hemoglobin_concentration_total__g/dL", r.GetTotalHemoglobinConcScaled(), fields)
	CreateFieldIfNotNan("hemoglobin_concentration_total_min__g/dL", r.GetTotalHemoglobinConcMinScaled(), fields)
	CreateFieldIfNotNan("hemoglobin_concentration_total_max__g/dL", r.GetTotalHemoglobinConcMaxScaled(), fields)
	CreateFieldIfNotNan("hemoglobin_saturated_total__%", r.GetSaturatedHemoglobinPercentScaled(), fields)
	CreateFieldIfNotNan("hemoglobin_saturated_total_min__%", r.GetSaturatedHemoglobinPercentMinScaled(), fields)
	CreateFieldIfNotNan("hemoglobin_saturated_total_max__%", r.GetSaturatedHemoglobinPercentMaxScaled(), fields)

	return fields
}

func (p *Parser) ActivityEventToMetric(e *fit.EventMsg) map[string]interface{} {
	fields := make(map[string]interface{})
	fields["event"] = e.Event.String()
	fields["event_type"] = e.EventType.String()
	return fields
}

func (p *Parser) ActivityLapToMetric(lap *fit.LapMsg) map[string]interface{} {
	fields := make(map[string]interface{})
	CreateFieldIfNotNan("elapsed_time__s", lap.GetTotalElapsedTimeScaled(), fields)
	CreateFieldIfNotNan("elapsed_time__s", lap.GetTotalMovingTimeScaled(), fields)
	CreateFieldIfNotNan("moving_time__s", lap.GetTotalMovingTimeScaled(), fields)
	CreateFieldIfNotNan("distance__m", lap.GetTotalDistanceScaled(), fields)
	CreateFieldIfNotNan("calories__kcal", lap.TotalCalories, fields)
	CreateFieldIfNotNan("heart_rate_avg__bpm", lap.AvgHeartRate, fields)
	CreateFieldIfNotNan("heart_rate_max__bpm", lap.MaxHeartRate, fields)
	CreateFieldIfNotNan("cadence_avg__rpm", lap.AvgCadence, fields)
	CreateFieldIfNotNan("speed_avg__m/s", lap.GetEnhancedAvgSpeedScaled(), fields)
	CreateFieldIfNotNan("speed_max__m/s", lap.GetEnhancedMaxSpeedScaled(), fields)
	CreateFieldIfNotNan("pedal_left_contribution__%", LeftContribution100(lap.LeftRightBalance), fields)

	return fields
}

// ParseLine takes a single string metric
// ie, "cpu.usage.idle 90"
// and parses it into a telegraf metric.
//
// Must be thread-safe.
func (p *Parser) ParseLine(_ string) (telegraf.Metric, error) {
	return nil, nil
}

// SetDefaultTags tells the parser to add all the given tags
// to each parsed metric.
// NOTE: do _not_ modify the map after you've passed it here!!
func (p *Parser) SetDefaultTags(tags map[string]string) {
	p.DefaultTags = tags
}

func init() {
	parsers.Add("fit",
		func(defaultMetricName string) telegraf.Parser {
			return &Parser{MetricName: defaultMetricName}
		})
}
