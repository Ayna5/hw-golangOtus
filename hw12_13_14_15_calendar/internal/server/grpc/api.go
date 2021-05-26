package grpc

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/internal/storage"
	calendar_pb "github.com/Ayna5/hw-golangOtus/hw12_13_14_15_calendar/pkg/calendar"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (a *Server) CreateEvent(ctx context.Context, req *calendar_pb.CreateEventRequest) (*calendar_pb.CreateEventResponse, error) { //nolint:stylecheck
	event := convertToStorageEvent(req.Event)
	if err := a.api.CreateEvent(event); err != nil {
		return nil, fmt.Errorf("cannot create event: %w", err)
	}
	return &calendar_pb.CreateEventResponse{}, nil
}

func (a *Server) UpdateEvent(ctx context.Context, req *calendar_pb.UpdateEventRequest) (*calendar_pb.UpdateEventResponse, error) {
	event := convertToStorageEvent(req.Event)
	if err := a.api.UpdateEvent(event); err != nil {
		return nil, fmt.Errorf("cannot update event: %w", err)
	}
	return &calendar_pb.UpdateEventResponse{}, nil
}

func (a *Server) GetEvents(ctx context.Context, req *calendar_pb.GetEventsRequest) (*calendar_pb.GetEventsResponse, error) {
	events, err := a.api.GetEvents(req.StartData.AsTime(), req.EndData.AsTime())
	if err != nil {
		return nil, fmt.Errorf("cannot get events: %w", err)
	}
	var newEvents []*calendar_pb.Event //nolint:prealloc
	for _, val := range events {
		var newEvent *calendar_pb.Event
		newEvent, err = convertToPBEvent(*val)
		if err != nil {
			return nil, fmt.Errorf("cannot convert to event: %w", err)
		}
		newEvents = append(newEvents, newEvent)
	}
	return &calendar_pb.GetEventsResponse{Event: newEvents}, nil
}

func (a *Server) DeleteEvent(ctx context.Context, req *calendar_pb.DeleteEventRequest) (*calendar_pb.DeleteEventResponse, error) {
	event := convertToStorageEvent(req.Event)
	if err := a.api.DeleteEvent(event); err != nil {
		return nil, fmt.Errorf("cannot delete event: %w", err)
	}
	return &calendar_pb.DeleteEventResponse{}, nil
}

func convertToStorageEvent(e *calendar_pb.Event) storage.Event {
	return storage.Event{
		ID:          strconv.Itoa(int(e.Id)),
		Title:       e.Title,
		StartData:   e.StartData.AsTime(),
		EndData:     e.EndData.AsTime(),
		Description: e.Description,
		OwnerID:     strconv.Itoa(int(e.OwnerId)),
		RemindIn:    e.RemindIn,
	}
}

func convertToPBEvent(e storage.Event) (*calendar_pb.Event, error) {
	id, err := strconv.Atoi(e.ID)
	if err != nil {
		return nil, fmt.Errorf("cannot convert to int: %w", err)
	}
	ownerID, err := strconv.Atoi(e.OwnerID)
	if err != nil {
		return nil, fmt.Errorf("cannot convert to int: %w", err)
	}
	return &calendar_pb.Event{
		Id:          uint64(id),
		Title:       e.Title,
		StartData:   &timestamppb.Timestamp{Seconds: int64(e.StartData.Second())},
		EndData:     &timestamppb.Timestamp{Seconds: int64(e.EndData.Second())},
		Description: e.Description,
		OwnerId:     uint64(ownerID),
		RemindIn:    e.RemindIn,
	}, nil
}
