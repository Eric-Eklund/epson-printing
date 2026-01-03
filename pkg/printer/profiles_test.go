package printer

import (
	"testing"
)

func TestGetPrintOptions(t *testing.T) {
	tests := []struct {
		name          string
		profile       PrintProfile
		expectError   bool
		expectedSize  string
		expectedTray  string
		expectedMedia string
		expectedQual  int
	}{
		{
			name:          "Default profile",
			profile:       ProfileDefault,
			expectError:   false,
			expectedSize:  "A4",
			expectedTray:  "Main",
			expectedMedia: "stationery",
			expectedQual:  3,
		},
		{
			name:          "Empty profile returns default",
			profile:       "",
			expectError:   false,
			expectedSize:  "A4",
			expectedTray:  "Main",
			expectedMedia: "stationery",
			expectedQual:  3,
		},
		{
			name:          "Photo 4x6 borderless glossy",
			profile:       ProfilePhoto4x6BorderlessGlossy,
			expectError:   false,
			expectedSize:  "4x6.Borderless",
			expectedTray:  "Photo",
			expectedMedia: "photographic-glossy",
			expectedQual:  5,
		},
		{
			name:          "Photo 4x6 borderless matte",
			profile:       ProfilePhoto4x6BorderlessMatte,
			expectError:   false,
			expectedSize:  "4x6.Borderless",
			expectedTray:  "Photo",
			expectedMedia: "photographic-matte",
			expectedQual:  5,
		},
		{
			name:          "Photo A3+ borderless matte",
			profile:       ProfilePhotoA3PlusBorderlessMatte,
			expectError:   false,
			expectedSize:  "13x19.Borderless",
			expectedTray:  "Rear",
			expectedMedia: "photographic-matte",
			expectedQual:  5,
		},
		{
			name:          "Photo A3 borderless glossy",
			profile:       ProfilePhotoA3BorderlessGlossy,
			expectError:   false,
			expectedSize:  "A3.Borderless",
			expectedTray:  "Rear",
			expectedMedia: "photographic-glossy",
			expectedQual:  5,
		},
		{
			name:          "Photo A4 borderless semigloss",
			profile:       ProfilePhotoA4BorderlessSemiGloss,
			expectError:   false,
			expectedSize:  "A4.Borderless",
			expectedTray:  "Auto",
			expectedMedia: "photographic-semi-gloss",
			expectedQual:  5,
		},
		{
			name:          "Document draft",
			profile:       ProfileDocumentDraft,
			expectError:   false,
			expectedSize:  "A4",
			expectedTray:  "Main",
			expectedMedia: "stationery",
			expectedQual:  3,
		},
		{
			name:          "Document normal",
			profile:       ProfileDocumentNormal,
			expectError:   false,
			expectedSize:  "A4",
			expectedTray:  "Main",
			expectedMedia: "stationery",
			expectedQual:  4,
		},
		{
			name:          "Document best",
			profile:       ProfileDocumentBest,
			expectError:   false,
			expectedSize:  "A4",
			expectedTray:  "Main",
			expectedMedia: "stationery-coated",
			expectedQual:  5,
		},
		{
			name:        "Unknown profile returns error",
			profile:     "unknown-profile",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, err := GetPrintOptions(tt.profile)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if opts.PaperSize != tt.expectedSize {
				t.Errorf("Expected paper size %s, got %s", tt.expectedSize, opts.PaperSize)
			}

			if opts.Tray != tt.expectedTray {
				t.Errorf("Expected tray %s, got %s", tt.expectedTray, opts.Tray)
			}

			if opts.MediaType != tt.expectedMedia {
				t.Errorf("Expected media type %s, got %s", tt.expectedMedia, opts.MediaType)
			}

			if opts.Quality != tt.expectedQual {
				t.Errorf("Expected quality %d, got %d", tt.expectedQual, opts.Quality)
			}

			// All profiles should have sensible defaults
			if opts.PageRange != "all" {
				t.Errorf("Expected page range 'all', got %s", opts.PageRange)
			}

			if opts.Copies != 1 {
				t.Errorf("Expected copies 1, got %d", opts.Copies)
			}
		})
	}
}

