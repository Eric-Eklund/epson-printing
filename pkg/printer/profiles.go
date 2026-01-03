package printer

import "fmt"

// PrintProfile represents a named print configuration profile
type PrintProfile string

// Predefined print profiles
const (
	// ProfileDefault is the default profile for test/draft printing.
	// Paper: A4 | Tray: Main | Media: Plain paper | Quality: 3 (draft)
	ProfileDefault PrintProfile = "default"

	// --- 4x6" Borderless Photo Profiles ---

	// ProfilePhoto4x6BorderlessGlossy prints 4x6" borderless photos on glossy paper.
	// Paper: 4x6.Borderless | Tray: Photo | Media: Glossy | Quality: 5 (best)
	ProfilePhoto4x6BorderlessGlossy PrintProfile = "photo-4x6-borderless-glossy"

	// ProfilePhoto4x6BorderlessMatte prints 4x6" borderless photos on matte paper.
	// Paper: 4x6.Borderless | Tray: Photo | Media: Matte | Quality: 5 (best)
	ProfilePhoto4x6BorderlessMatte PrintProfile = "photo-4x6-borderless-matte"

	// ProfilePhoto4x6BorderlessSemiGloss prints 4x6" borderless photos on semi-gloss paper.
	// Paper: 4x6.Borderless | Tray: Photo | Media: Semi-gloss | Quality: 5 (best)
	ProfilePhoto4x6BorderlessSemiGloss PrintProfile = "photo-4x6-borderless-semigloss"

	// --- 5x7" Borderless Photo Profiles ---

	// ProfilePhoto5x7BorderlessGlossy prints 5x7" borderless photos on glossy paper.
	// Paper: 5x7.Borderless | Tray: Photo | Media: Glossy | Quality: 5 (best)
	ProfilePhoto5x7BorderlessGlossy PrintProfile = "photo-5x7-borderless-glossy"

	// ProfilePhoto5x7BorderlessMatte prints 5x7" borderless photos on matte paper.
	// Paper: 5x7.Borderless | Tray: Photo | Media: Matte | Quality: 5 (best)
	ProfilePhoto5x7BorderlessMatte PrintProfile = "photo-5x7-borderless-matte"

	// ProfilePhoto5x7BorderlessSemiGloss prints 5x7" borderless photos on semi-gloss paper.
	// Paper: 5x7.Borderless | Tray: Photo | Media: Semi-gloss | Quality: 5 (best)
	ProfilePhoto5x7BorderlessSemiGloss PrintProfile = "photo-5x7-borderless-semigloss"

	// --- A4 Borderless Photo Profiles ---

	// ProfilePhotoA4BorderlessGlossy prints A4 borderless photos on glossy paper.
	// Paper: A4.Borderless | Tray: Auto | Media: Glossy | Quality: 5 (best)
	ProfilePhotoA4BorderlessGlossy PrintProfile = "photo-a4-borderless-glossy"

	// ProfilePhotoA4BorderlessMatte prints A4 borderless photos on matte paper.
	// Paper: A4.Borderless | Tray: Auto | Media: Matte | Quality: 5 (best)
	ProfilePhotoA4BorderlessMatte PrintProfile = "photo-a4-borderless-matte"

	// ProfilePhotoA4BorderlessSemiGloss prints A4 borderless photos on semi-gloss paper.
	// Paper: A4.Borderless | Tray: Auto | Media: Semi-gloss | Quality: 5 (best)
	ProfilePhotoA4BorderlessSemiGloss PrintProfile = "photo-a4-borderless-semigloss"

	// --- A3 Borderless Photo Profiles ---

	// ProfilePhotoA3BorderlessGlossy prints A3 borderless photos on glossy paper.
	// Paper: A3.Borderless | Tray: Rear | Media: Glossy | Quality: 5 (best)
	ProfilePhotoA3BorderlessGlossy PrintProfile = "photo-a3-borderless-glossy"

	// ProfilePhotoA3BorderlessMatte prints A3 borderless photos on matte paper.
	// Paper: A3.Borderless | Tray: Rear | Media: Matte | Quality: 5 (best)
	ProfilePhotoA3BorderlessMatte PrintProfile = "photo-a3-borderless-matte"

	// ProfilePhotoA3BorderlessSemiGloss prints A3 borderless photos on semi-gloss paper.
	// Paper: A3.Borderless | Tray: Rear | Media: Semi-gloss | Quality: 5 (best)
	ProfilePhotoA3BorderlessSemiGloss PrintProfile = "photo-a3-borderless-semigloss"

	// --- A3+ (13x19") Borderless Photo Profiles ---

	// ProfilePhotoA3PlusBorderlessGlossy prints A3+ (13x19") borderless photos on glossy paper.
	// Paper: 13x19.Borderless | Tray: Rear | Media: Glossy | Quality: 5 (best)
	ProfilePhotoA3PlusBorderlessGlossy PrintProfile = "photo-a3plus-borderless-glossy"

	// ProfilePhotoA3PlusBorderlessMatte prints A3+ (13x19") borderless photos on matte paper.
	// Paper: 13x19.Borderless | Tray: Rear | Media: Matte | Quality: 5 (best)
	ProfilePhotoA3PlusBorderlessMatte PrintProfile = "photo-a3plus-borderless-matte"

	// ProfilePhotoA3PlusBorderlessSemiGloss prints A3+ (13x19") borderless photos on semi-gloss paper.
	// Paper: 13x19.Borderless | Tray: Rear | Media: Semi-gloss | Quality: 5 (best)
	ProfilePhotoA3PlusBorderlessSemiGloss PrintProfile = "photo-a3plus-borderless-semigloss"

	// --- Document Profiles ---

	// ProfileDocumentDraft prints documents in draft quality on plain paper.
	// Paper: A4 | Tray: Main | Media: Plain paper | Quality: 3 (draft)
	ProfileDocumentDraft PrintProfile = "document-draft"

	// ProfileDocumentNormal prints documents in normal quality on plain paper.
	// Paper: A4 | Tray: Main | Media: Plain paper | Quality: 4 (normal/high)
	ProfileDocumentNormal PrintProfile = "document-normal"

	// ProfileDocumentBest prints documents in best quality on coated paper.
	// Paper: A4 | Tray: Main | Media: Coated paper | Quality: 5 (best)
	ProfileDocumentBest PrintProfile = "document-best"
)

