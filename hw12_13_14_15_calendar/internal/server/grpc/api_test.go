package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/server/grpc/mocks"
	calendar_pb "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/pkg/calendar"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	ctx    = context.Background()
)

func TestServerGRPC(t *testing.T) {
	svc := mocks.NewCalendarClientMock(t)

	t.Run("Create event", func(t *testing.T) {
		request := &calendar_pb.CreateEventRequest{Event: &calendar_pb.Event{
			Id:          1,
			Title:       "event",
			StartData:   timestamppb.Now(),
			EndData:     timestamppb.New(time.Now().AddDate(0, 0, 2)),
			Description: "event",
			OwnerId:     1,
			RemindIn:    "2",
		}}

		svc.CreateEventMock.Set(func(ctx context.Context, in *calendar_pb.CreateEventRequest, opts ...grpc.CallOption) (cp1 *calendar_pb.CreateEventResponse, err error) {
			return &calendar_pb.CreateEventResponse{}, nil
		})

		response, err := svc.CreateEvent(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
	})
	t.Run("Update event", func(t *testing.T) {
		request := &calendar_pb.UpdateEventRequest{Event: &calendar_pb.Event{
			Id:          1,
			Title:       "event new",
			StartData:   timestamppb.Now(),
			EndData:     timestamppb.New(time.Now().AddDate(0, 0, 20)),
			Description: "event new",
			OwnerId:     1,
			RemindIn:    "2",
		}}

		svc.UpdateEventMock.Set(func(ctx context.Context, in *calendar_pb.UpdateEventRequest, opts ...grpc.CallOption) (up1 *calendar_pb.UpdateEventResponse, err error) {
			return &calendar_pb.UpdateEventResponse{}, nil
		})

		response, err := svc.UpdateEvent(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
	})
	t.Run("Get events", func(t *testing.T) {
		request := &calendar_pb.GetEventsRequest{
			StartData: timestamppb.New(time.Now().Add(-1 * time.Hour)),
			EndData:   timestamppb.New(time.Now().AddDate(0, 0, 3)),
		}

		svc.GetEventsMock.Set(func(ctx context.Context, in *calendar_pb.GetEventsRequest, opts ...grpc.CallOption) (gp1 *calendar_pb.GetEventsResponse, err error) {
			return &calendar_pb.GetEventsResponse{Event: []*calendar_pb.Event{{Id: 1}}}, nil
		})

		response, err := svc.GetEvents(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
		require.Equal(t, uint64(1), response.Event[0].Id)
	})
	t.Run("Delete event", func(t *testing.T) {
		request := &calendar_pb.DeleteEventRequest{Event: &calendar_pb.Event{
			Id:          1,
			Title:       "event new",
			StartData:   timestamppb.Now(),
			EndData:     timestamppb.New(time.Now().AddDate(0, 0, 20)),
			Description: "event new",
			OwnerId:     1,
			RemindIn:    "2",
		}}

		svc.DeleteEventMock.Set(func(ctx context.Context, in *calendar_pb.DeleteEventRequest, opts ...grpc.CallOption) (dp1 *calendar_pb.DeleteEventResponse, err error) {
			return &calendar_pb.DeleteEventResponse{}, nil
		})

		response, err := svc.DeleteEvent(ctx, request)

		require.NoError(t, err)
		require.NotNil(t, response)
	})
}
