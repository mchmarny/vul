package metric

import (
	"context"

	"fmt"
	"time"

	"github.com/mchmarny/vul/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredres "google.golang.org/genproto/googleapis/api/monitoredres"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	metricBaseType = "custom.googleapis.com"
	resourceType   = "global"

	serviceEnvVar  = "K_SERVICE"
	revisionEnvVar = "K_REVISION"
)

// Recorder is a wrapper around the metric service.
type Recorder interface {
	One(ctx context.Context, name string)
	OneWithLabels(ctx context.Context, name string, labels map[string]string)
	WithLabels(ctx context.Context, name string, value interface{}, labels map[string]string)
}

// New creates a new metric service instance with the specified configuration.
func New(projectID, serviceName, serviceVersion string, send bool) (Recorder, error) {
	if projectID == "" || serviceName == "" || serviceVersion == "" {
		log.Error().
			Str("projectID", projectID).
			Str("serviceName", serviceName).
			Str("serviceVersion", serviceVersion).
			Msg("required parameters not provided")
		return nil, errors.New("required parameters not provided")
	}

	return &TimeSeriesService{
		projectID:      projectID,
		send:           send,
		serviceName:    serviceName,
		serviceVersion: serviceVersion,
	}, nil
}

// TimeSeriesService provides twitter service.
type TimeSeriesService struct {
	projectID      string
	serviceName    string
	serviceVersion string
	send           bool
}

// RecordOne adds one to the specified metric without labels.
func (s *TimeSeriesService) One(ctx context.Context, name string) {
	s.OneWithLabels(ctx, name, nil)
}

// RecordOneWithLabels adds one to the specified metric with labels.
func (s *TimeSeriesService) OneWithLabels(ctx context.Context, name string, labels map[string]string) {
	s.WithLabels(ctx, name, int64(1), labels)
}

// Record adds the specified value to the specified metric.
func (s *TimeSeriesService) WithLabels(ctx context.Context, name string, value interface{}, labels map[string]string) {
	v := &monitoringpb.TypedValue{}

	switch t := value.(type) {
	case bool:
		v.Value = &monitoringpb.TypedValue_BoolValue{BoolValue: value.(bool)}
	case int64:
		v.Value = &monitoringpb.TypedValue_Int64Value{Int64Value: value.(int64)}
	case float64:
		v.Value = &monitoringpb.TypedValue_DoubleValue{DoubleValue: value.(float64)}
	case string:
		v.Value = &monitoringpb.TypedValue_StringValue{StringValue: value.(string)}
	default:
		log.Error().Msgf("unsupported metric value type: %T", t)
		return
	}

	s.record(ctx, name, v, labels)
}

// Count adds the specified value to the specified metric.
func (s *TimeSeriesService) record(ctx context.Context, name string, value *monitoringpb.TypedValue, labels map[string]string) {
	if name == "" {
		log.Error().Msg("nil metric name")
		return
	}

	mtp := fmt.Sprintf("%s/%s", metricBaseType, name)
	now := timestamppb.New(time.Now().UTC())
	lbs := map[string]string{
		// HACK: prevents time series from being overwritten for timespan which leads to errors on write
		"nanos": fmt.Sprintf("%d", now.GetNanos()),
	}

	if len(labels) > 0 {
		for k, v := range labels {
			lbs[k] = v
		}
	}

	req := &monitoringpb.CreateTimeSeriesRequest{
		Name: fmt.Sprintf("projects/%s", s.projectID),
		TimeSeries: []*monitoringpb.TimeSeries{{
			Metric: &metricpb.Metric{Type: mtp, Labels: lbs},
			Resource: &monitoredres.MonitoredResource{
				Type:   resourceType,
				Labels: map[string]string{"project_id": s.projectID},
			},
			Points: []*monitoringpb.Point{{
				Interval: &monitoringpb.TimeInterval{StartTime: now, EndTime: now},
				Value:    value,
			}},
		}},
	}

	// use the metrics context to avid it being canceled unexpectedly
	go s.post(ToMetricContext(ctx), req)
}

func (s *TimeSeriesService) post(ctx context.Context, req *monitoringpb.CreateTimeSeriesRequest) {
	if req == nil {
		log.Error().Msg("nil request in createTimeSeries")
		return
	}

	// don't create client and do not send if not configured
	if !s.send {
		return
	}

	c, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error creating metric client")
	}
	defer c.Close()

	if err = c.CreateTimeSeries(ctx, req); err != nil {
		// debug only because this is a best effort anyway
		// and there is so many errors due to too frequent metric writes
		log.Error().Err(err).Msg("error create time series")
	}
}

func (s *TimeSeriesService) GetDefaultMetricLabels() map[string]string {
	return map[string]string{
		"name":     s.serviceName,
		"version":  s.serviceVersion,
		"service":  config.GetEnv(serviceEnvVar, s.serviceName),
		"revision": config.GetEnv(revisionEnvVar, fmt.Sprintf("%s.1", s.serviceName)),
	}
}

// Count adds the specified value to the specified metric.
func (s *TimeSeriesService) Count(ctx context.Context, name string) {
	s.OneWithLabels(ctx, name, s.GetDefaultMetricLabels())
}

// CountWithLabels adds one to the specified metric with labels.
func (s *TimeSeriesService) CountWithLabels(ctx context.Context, name string, labels map[string]string) {
	l := s.GetDefaultMetricLabels()
	for k, v := range labels {
		l[k] = v
	}
	s.OneWithLabels(ctx, name, l)
}

// MeterWithLabels adds the specified value to the specified metric.
func (s *TimeSeriesService) MeterWithLabels(ctx context.Context, name string, v interface{}, labels map[string]string) {
	l := s.GetDefaultMetricLabels()
	for k, v := range labels {
		l[k] = v
	}
	s.WithLabels(ctx, name, v, l)
}

// Meter adds the specified value to the specified metric.
func (s *TimeSeriesService) Meter(ctx context.Context, name string, v interface{}) {
	l := s.GetDefaultMetricLabels()
	s.WithLabels(ctx, name, v, l)
}
