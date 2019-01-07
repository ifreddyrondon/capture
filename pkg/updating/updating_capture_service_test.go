package updating_test

//func TestUpdateCapture(t *testing.T) {
//	app, teardown := setup(t)
//	defer teardown()
//
//	e := bastion.Tester(t, app)
//	capPayload := map[string]interface{}{
//		"location": map[string]float64{
//			"lat":       1,
//			"lng":       12,
//			"elevation": 50,
//		},
//		"timestamp": "1989-12-26T06:01:00Z",
//		"tags":      []string{},
//		"payload": []map[string]interface{}{
//			{
//				"name":  "power",
//				"value": []interface{}{-70.0, -100.1, 3.1},
//			},
//		},
//	}
//
//	tt := []struct {
//		name          string
//		updatePayload map[string]interface{}
//	}{
//		{
//			"update lat",
//			map[string]interface{}{
//				"location": map[string]float64{
//					"lat":       1,
//					"lng":       capPayload["location"].(map[string]float64)["lng"],
//					"elevation": capPayload["location"].(map[string]float64)["elevation"],
//				},
//				"elevation": capPayload["elevation"],
//				"timestamp": capPayload["timestamp"],
//				"tags":      capPayload["tags"],
//				"payload":   capPayload["payload"],
//			},
//		},
//		{
//			"update lng",
//			map[string]interface{}{
//				"location": map[string]float64{
//					"lat":       capPayload["location"].(map[string]float64)["lat"],
//					"lng":       30,
//					"elevation": capPayload["location"].(map[string]float64)["elevation"],
//				},
//				"timestamp": capPayload["timestamp"],
//				"tags":      capPayload["tags"],
//				"payload":   capPayload["payload"],
//			},
//		},
//		{
//			"update elevation",
//			map[string]interface{}{
//				"location": map[string]float64{
//					"lat":       capPayload["location"].(map[string]float64)["lat"],
//					"lng":       capPayload["location"].(map[string]float64)["lng"],
//					"elevation": 100,
//				},
//				"timestamp": capPayload["timestamp"],
//				"tags":      capPayload["tags"],
//				"payload":   capPayload["payload"],
//			},
//		},
//		{
//			"update timestamp",
//			map[string]interface{}{
//				"location":  capPayload["location"],
//				"timestamp": "2006-07-12T06:01:00Z",
//				"tags":      capPayload["tags"],
//				"payload":   capPayload["payload"],
//			},
//		},
//		{
//			"update payload",
//			map[string]interface{}{
//				"location":  capPayload["location"],
//				"timestamp": capPayload["timestamp"],
//				"tags":      capPayload["tags"],
//				"payload": []map[string]interface{}{
//					{
//						"name":  "power",
//						"value": []interface{}{1},
//					},
//				},
//			},
//		},
//		{
//			"update tags",
//			map[string]interface{}{
//				"location":  capPayload["location"],
//				"timestamp": capPayload["timestamp"],
//				"payload":   capPayload["payload"],
//				"tags":      []string{"tag1", "tag2"},
//			},
//		},
//		{
//			"do not update id",
//			map[string]interface{}{
//				"id":        "123",
//				"location":  capPayload["location"],
//				"timestamp": capPayload["timestamp"],
//				"tags":      capPayload["tags"],
//				"payload":   capPayload["payload"],
//			},
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			createdObj := e.POST("/test").WithJSON(capPayload).Expect().
//				Status(http.StatusCreated).JSON().Object().Raw()
//
//			e.GET(fmt.Sprintf("/%v", createdObj["id"])).Expect().Status(http.StatusOK)
//			tc.updatePayload["id"] = createdObj["id"]
//
//			updatedObj := e.PUT(fmt.Sprintf("/%v", createdObj["id"])).WithJSON(tc.updatePayload).Expect().
//				Status(http.StatusOK).
//				JSON().Object().
//				ContainsKey("id").ValueEqual("id", tc.updatePayload["id"]).
//				ContainsKey("location").ValueEqual("location", tc.updatePayload["location"]).
//				ContainsKey("timestamp").ValueEqual("timestamp", tc.updatePayload["timestamp"]).
//				ContainsKey("payload").ValueEqual("payload", tc.updatePayload["payload"]).
//				ContainsKey("tags").ValueEqual("tags", tc.updatePayload["tags"]).
//				ContainsKey("createdAt").NotEmpty().
//				ContainsKey("updatedAt").NotEmpty().
//				Raw()
//
//			// updatedAt from put should be after updatedAt from post
//			updatedAtFromCreate, err := dateparse.ParseAny(createdObj["updatedAt"].(string))
//			assert.Nil(t, err)
//			updatedAtFromUpdate, err := dateparse.ParseAny(updatedObj["updatedAt"].(string))
//			assert.Nil(t, err)
//			assert.True(t, updatedAtFromUpdate.After(updatedAtFromCreate))
//		})
//	}
//}
//
//func TestUpdateCaptureFailsBadRequest(t *testing.T) {
//	app, teardown := setup(t)
//	defer teardown()
//
//	e := bastion.Tester(t, app)
//	capPayload := map[string]interface{}{
//		"location": map[string]float64{
//			"lat": 1,
//			"lng": 12,
//		},
//		"timestamp": "1989-12-26T06:01:00Z",
//		"payload": []map[string]interface{}{
//			{
//				"name":  "power",
//				"value": []interface{}{-70.0, -100.1, 3.1},
//			},
//		},
//	}
//
//	tt := []struct {
//		name          string
//		updatePayload map[string]interface{}
//		response      map[string]interface{}
//	}{
//		{
//			"lat out of range",
//			map[string]interface{}{
//				"location": map[string]float64{
//					"lat": 200,
//					"lng": capPayload["location"].(map[string]float64)["lng"],
//				},
//				"timestamp": capPayload["timestamp"],
//				"payload":   capPayload["payload"],
//			},
//			map[string]interface{}{
//				"status":  400.0,
//				"error":   "Bad Request",
//				"message": "latitude out of boundaries, may range from -90.0 to 90.0",
//			},
//		},
//		{
//			"missing payload",
//			map[string]interface{}{
//				"location": map[string]float64{
//					"lat": capPayload["location"].(map[string]float64)["lat"],
//					"lng": capPayload["location"].(map[string]float64)["lng"],
//				},
//				"timestamp": capPayload["timestamp"],
//			},
//			map[string]interface{}{
//				"status":  400.0,
//				"error":   "Bad Request",
//				"message": "payload value must not be blank",
//			},
//		},
//	}
//
//	for _, tc := range tt {
//		t.Run(tc.name, func(t *testing.T) {
//			createdObj := e.POST("/test").WithJSON(capPayload).Expect().
//				Status(http.StatusCreated).JSON().Object().Raw()
//			e.GET(fmt.Sprintf("/%v", createdObj["id"])).Expect().Status(http.StatusOK)
//			tc.updatePayload["id"] = createdObj["id"]
//			e.PUT(fmt.Sprintf("/%v", createdObj["id"])).WithJSON(tc.updatePayload).Expect().
//				Status(http.StatusBadRequest).
//				JSON().Object().Equal(tc.response)
//		})
//	}
//}
//
//func TestUpdateCaptureFailsBadRequestWhenMissingID(t *testing.T) {
//	app, teardown := setup(t)
//	defer teardown()
//
//	e := bastion.Tester(t, app)
//	capPayload := map[string]interface{}{
//		"location": map[string]float64{
//			"lat": 1,
//			"lng": 12,
//		},
//		"timestamp": "1989-12-26T06:01:00Z",
//		"payload": []map[string]interface{}{
//			{
//				"name":  "power",
//				"value": []interface{}{-70.0, -100.1, 3.1},
//			},
//		},
//	}
//
//	response := map[string]interface{}{
//		"status":  400.0,
//		"error":   "Bad Request",
//		"message": "capture id must not be blank",
//	}
//
//	createdObj := e.POST("/test").WithJSON(capPayload).Expect().
//		Status(http.StatusCreated).JSON().Object().Raw()
//	e.GET(fmt.Sprintf("/%v", createdObj["id"])).Expect().Status(http.StatusOK)
//
//	e.PUT(fmt.Sprintf("/%v", createdObj["id"])).WithJSON(capPayload).Expect().
//		Status(http.StatusBadRequest).
//		JSON().Object().Equal(response)
//}
