package core

import (
	"fmt"
	"log"
	"math"
	"rag_system/models"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	// Adaptive thresholds
	minMeaningfulChunkSize = 200  // Minimum chars for a meaningful chunk
	maxChunkSize           = 1500 // Maximum chunk size
	preferredChunkSize     = 800  // Preferred chunk size
	overlapRatio           = 0.15 // 15% overlap

	// Document size categories
	verySmallDoc = 1000  // < 1KB - keep as single chunk or minimal splits
	smallDoc     = 3000  // < 3KB - conservative chunking
	mediumDoc    = 10000 // < 10KB - normal chunking
	largeDoc     = 50000 // < 50KB - aggressive chunking

	// Minimum chunks before splitting
	minChunksThreshold = 3
)

// DocumentCharacteristics analyzes document properties
type DocumentCharacteristics struct {
	Length        int
	Category      DocumentCategory
	HasStructure  bool
	StructureType DocumentStructureType
	Language      string
	Complexity    float64
}

type DocumentCategory string
type DocumentStructureType string

const (
	VerySmallDocument DocumentCategory = "very_small"
	SmallDocument     DocumentCategory = "small"
	MediumDocument    DocumentCategory = "medium"
	LargeDocument     DocumentCategory = "large"
	VeryLargeDocument DocumentCategory = "very_large"

	NoStructure           DocumentStructureType = "none"
	SimpleStructure       DocumentStructureType = "simple"
	SectionedStructure    DocumentStructureType = "sectioned"
	HierarchicalStructure DocumentStructureType = "hierarchical"
)

// DocumentProcessor handles advanced document processing and chunking
type DocumentProcessor struct{}

// NewDocumentProcessor creates a new document processor
// a constructor function pattern in Go.
func NewDocumentProcessor() *DocumentProcessor {
	return &DocumentProcessor{}
}

func analyzeStructure(content string) (DocumentStructureType, bool) {
	// Check for hierarchical patterns (multiple heading levels)
	hierarchicalPatterns := []string{
		`(?m)^#+\s+`,            // Markdown headers
		`(?m)^[A-Z][A-Z\s]+:?$`, // ALL CAPS sections
		`(?m)^\d+\.\s+[A-Z]`,    // Numbered sections
		`(?m)^[IVX]+\.\s+`,      // Roman numerals
	}

	structureCount := 0
	for _, pattern := range hierarchicalPatterns {
		matched, _ := regexp.MatchString(pattern, content)

		if matched {
			structureCount++
		}
	}

	// Count section-like patterns
	sectionPatterns := []string{
		`(?i)\b(experience|education|skills|summary|objective|projects|achievements|awards|certifications|languages|references|contact|about)\b`,
		`(?m)^[A-Z][A-Z\s]{3,}:?\s*$`,
		`(?m)^.{1,50}:$`,
	}
	sectionCount := 0

	for _, pattern := range sectionPatterns {
		re := regexp.MustCompile(pattern)
		sectionCount += len(re.FindAllString(content, -1))
	}

	if structureCount >= 3 || sectionCount >= 5 {
		return HierarchicalStructure, true
	} else if structureCount >= 1 || sectionCount >= 2 {
		return SectionedStructure, true
	} else if strings.Count(content, "\n\n") >= 3 {
		return SimpleStructure, true
	}

	return NoStructure, false
}

func calculateComplexity(content string) float64 {
	words := strings.Fields(content)

	if len(words) == 0 {
		return 0.0
	}
	sentences := strings.Split(content, ".")

	avgWordsPerSentence := float64(len(words)) / float64(len(sentences))

	// Simple complexity score based on sentence length
	complexity := math.Min(avgWordsPerSentence/15.0, 1.0)
	return complexity
}

// calculateOptimalChunkCount determines the ideal number of chunks for a document
func calculateOptimalChunkCount(length int) int {
	switch {
	case length < 600:
		return 1 // Single chunk
	case length < 1200:
		return 2 // Two chunks
	case length < 2000:
		return 3 // Three chunks
	case length < 4000:
		return 4 // Four chunks
	case length < 8000:
		return int(math.Ceil(float64(length) / 1500)) // ~1500 chars per chunk
	default:
		return int(math.Ceil(float64(length) / 1000)) // ~1000 chars per chunk for larger docs
	}
}

func analyzeDocument(content string) DocumentCharacteristics {
	length := len(content)
	var category DocumentCategory

	switch {
	case length < verySmallDoc:
		category = VerySmallDocument
	case length < smallDoc:
		category = SmallDocument
	case length < mediumDoc:
		category = MediumDocument
	case length < largeDoc:
		category = LargeDocument
	default:
		category = VeryLargeDocument
	}

	structureType, hasStructure := analyzeStructure(content)

	complexity := calculateComplexity(content)

	return DocumentCharacteristics{
		Length:        length,
		Category:      category,
		HasStructure:  hasStructure,
		StructureType: structureType,
		Language:      "en",
		Complexity:    complexity,
	}
}

