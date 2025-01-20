package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRegisterDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockDeviceRepository(ctrl)
	device := Device{
		DeviceName:   "Device1",
		DeviceType:   "TypeA",
		Brand:        "BrandX",
		Model:        "ModelY",
		Os:           "OSZ",
		OsVersion:    "1.0",
		PurchaseDate: "2022-01-01",
		WarrantyEnd:  "2024-01-01",
		Status:       "Active",
		Price:        1000,
	}

	mockRepo.EXPECT().Create(&device).Return(nil)

	body, _ := json.Marshal(device)
	req, _ := http.NewRequest(http.MethodPost, "/device", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestUpdateDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockDeviceRepository(ctrl)
	updatedDevice := Device{
		ID:           1,
		DeviceName:   "UpdatedDevice",
		DeviceType:   "TypeA",
		Brand:        "BrandX",
		Model:        "ModelY",
		Os:           "OSZ",
		OsVersion:    "1.0",
		PurchaseDate: "2022-01-01",
		WarrantyEnd:  "2024-01-01",
		Status:       "Active",
		Price:        1200,
	}

	mockRepo.EXPECT().Update(&updatedDevice).Return(nil)

	body, _ := json.Marshal(updatedDevice)
	req, _ := http.NewRequest(http.MethodPut, "/device/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListDevices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockDeviceRepository(ctrl)
	devices := []Device{
		{ID: 1, DeviceName: "Device1"},
		{ID: 2, DeviceName: "Device2"},
	}

	mockRepo.EXPECT().FindAll(10, 0).Return(devices, nil)

	req, _ := http.NewRequest(http.MethodGet, "/device?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []Device
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, devices, response)
}

func TestGetDeviceByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockDeviceRepository(ctrl)
	device := Device{
		ID:           1,
		DeviceName:   "Device1",
		DeviceType:   "TypeA",
		Brand:        "BrandX",
		Model:        "ModelY",
		Os:           "OSZ",
		OsVersion:    "1.0",
		PurchaseDate: "2022-01-01",
		WarrantyEnd:  "2024-01-01",
		Status:       "Active",
		Price:        1000,
	}

	mockRepo.EXPECT().FindByID(uint(1)).Return(&device, nil)

	req, _ := http.NewRequest(http.MethodGet, "/device/1", nil)
	w := httptest.NewRecorder()

	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response Device
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, device, response)
}

func TestDeleteDevice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockDeviceRepository(ctrl)

	mockRepo.EXPECT().Delete(uint(1)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/device/1", nil)
	w := httptest.NewRecorder()

	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRegisterDeviceInvalidInput(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/device", bytes.NewBuffer([]byte(`{"invalid":"data"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDeviceByIDNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockDeviceRepository(ctrl)
	mockRepo.EXPECT().FindByID(uint(999)).Return(nil, gorm.ErrRecordNotFound)

	req, _ := http.NewRequest(http.MethodGet, "/device/999", nil)
	w := httptest.NewRecorder()

	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListDevicesEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockDeviceRepository(ctrl)
	mockRepo.EXPECT().FindAll(10, 0).Return([]Device{}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/device?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	r := setupRouter()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []Device
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Empty(t, response)
}
