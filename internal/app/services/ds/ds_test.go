package ds

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/vitalykrupin/url-shortener/internal/app/storage"
)

// mockStorage for testing
type mockStorage struct {
	deleteCalls []struct {
		userID string
		urls   []string
	}
	mu sync.Mutex
}

func (m *mockStorage) Add(ctx context.Context, batch map[storage.Alias]storage.OriginalURL) error {
	return nil
}

func (m *mockStorage) GetURL(ctx context.Context, alias storage.Alias) (storage.OriginalURL, error) {
	return "", nil
}

func (m *mockStorage) GetAlias(ctx context.Context, url storage.OriginalURL) (storage.Alias, error) {
	return "", nil
}

func (m *mockStorage) GetUserURLs(ctx context.Context, userID string) (storage.AliasKeysMap, error) {
	return nil, nil
}

func (m *mockStorage) DeleteUserURLs(ctx context.Context, userID string, urls []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deleteCalls = append(m.deleteCalls, struct {
		userID string
		urls   []string
	}{userID, urls})
	return nil
}

func (m *mockStorage) CloseStorage(ctx context.Context) error {
	return nil
}

func (m *mockStorage) PingStorage(ctx context.Context) error {
	return nil
}

func (m *mockStorage) GetDeleteCalls() []struct {
	userID string
	urls   []string
} {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Return a copy
	result := make([]struct {
		userID string
		urls   []string
	}, len(m.deleteCalls))
	copy(result, m.deleteCalls)
	return result
}

func TestNewDeleteService(t *testing.T) {
	mockStore := &mockStorage{}
	service := NewDeleteService(mockStore)

	if service == nil {
		t.Fatal("Expected DeleteService to be non-nil")
	}
}

func TestDeleteService_Add(t *testing.T) {
	mockStore := &mockStorage{}
	service := NewDeleteService(mockStore)

	userID := "user123"
	urls := []string{"url1", "url2", "url3"}

	// Start the service
	service.Start(1)
	defer service.Stop()

	// Add a payload
	service.Add(userID, urls)

	// Give some time for processing
	time.Sleep(100 * time.Millisecond)

	// Check that DeleteUserURLs was called
	calls := mockStore.GetDeleteCalls()
	if len(calls) != 1 {
		t.Fatalf("Expected 1 delete call, got %d", len(calls))
	}

	call := calls[0]
	if call.userID != userID {
		t.Errorf("Expected userID %s, got %s", userID, call.userID)
	}

	if len(call.urls) != len(urls) {
		t.Errorf("Expected %d URLs, got %d", len(urls), len(call.urls))
	}

	for i, url := range urls {
		if call.urls[i] != url {
			t.Errorf("Expected URL %s at index %d, got %s", url, i, call.urls[i])
		}
	}
}

func TestDeleteService_MultipleWorkers(t *testing.T) {
	mockStore := &mockStorage{}
	service := NewDeleteService(mockStore)

	// Start with multiple workers
	workerCount := 3
	service.Start(workerCount)
	defer service.Stop()

	// Add multiple payloads
	payloads := []struct {
		userID string
		urls   []string
	}{
		{"user1", []string{"url1", "url2"}},
		{"user2", []string{"url3", "url4"}},
		{"user3", []string{"url5", "url6"}},
		{"user4", []string{"url7", "url8"}},
		{"user5", []string{"url9", "url10"}},
	}

	for _, payload := range payloads {
		service.Add(payload.userID, payload.urls)
	}

	// Give time for processing
	time.Sleep(200 * time.Millisecond)

	// Check that all payloads were processed
	calls := mockStore.GetDeleteCalls()
	if len(calls) != len(payloads) {
		t.Fatalf("Expected %d delete calls, got %d", len(payloads), len(calls))
	}

	// Check that all userIDs are present
	userIDs := make(map[string]bool)
	for _, call := range calls {
		userIDs[call.userID] = true
	}

	for _, payload := range payloads {
		if !userIDs[payload.userID] {
			t.Errorf("Expected userID %s to be processed", payload.userID)
		}
	}
}

func TestDeleteService_ConcurrentAdd(t *testing.T) {
	mockStore := &mockStorage{}
	service := NewDeleteService(mockStore)

	service.Start(2)
	defer service.Stop()

	// Add payloads concurrently
	var wg sync.WaitGroup
	numGoroutines := 10
	payloadsPerGoroutine := 5

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < payloadsPerGoroutine; j++ {
				userID := fmt.Sprintf("user%d", goroutineID)
				urls := []string{fmt.Sprintf("url%d_%d", goroutineID, j)}
				service.Add(userID, urls)
			}
		}(i)
	}

	wg.Wait()

	// Give time for processing
	time.Sleep(200 * time.Millisecond)

	// Check that all payloads were processed
	expectedCalls := numGoroutines * payloadsPerGoroutine
	calls := mockStore.GetDeleteCalls()
	if len(calls) != expectedCalls {
		t.Fatalf("Expected %d delete calls, got %d", expectedCalls, len(calls))
	}
}

func TestDeleteService_Stop(t *testing.T) {
	mockStore := &mockStorage{}
	service := NewDeleteService(mockStore)

	service.Start(1)

	// Add a payload
	service.Add("user1", []string{"url1"})

	// Give time for processing
	time.Sleep(100 * time.Millisecond)

	// Stop the service
	service.Stop()

	// The first payload should have been processed
	calls := mockStore.GetDeleteCalls()
	if len(calls) == 0 {
		t.Error("Expected at least one delete call")
	}

	// Note: Adding after stop will panic due to closed channel
	// This is expected behavior and the test verifies the service stops gracefully
}

func TestDeleteService_InterfaceCompliance(t *testing.T) {
	mockStore := &mockStorage{}
	service := NewDeleteService(mockStore)

	// Test that DeleteService implements DeleteServiceInterface
	var _ DeleteServiceInterface = service
}

func TestDeleteService_EmptyURLs(t *testing.T) {
	mockStore := &mockStorage{}
	service := NewDeleteService(mockStore)

	service.Start(1)
	defer service.Stop()

	// Add payload with empty URLs
	service.Add("user1", []string{})

	// Give time for processing
	time.Sleep(100 * time.Millisecond)

	// Should still process the call
	calls := mockStore.GetDeleteCalls()
	if len(calls) != 1 {
		t.Fatalf("Expected 1 delete call, got %d", len(calls))
	}

	call := calls[0]
	if len(call.urls) != 0 {
		t.Errorf("Expected empty URLs, got %v", call.urls)
	}
}