func adaptiveChunkingStrategy(characteristics DocumentCharacteristics, config *models.ChunkingConfig) *models.ChunkingConfig {
	if config == nil {
		config = &models.ChunkingConfig{}
	}

	// Copy existing config
	adaptiveConfig := *config

	optimalChunkCount := calculateOptimalChunkCount(characteristics.Length)

	log.Printf("Document length: %d chars, optimal chunk count: %d", characteristics.Length, optimalChunkCount)

	// Override strategy based on document characteristics with smart thresholds
	switch characteristics.Category {
	case VerySmallDocument:
		// For very small documents (< 1000 chars), be very conservative
		if characteristics.Length < 600 {
			// Single chunk for very small documents
			adaptiveConfig.Strategy = models.FixedSizeStrategy
			adaptiveConfig.FixedSize = characteristics.Length
			adaptiveConfig.Overlap = 0
			adaptiveConfig.MinChunkSize = characteristics.Length
			log.Printf("Very small document: keeping as single chunk")
		} else {
			// Max 2-3 chunks for small documents
			adaptiveConfig.Strategy = models.StructuralStrategy
			adaptiveConfig.MinChunkSize = characteristics.Length / 3 // Ensure max 3 chunks
			adaptiveConfig.MaxChunkSize = characteristics.Length / 2 // Ensure min 2 chunks
			if adaptiveConfig.MinChunkSize < 250 {
				adaptiveConfig.MinChunkSize = 250
			}
			log.Printf("Small document: conservative chunking with min=%d, max=%d",
				adaptiveConfig.MinChunkSize, adaptiveConfig.MaxChunkSize)
		}

	case SmallDocument:
		// For small documents (1-3KB), aim for 3-5 meaningful chunks
		targetChunkSize := characteristics.Length / optimalChunkCount
		if targetChunkSize < 400 {
			targetChunkSize = 400
		}

		if characteristics.HasStructure {
			adaptiveConfig.Strategy = models.StructuralStrategy
			adaptiveConfig.MinChunkSize = targetChunkSize
			adaptiveConfig.MaxChunkSize = targetChunkSize + 300
		} else {
			adaptiveConfig.Strategy = models.SentenceWindowStrategy
			adaptiveConfig.SentenceWindowSize = 4
			adaptiveConfig.MinChunkSize = targetChunkSize
		}
		log.Printf("Small document: targeting %d chunks with size ~%d", optimalChunkCount, targetChunkSize)

	case MediumDocument:
		// For medium documents, use normal strategies
		if characteristics.StructureType == HierarchicalStructure {
			adaptiveConfig.Strategy = models.ParentDocumentStrategy
		} else if characteristics.HasStructure {
			adaptiveConfig.Strategy = models.StructuralStrategy
		} else {
			adaptiveConfig.Strategy = models.SemanticStrategy
		}

	case LargeDocument, VeryLargeDocument:
		// For large documents, use aggressive chunking
		adaptiveConfig.Strategy = models.ParentDocumentStrategy
		adaptiveConfig.MaxChunkSize = 1200
		adaptiveConfig.MinChunkSize = 400
	}

	// Set intelligent defaults based on document size
	if adaptiveConfig.MinChunkSize == 0 {
		if characteristics.Length < 2000 {
			adaptiveConfig.MinChunkSize = characteristics.Length / 4 // Max 4 chunks for small docs
			if adaptiveConfig.MinChunkSize < 200 {
				adaptiveConfig.MinChunkSize = 200
			}
		} else {
			adaptiveConfig.MinChunkSize = minMeaningfulChunkSize
		}
	}

	if adaptiveConfig.MaxChunkSize == 0 {
		if characteristics.Length < 3000 {
			adaptiveConfig.MaxChunkSize = characteristics.Length / 2 // Min 2 chunks for small docs
		} else {
			adaptiveConfig.MaxChunkSize = maxChunkSize
		}
	}

	if adaptiveConfig.FixedSize == 0 {
		if characteristics.Length < 2000 {
			adaptiveConfig.FixedSize = characteristics.Length / optimalChunkCount
		} else {
			adaptiveConfig.FixedSize = preferredChunkSize
		}
	}

	if adaptiveConfig.Overlap == 0 {
		// Reduce overlap for very small documents
		if characteristics.Length < 1500 {
			adaptiveConfig.Overlap = int(float64(adaptiveConfig.FixedSize) * 0.1) // 10% overlap
		} else {
			adaptiveConfig.Overlap = int(float64(adaptiveConfig.FixedSize) * overlapRatio)
		}
	}

	adaptiveConfig.PreserveParagraphs = true
	adaptiveConfig.ExtractKeywords = true

	return &adaptiveConfig
}

func ProcessDocumentContent(content string, source string, docType string, config *models.ChunkingConfig) (*models.Document, error) {
	if content == "" {
		return nil, fmt.Errorf("Content cannot be empty")
	}
	// Analyze document characteristics
	characteristics := analyzeDocument(content)

	// Override config with adaptive strategy if needed
	adaptiveConfig := adaptiveChunkingStrategy(characteristics, config)

	log.Printf("Document analysis: %d chars, category: %s, structure: %s, strategy: %s",
		characteristics.Length, characteristics.Category, characteristics.StructureType, adaptiveConfig.Strategy)

	doc := &models.Document{
		ID:      uuid.New().String(),
		Content: content,
		Source:  source,
		DocType: docType,
		Metadata: map[string]interface{}{
			"chunking_strategy": string(adaptiveConfig.Strategy),
			"document_length":   characteristics.Length,
			"document_category": string(characteristics.Category),
			"structure_type":    string(characteristics.StructureType),
			"chunk_count":       0, // Will be updated after chunking
		},
	}

	var chunks []*models.EnhancedChunk
	var err error

}
