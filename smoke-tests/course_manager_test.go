package smoke_tests

import (
	"net/http"
	"testing"
)

func TestCourseManager_HappyPath(t *testing.T) {

	rp := RequestParams{
		BaseUrl: "http://localhost:8000/v1",
	}
	err := rp.Do()
	if err != nil {
		t.Skip("course manager service is not running, skipping the smoke tests")
	}

	rp.Payload = map[string]interface{}{
		"course": map[string]interface{}{
			"name": "Microservices with Go",
			"tutor": map[string]string{
				"name":       "John",
				"lastname":   "Stone",
				"uuid":       "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				"faculty":    "Computer Science",
				"lecturerOf": "Golang",
			},
		},
	}

	rp.Path = "/createCourse"
	rp.Method = http.MethodPost
	err = rp.Do()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rp.StatusCode != http.StatusCreated {
		t.Errorf("expected %v, got %v", http.StatusCreated, rp.StatusCode)
		t.Log(rp.ResponseBody)
	}

	courseUUID := rp.ResponseBody.(map[string]interface{})["uuid"].(string)
	rp.Path = "/getCourse/" + courseUUID
	rp.Method = http.MethodGet
	rp.Payload = map[string]interface{}{}
	err = rp.Do()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rp.StatusCode != http.StatusOK && rp.ResponseBody.(map[string]interface{})["uuid"].(string) != courseUUID {
		t.Errorf("expected %v, got %v", http.StatusOK, rp.StatusCode)
		t.Log(rp.ResponseBody)
	}

	rp.Path = "/listCourses"
	rp.Method = http.MethodGet
	err = rp.Do()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rp.StatusCode != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, rp.StatusCode)
		t.Log(rp.ResponseBody)
	}
	var exists bool
	for _, u := range rp.ResponseBody.([]interface{}) {
		if u.(map[string]interface{})["uuid"].(string) == courseUUID {
			exists = true
		}
	}
	if !exists {
		t.Errorf("course list does not contain expected course %v", courseUUID)
	}

	rp.Path = "/registerStudent/" + courseUUID
	rp.Method = http.MethodPut
	rp.Payload = map[string]interface{}{
		"student": map[string]string{
			"name":     "Alice J",
			"lastname": "Smith",
			"uuid":     "3fa85f64-5717-4562-b3fc-2c963f66afa7",
			"faculty":  "Computer Science",
		},
	}
	err = rp.Do()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rp.StatusCode != http.StatusNoContent {
		t.Errorf("expected %v, got %v", http.StatusNoContent, rp.StatusCode)
		t.Log(rp.ResponseBody)
	}

	rp.Path = "/unregisterStudent/" + courseUUID
	rp.Method = http.MethodPut
	rp.Payload = map[string]interface{}{
		"studentUUID": "3fa85f64-5717-4562-b3fc-2c963f66afa7",
	}
	err = rp.Do()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rp.StatusCode != http.StatusNoContent {
		t.Errorf("expected %v, got %v", http.StatusNoContent, rp.StatusCode)
		t.Log(rp.ResponseBody)
	}

	rp.Path = "/deleteCourse/" + courseUUID
	rp.Method = http.MethodDelete
	rp.Payload = map[string]interface{}{}
	err = rp.Do()
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if rp.StatusCode != http.StatusNoContent {
		t.Errorf("expected %v, got %v", http.StatusNoContent, rp.StatusCode)
		t.Log(rp.ResponseBody)
	}
}