func TestMustGetPrintOptions(t *testing.T) {
	// Should not panic for valid profile
	opts := MustGetPrintOptions(ProfileDefault)
	if opts.PaperSize != "A4" {
		t.Errorf("Expected A4, got %s", opts.PaperSize)
	}

	// Should panic for invalid profile
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid profile, but didn't get one")
		}
	}()
	MustGetPrintOptions("invalid-profile")
}

func TestListProfiles(t *testing.T) {
	profiles := ListProfiles()

	if len(profiles) == 0 {
		t.Error("Expected at least one profile")
	}

	// Should contain some known profiles
	expectedProfiles := []PrintProfile{
		ProfileDefault,
		ProfilePhoto4x6BorderlessGlossy,
		ProfilePhotoA3PlusBorderlessMatte,
		ProfileDocumentDraft,
	}

	profileMap := make(map[PrintProfile]bool)
	for _, p := range profiles {
		profileMap[p] = true
	}

	for _, expected := range expectedProfiles {
		if !profileMap[expected] {
			t.Errorf("Expected profile %s not found in list", expected)
		}
	}
}

func TestRegisterProfile(t *testing.T) {
	customProfile := PrintProfile("custom-test-profile")
	customOpts := PrintOptions{
		PaperSize: "A5",
		Tray:      "Auto",
		MediaType: "stationery",
		Quality:   4,
		PageRange: "all",
		Copies:    2,
	}

	// Register custom profile
	RegisterProfile(customProfile, customOpts)

	// Should be able to retrieve it
	opts, err := GetPrintOptions(customProfile)
	if err != nil {
		t.Fatalf("Failed to get custom profile: %v", err)
	}

	if opts.PaperSize != "A5" {
		t.Errorf("Expected A5, got %s", opts.PaperSize)
	}

	if opts.Copies != 2 {
		t.Errorf("Expected 2 copies, got %d", opts.Copies)
	}
}

func TestGetProfileDescription(t *testing.T) {
	tests := []struct {
		name     string
		profile  PrintProfile
		contains []string
	}{
		{
			name:     "Photo 4x6 glossy",
			profile:  ProfilePhoto4x6BorderlessGlossy,
			contains: []string{"4x6.Borderless", "photographic-glossy", "Photo", "quality: 5"},
		},
		{
			name:     "Document draft",
			profile:  ProfileDocumentDraft,
			contains: []string{"A4", "stationery", "Main", "quality: 3"},
		},
		{
			name:     "A3+ matte",
			profile:  ProfilePhotoA3PlusBorderlessMatte,
			contains: []string{"13x19.Borderless", "photographic-matte", "Rear", "quality: 5"},
		},
		{
			name:     "Unknown profile",
			profile:  "unknown",
			contains: []string{"Unknown profile"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := GetProfileDescription(tt.profile)

			for _, expected := range tt.contains {
				if !contains(desc, expected) {
					t.Errorf("Expected description to contain '%s', got: %s", expected, desc)
				}
			}
		})
	}
}

func TestAllPhotoProfilesUseBestQuality(t *testing.T) {
	photoProfiles := []PrintProfile{
		ProfilePhoto4x6BorderlessGlossy,
		ProfilePhoto4x6BorderlessMatte,
		ProfilePhoto4x6BorderlessSemiGloss,
		ProfilePhoto5x7BorderlessGlossy,
		ProfilePhoto5x7BorderlessMatte,
		ProfilePhoto5x7BorderlessSemiGloss,
		ProfilePhotoA4BorderlessGlossy,
		ProfilePhotoA4BorderlessMatte,
		ProfilePhotoA4BorderlessSemiGloss,
		ProfilePhotoA3BorderlessGlossy,
		ProfilePhotoA3BorderlessMatte,
		ProfilePhotoA3BorderlessSemiGloss,
		ProfilePhotoA3PlusBorderlessGlossy,
		ProfilePhotoA3PlusBorderlessMatte,
		ProfilePhotoA3PlusBorderlessSemiGloss,
	}

	for _, profile := range photoProfiles {
		t.Run(string(profile), func(t *testing.T) {
			opts := MustGetPrintOptions(profile)
			if opts.Quality != 5 {
				t.Errorf("Photo profile %s should use quality 5 (best), got %d", profile, opts.Quality)
			}
		})
	}
}

