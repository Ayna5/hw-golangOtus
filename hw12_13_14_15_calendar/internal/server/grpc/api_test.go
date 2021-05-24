package grpc

//var (
//	ctx    = context.Background()
//	client calendar_pb.CalendarClient
//)
//
//func TestServerGRPC(t *testing.T) {
//	host := "localhost:14545"
//	conn, err := grpc.Dial(host, grpc.WithInsecure())
//	if err != nil {
//		fmt.Println(err)
//	}
//	defer conn.Close()
//
//	client = calendar_pb.NewCalendarClient(conn)
//	t.Run("Create event", func(t *testing.T) {
//		request := &calendar_pb.CreateEventRequest{Event: &calendar_pb.Event{
//			Id:          1,
//			Title:       "event",
//			StartData:   timestamppb.Now(),
//			EndData:     timestamppb.New(time.Now().AddDate(0, 0, 2)),
//			Description: "event",
//			OwnerId:     1,
//			RemindIn:    "2",
//		}}
//		response, err := client.CreateEvent(ctx, request)
//
//		require.NoError(t, err)
//		require.NotNil(t, response)
//	})
//	t.Run("Update event", func(t *testing.T) {
//		request := &calendar_pb.UpdateEventRequest{Event: &calendar_pb.Event{
//			Id:          1,
//			Title:       "event new",
//			StartData:   timestamppb.Now(),
//			EndData:     timestamppb.New(time.Now().AddDate(0, 0, 20)),
//			Description: "event new",
//			OwnerId:     1,
//			RemindIn:    "2",
//		}}
//		response, err := client.UpdateEvent(ctx, request)
//
//		require.NoError(t, err)
//		require.NotNil(t, response)
//	})
//	t.Run("Get events", func(t *testing.T) {
//		request := &calendar_pb.GetEventsRequest{
//			StartData:   timestamppb.New(time.Now().Add(-1*time.Hour)),
//			EndData:     timestamppb.New(time.Now().AddDate(0, 0, 3)),
//		}
//		response, err := client.GetEvents(ctx, request)
//
//		require.NoError(t, err)
//		require.NotNil(t, response)
//		require.Equal(t, uint64(1), response.Event[0].Id)
//	})
//	t.Run("Delete event", func(t *testing.T) {
//		request := &calendar_pb.DeleteEventRequest{Event: &calendar_pb.Event{
//			Id:          1,
//			Title:       "event new",
//			StartData:   timestamppb.Now(),
//			EndData:     timestamppb.New(time.Now().AddDate(0, 0, 20)),
//			Description: "event new",
//			OwnerId:     1,
//			RemindIn:    "2",
//		}}
//		response, err := client.DeleteEvent(ctx, request)
//
//		require.NoError(t, err)
//		require.NotNil(t, response)
//	})
//}
