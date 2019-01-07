package domain_test

//func TestEmptyBranch(t *testing.T) {
//	t.Parallel()
//
//	p := []byte(`[]`)
//	var c pkg.Captures
//	err := c.UnmarshalJSON(p)
//	require.EqualError(t, err, "cannot unmarshal json into valid captures, it needs at least one valid capture")
//}
//
//func TestCapturesOKUnmarshalJSON(t *testing.T) {
//	t.Parallel()
//
//	payl := payload.Payload{
//		&payload.Metric{Name: "power", Value: []interface{}{-70.0, -100.1, 3.1}},
//	}
//	tt := []struct {
//		name    string
//		payload []byte
//		result  capture.Captures
//	}{
//		{
//			"captures of len 1",
//			[]byte(`[{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}], "lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
//			capture.Captures{getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 1)},
//		},
//		{
//			"path of len 2",
//			[]byte(`[
//						{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}], "lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
//						{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}], "lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
//			capture.Captures{getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 1), getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 2)},
//		},
//		{
//			"invalid capture into path of len 2",
//			[]byte(`[
//						{"payload":[{"name": "power", "value": [-70, -100.1, 3.1]}], "lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"},
//						{"lat": 1, "lng": 2, "date": "1989-12-26T06:01:00.00Z"}]`),
//			capture.Captures{getCapture(payl, "1989-12-26T06:01:00.00Z", 1, 2)},
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			var c capture.Captures
//			err := c.UnmarshalJSON(tc.payload)
//			require.Nil(t, err)
//			assert.Len(t, c, len(tc.result))
//		})
//	}
//}
//
//func TestCapturesBadUnmarshalJSON(t *testing.T) {
//	t.Parallel()
//
//	expectedErr := "cannot unmarshal json into valid captures, it needs at least one valid capture"
//
//	tt := []struct {
//		name    string
//		payload []byte
//	}{
//		{
//			"invalid capture (lat) into path of len 1",
//			[]byte(`[{"lat": -101, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
//		},
//		{
//			"invalid capture (missing payload) into path of len 1",
//			[]byte(`[{"lat": 1, "lng": 1, "date": "1989-12-26T06:01:00.00Z"}]`),
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			var c capture.Captures
//			err := c.UnmarshalJSON(tc.payload)
//			assert.EqualError(t, err, expectedErr)
//		})
//	}
//}

//func TestUnmarshalCapturesFail(t *testing.T) {
//	app, teardown := setup(t)
//	defer teardown()
//
//	e := bastion.Tester(t, app)
//	tt := []struct {
//		name     string
//		payload  []interface{}
//		response map[string]interface{}
//	}{
//		{
//			name:    "bad request, missing body",
//			payload: []interface{}{},
//			response: map[string]interface{}{
//				"status":  400.0,
//				"error":   "Bad Request",
//				"message": "cannot unmarshal json into valid captures, it needs at least one valid capture",
//			},
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			e.POST("/test").
//				WithJSON(tc.payload).
//				Expect().
//				Status(http.StatusBadRequest).
//				JSON().Object().Equal(tc.response)
//		})
//	}