func TestLargeFormatProfilesUseRearTray(t *testing.T) {
	largeProfiles := []PrintProfile{
		ProfilePhotoA3BorderlessGlossy,
		ProfilePhotoA3BorderlessMatte,
		ProfilePhotoA3BorderlessSemiGloss,
		ProfilePhotoA3PlusBorderlessGlossy,
		ProfilePhotoA3PlusBorderlessMatte,
		ProfilePhotoA3PlusBorderlessSemiGloss,
	}

	for _, profile := range largeProfiles {
		t.Run(string(profile), func(t *testing.T) {
			opts := MustGetPrintOptions(profile)
			if opts.Tray != "Rear" {
				t.Errorf("Large format profile %s should use Rear tray, got %s", profile, opts.Tray)
			}
		})
	}
}

func TestSmallPhotoProfilesUsePhotoTray(t *testing.T) {
	smallProfiles := []PrintProfile{
		ProfilePhoto4x6BorderlessGlossy,
		ProfilePhoto4x6BorderlessMatte,
		ProfilePhoto4x6BorderlessSemiGloss,
		ProfilePhoto5x7BorderlessGlossy,
		ProfilePhoto5x7BorderlessMatte,
		ProfilePhoto5x7BorderlessSemiGloss,
	}

	for _, profile := range smallProfiles {
		t.Run(string(profile), func(t *testing.T) {
			opts := MustGetPrintOptions(profile)
			if opts.Tray != "Photo" {
				t.Errorf("Small photo profile %s should use Photo tray, got %s", profile, opts.Tray)
			}
		})
	}
}

func TestGetProfileID(t *testing.T) {
	tests := []struct {
		name       string
		profile    PrintProfile
		expectedID int
	}{
		{
			name:       "Default profile has ID 0",
			profile:    ProfileDefault,
			expectedID: 0,
		},
		{
			name:       "Photo 4x6 glossy has ID 1",
			profile:    ProfilePhoto4x6BorderlessGlossy,
			expectedID: 1,
		},
		{
			name:       "Photo A3+ matte has ID 14",
			profile:    ProfilePhotoA3PlusBorderlessMatte,
			expectedID: 14,
		},
		{
			name:       "Document normal has ID 17",
			profile:    ProfileDocumentNormal,
			expectedID: 17,
		},
		{
			name:       "Document best has ID 18",
			profile:    ProfileDocumentBest,
			expectedID: 18,
		},
		{
			name:       "Unknown profile returns -1",
			profile:    "unknown-profile",
			expectedID: -1,
		},
		{
			name:       "Empty profile returns -1",
			profile:    "",
			expectedID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := GetProfileID(tt.profile)
			if id != tt.expectedID {
				t.Errorf("GetProfileID(%s) = %d, expected %d", tt.profile, id, tt.expectedID)
			}
		})
	}
}

func TestGetProfileIDRoundTrip(t *testing.T) {
	// Test that GetProfileID and GetProfileByID are inverses of each other
	for id := 0; id <= 18; id++ {
		profile, err := GetProfileByID(id)
		if err != nil {
			t.Fatalf("GetProfileByID(%d) failed: %v", id, err)
		}

		gotID := GetProfileID(profile)
		if gotID != id {
			t.Errorf("Round trip failed: GetProfileByID(%d) = %s, GetProfileID(%s) = %d, expected %d",
				id, profile, profile, gotID, id)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