// printProfiles is the registry of all available print profiles
var printProfiles = map[PrintProfile]PrintOptions{
	// Default profile - A4 draft on plain paper
	ProfileDefault: {
		PaperSize: "A4",
		Tray:      "Main",
		MediaType: "stationery",
		Quality:   3, // Draft
		PageRange: "all",
		Copies:    1,
	},

	// 4x6 Borderless profiles
	ProfilePhoto4x6BorderlessGlossy: {
		PaperSize: "4x6.Borderless",
		Tray:      "Photo",
		MediaType: "photographic-glossy",
		Quality:   5, // Best quality for photos
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhoto4x6BorderlessMatte: {
		PaperSize: "4x6.Borderless",
		Tray:      "Photo",
		MediaType: "photographic-matte",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhoto4x6BorderlessSemiGloss: {
		PaperSize: "4x6.Borderless",
		Tray:      "Photo",
		MediaType: "photographic-semi-gloss",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},

	// 5x7 Borderless profiles
	ProfilePhoto5x7BorderlessGlossy: {
		PaperSize: "5x7.Borderless",
		Tray:      "Photo",
		MediaType: "photographic-glossy",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhoto5x7BorderlessMatte: {
		PaperSize: "5x7.Borderless",
		Tray:      "Photo",
		MediaType: "photographic-matte",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhoto5x7BorderlessSemiGloss: {
		PaperSize: "5x7.Borderless",
		Tray:      "Photo",
		MediaType: "photographic-semi-gloss",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},

	// A4 Borderless profiles
	ProfilePhotoA4BorderlessGlossy: {
		PaperSize: "A4.Borderless",
		Tray:      "Auto",
		MediaType: "photographic-glossy",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhotoA4BorderlessMatte: {
		PaperSize: "A4.Borderless",
		Tray:      "Auto",
		MediaType: "photographic-matte",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhotoA4BorderlessSemiGloss: {
		PaperSize: "A4.Borderless",
		Tray:      "Auto",
		MediaType: "photographic-semi-gloss",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},

	// A3 Borderless profiles
	ProfilePhotoA3BorderlessGlossy: {
		PaperSize: "A3.Borderless",
		Tray:      "Rear", // A3 requires rear tray
		MediaType: "photographic-glossy",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhotoA3BorderlessMatte: {
		PaperSize: "A3.Borderless",
		Tray:      "Rear",
		MediaType: "photographic-matte",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhotoA3BorderlessSemiGloss: {
		PaperSize: "A3.Borderless",
		Tray:      "Rear",
		MediaType: "photographic-semi-gloss",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},

	// A3+ Borderless profiles (13x19)
	ProfilePhotoA3PlusBorderlessGlossy: {
		PaperSize: "13x19.Borderless",
		Tray:      "Rear", // A3+ requires rear tray
		MediaType: "photographic-glossy",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhotoA3PlusBorderlessMatte: {
		PaperSize: "13x19.Borderless",
		Tray:      "Rear",
		MediaType: "photographic-matte",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},
	ProfilePhotoA3PlusBorderlessSemiGloss: {
		PaperSize: "13x19.Borderless",
		Tray:      "Rear",
		MediaType: "photographic-semi-gloss",
		Quality:   5,
		PageRange: "all",
		Copies:    1,
	},

	// Document profiles
	ProfileDocumentDraft: {
		PaperSize: "A4",
		Tray:      "Main",
		MediaType: "stationery",
		Quality:   3, // Draft quality
		PageRange: "all",
		Copies:    1,
	},
	ProfileDocumentNormal: {
		PaperSize: "A4",
		Tray:      "Main",
		MediaType: "stationery",
		Quality:   4, // Normal/high quality
		PageRange: "all",
		Copies:    1,
	},
	ProfileDocumentBest: {
		PaperSize: "A4",
		Tray:      "Main",
		MediaType: "stationery-coated",
		Quality:   5, // Best quality
		PageRange: "all",
		Copies:    1,
	},
}

// GetPrintOptions returns the PrintOptions for a given profile
// If profile is empty or "default", returns the default test/draft settings
func GetPrintOptions(profile PrintProfile) (PrintOptions, error) {
	// Handle empty profile - use default
	if profile == "" {
		profile = ProfileDefault
	}

	opts, exists := printProfiles[profile]
	if !exists {
		return PrintOptions{}, fmt.Errorf("unknown print profile: %s", profile)
	}

	return opts, nil
}

// MustGetPrintOptions returns the PrintOptions for a profile, panicking on error
// Use this only when you're certain the profile exists
func MustGetPrintOptions(profile PrintProfile) PrintOptions {
	opts, err := GetPrintOptions(profile)
	if err != nil {
		panic(err)
	}
	return opts
}

// ListProfiles returns a list of all available profile names
func ListProfiles() []PrintProfile {
	profiles := make([]PrintProfile, 0, len(printProfiles))
	for profile := range printProfiles {
		profiles = append(profiles, profile)
	}
	return profiles
}

// RegisterProfile allows registering a custom print profile
// This can be used to add user-defined profiles at runtime
func RegisterProfile(name PrintProfile, opts PrintOptions) {
	printProfiles[name] = opts
}

// GetProfileDescription returns a human-readable description of a profile
func GetProfileDescription(profile PrintProfile) string {
	opts, err := GetPrintOptions(profile)
	if err != nil {
		return "Unknown profile"
	}

	return fmt.Sprintf("%s on %s (%s, quality: %d)",
		opts.PaperSize, opts.MediaType, opts.Tray, opts.Quality)
}

// profileIDs maps numeric IDs to profile names for easier command-line usage
var profileIDs = map[int]PrintProfile{
	0: ProfileDefault,

	// 4x6" Borderless
	1: ProfilePhoto4x6BorderlessGlossy,
	2: ProfilePhoto4x6BorderlessMatte,
	3: ProfilePhoto4x6BorderlessSemiGloss,

	// 5x7" Borderless
	4: ProfilePhoto5x7BorderlessGlossy,
	5: ProfilePhoto5x7BorderlessMatte,
	6: ProfilePhoto5x7BorderlessSemiGloss,

	// A4 Borderless
	7: ProfilePhotoA4BorderlessGlossy,
	8: ProfilePhotoA4BorderlessMatte,
	9: ProfilePhotoA4BorderlessSemiGloss,

	// A3 Borderless
	10: ProfilePhotoA3BorderlessGlossy,
	11: ProfilePhotoA3BorderlessMatte,
	12: ProfilePhotoA3BorderlessSemiGloss,

	// A3+ Borderless
	13: ProfilePhotoA3PlusBorderlessGlossy,
	14: ProfilePhotoA3PlusBorderlessMatte,
	15: ProfilePhotoA3PlusBorderlessSemiGloss,

	// Documents
	16: ProfileDocumentDraft,
	17: ProfileDocumentNormal,
	18: ProfileDocumentBest,
}

// profileNameToID is the reverse mapping for looking up IDs by name
var profileNameToID map[PrintProfile]int

func init() {
	// Build reverse mapping
	profileNameToID = make(map[PrintProfile]int, len(profileIDs))
	for id, name := range profileIDs {
		profileNameToID[name] = id
	}
}

// GetProfileByID returns the profile name for a given numeric ID
func GetProfileByID(id int) (PrintProfile, error) {
	profile, exists := profileIDs[id]
	if !exists {
		return "", fmt.Errorf("unknown profile ID: %d (valid IDs: 0-18)", id)
	}
	return profile, nil
}

// GetProfileID returns the numeric ID for a given profile name
// Returns -1 if the profile name is not found
func GetProfileID(profile PrintProfile) int {
	if id, exists := profileNameToID[profile]; exists {
		return id
	}
	return -1
}

// ProfileInfo contains information about a profile for display
type ProfileInfo struct {
	ID          int
	Name        PrintProfile
	PaperSize   string
	Tray        string
	MediaType   string
	Quality     int
	Description string
}

// ListProfilesWithInfo returns detailed information about all profiles
// sorted by ID for display purposes
func ListProfilesWithInfo() []ProfileInfo {
	infos := make([]ProfileInfo, 0, len(profileIDs))

	for id := 0; id < len(profileIDs); id++ {
		profile, err := GetProfileByID(id)
		if err != nil {
			continue
		}

		opts, err := GetPrintOptions(profile)
		if err != nil {
			continue
		}

		info := ProfileInfo{
			ID:          id,
			Name:        profile,
			PaperSize:   opts.PaperSize,
			Tray:        opts.Tray,
			MediaType:   opts.MediaType,
			Quality:     opts.Quality,
			Description: GetProfileDescription(profile),
		}
		infos = append(infos, info)
	}

	return infos
}
