package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rag_system/models"

	"github.com/qdrant/go-client/qdrant"
)

// VectorDB wraps the Qdrant client
type VectorDB struct {
	client *qdrant.Client
	ctx    context.Context
}

// NewVectorDB creates a new Qdrant-backed VectorDB.
// Reads QDRANT_HOST and QDRANT_API_KEY from environment variables.
// The dbPath argument is ignored (kept for signature compatibility).
func NewVectorDB(dbPath string) (*VectorDB, error) {
	host := os.Getenv("QDRANT_HOST")
	apiKey := os.Getenv("QDRANT_API_KEY")

	if host == "" {
		return nil, fmt.Errorf("QDRANT_HOST environment variable is not set")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("QDRANT_API_KEY environment variable is not set")
	}

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   6334,
		APIKey: apiKey,
		UseTLS: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Qdrant: %w", err)
	}

	ctx := context.Background()

	info, err := client.HealthCheck(ctx)
	if err != nil {
		return nil, fmt.Errorf("sqlite-vec not available: %w", err)
	}

	log.Printf("Connected to Qdrant version: %s (host: %s)", info.GetVersion(), host)

	db := &VectorDB{client: client, ctx: ctx}

	// Ensure payload indexes exist on all existing collections
	if cols, err := client.ListCollections(ctx); err == nil {
		for _, col := range cols {
			db.createPayloadIndexes(col)
			log.Printf("Ensured payload indexes for collection: %s", col)
		}
	}

	return db, nil
}

// CreateCollection creates a Qdrant collection with default dimension 1024.
func (db *VectorDB) CreateCollection(name, description string) error {
	exists, err := db.collectionExists(name)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	err = db.client.CreateCollection(db.ctx, &qdrant.CreateCollection{
		CollectionName: name,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     1024,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	// Qdrant has no collection-level metadata, store description as a special meta point
	zeroVec := make([]float32, 1024)
	_, err = db.client.Upsert(db.ctx, &qdrant.UpsertPoints{
		CollectionName: name,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(0),
				Vectors: qdrant.NewVectors(zeroVec...),
				Payload: qdrant.NewValueMap(map[string]interface{}{
					"_meta":       true,
					"name":        name,
					"description": description,
				}),
			},
		},
	})
	if err != nil {
		log.Printf("Warning: could not store collection metadata: %v", err)
	}

	db.createPayloadIndexes(name)
	return nil
}

// AddDocument stores all chunks of a document into the collection with zero vectors.
// Real vectors are added later via AddEmbeddings.
func (db *VectorDB) AddDocument(collectionName string, doc *models.Document) error {
	exists, err := db.collectionExists(collectionName)
	if err != nil {
		return err
	}
	if !exists {
		if err := db.CreateCollection(collectionName, ""); err != nil {
			return err
		}
	}

	// Set collection name in chunk metadata so AddEmbeddings stores to correct collection
	for _, chunk := range doc.Chunks {
		if chunk.Metadata == nil {
			chunk.Metadata = make(map[string]interface{})
		}
		chunk.Metadata["collection_name"] = collectionName
	}

	var points []*qdrant.PointStruct
	for i, chunk := range doc.Chunks {
		payload := db.chunkToPayload(chunk, doc)
		zeroVec := make([]float32, 1024)

		points = append(points, &qdrant.PointStruct{
			Id:      qdrant.NewIDUUID(chunk.ID),
			Vectors: qdrant.NewVectors(zeroVec...),
			Payload: qdrant.NewValueMap(payload),
		})

		if len(points) == 100 || i == len(doc.Chunks)-1 {
			_, err := db.client.Upsert(db.ctx, &qdrant.UpsertPoints{
				CollectionName: collectionName,
				Points:         points,
			})
			if err != nil {
				return fmt.Errorf("failed to insert chunks: %w", err)
			}
			points = nil
		}
	}

	log.Printf("Stored document %s with %d chunks into collection %s", doc.ID, len(doc.Chunks), collectionName)
	return nil
}

// AddEmbeddings upserts chunks with their real embedding vectors into Qdrant.
func (db *VectorDB) AddEmbeddings(chunks []*models.EnhancedChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	// Determine embedding dimension from first non-empty chunk
	var embeddingDim uint64
	for _, chunk := range chunks {
		if len(chunk.Embedding) > 0 {
			embeddingDim = uint64(len(chunk.Embedding))
			break
		}
	}
	if embeddingDim == 0 {
		return fmt.Errorf("no valid embeddings found in chunks")
	}

	// Group chunks by collection name stored in metadata, fallback to "default"
	byCollection := map[string][]*models.EnhancedChunk{}
	for _, chunk := range chunks {
		col := "default"
		if chunk.Metadata != nil {
			if c, ok := chunk.Metadata["collection_name"].(string); ok && c != "" {
				col = c
			}
		}
		byCollection[col] = append(byCollection[col], chunk)
	}

	for collectionName, colChunks := range byCollection {
		if err := db.ensureCollectionDimension(collectionName, embeddingDim); err != nil {
			return err
		}

		var points []*qdrant.PointStruct
		for i, chunk := range colChunks {
			if len(chunk.Embedding) == 0 {
				continue
			}

			payload := db.chunkToPayload(chunk, nil)
			payload["collection_name"] = collectionName

			points = append(points, &qdrant.PointStruct{
				Id:      qdrant.NewIDUUID(chunk.ID),
				Vectors: qdrant.NewVectors(chunk.Embedding...),
				Payload: qdrant.NewValueMap(payload),
			})

			if len(points) == 100 || i == len(colChunks)-1 {
				_, err := db.client.Upsert(db.ctx, &qdrant.UpsertPoints{
					CollectionName: collectionName,
					Points:         points,
				})
				if err != nil {
					return fmt.Errorf("failed to upsert embeddings: %w", err)
				}
				points = nil
			}
		}
		log.Printf("Upserted %d embeddings into collection %s", len(colChunks), collectionName)
	}

	return nil
}

// QuerySimilarChunks performs a vector similarity search in Qdrant.
func (db *VectorDB) QuerySimilarChunks(collectionName string, queryEmbedding []float32, topK int, filters map[string]interface{}) ([]*models.EnhancedChunk, []float64, error) {
	mustNot := []*qdrant.Condition{
		qdrant.NewMatch("chunk_type", "meta"),
	}

	var must []*qdrant.Condition
	for key, value := range filters {
		switch key {
		case "chunk_type":
			must = append(must, qdrant.NewMatch("chunk_type", fmt.Sprintf("%v", value)))
		case "section":
			must = append(must, qdrant.NewMatch("section", fmt.Sprintf("%v", value)))
		case "doc_type":
			must = append(must, qdrant.NewMatch("doc_type", fmt.Sprintf("%v", value)))
		}
	}

	limit := uint64(topK)
	results, err := db.client.Query(db.ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(queryEmbedding...),
		Filter: &qdrant.Filter{
			Must:    must,
			MustNot: mustNot,
		},
		Limit:       &limit,
		WithPayload: qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query similar chunks: %w", err)
	}

	var chunks []*models.EnhancedChunk
	var scores []float64
	for _, result := range results {
		chunk := db.payloadToChunk(result.GetPayload())
		chunks = append(chunks, chunk)
		scores = append(scores, float64(result.Score))
	}

	return chunks, scores, nil
}

// GetChunkWithParents retrieves a chunk and walks up its parent hierarchy.
func (db *VectorDB) GetChunkWithParents(chunkID string) ([]*models.EnhancedChunk, error) {
	collections, err := db.client.ListCollections(db.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	// Find the chunk across all collections
	var firstChunk *models.EnhancedChunk
	var foundCollection string
	for _, col := range collections {
		pts, err := db.client.Get(db.ctx, &qdrant.GetPoints{
			CollectionName: col,
			Ids:            []*qdrant.PointId{qdrant.NewIDUUID(chunkID)},
			WithPayload:    qdrant.NewWithPayload(true),
		})
		if err == nil && len(pts) > 0 {
			firstChunk = db.payloadToChunk(pts[0].GetPayload())
			foundCollection = col
			break
		}
	}

	if firstChunk == nil {
		return nil, fmt.Errorf("chunk %s not found", chunkID)
	}

	// Walk up the parent hierarchy, prepending so order is root → leaf
	var hierarchy []*models.EnhancedChunk
	current := firstChunk
	for current != nil {
		hierarchy = append([]*models.EnhancedChunk{current}, hierarchy...)
		if current.ParentChunkID == nil || *current.ParentChunkID == "" {
			break
		}
		pts, err := db.client.Get(db.ctx, &qdrant.GetPoints{
			CollectionName: foundCollection,
			Ids:            []*qdrant.PointId{qdrant.NewIDUUID(*current.ParentChunkID)},
			WithPayload:    qdrant.NewWithPayload(true),
		})
		if err != nil || len(pts) == 0 {
			break
		}
		current = db.payloadToChunk(pts[0].GetPayload())
	}

	return hierarchy, nil
}

// AddChunk is legacy support.
func (db *VectorDB) AddChunk(collectionName string, chunk *models.DocumentChunk) error {
	enhanced := &models.EnhancedChunk{
		ID:         chunk.ID,
		DocumentID: chunk.DocumentID,
		Text:       chunk.Text,
		Embedding:  chunk.Embedding,
		ChunkType:  "legacy",
		Metadata:   map[string]interface{}{"collection_name": collectionName},
	}
	return db.AddEmbeddings([]*models.EnhancedChunk{enhanced})
}

// QuerySimilar is legacy support.
func (db *VectorDB) QuerySimilar(collectionName string, queryEmbedding []float32, topK int) ([]*models.DocumentChunk, error) {
	enhancedChunks, _, err := db.QuerySimilarChunks(collectionName, queryEmbedding, topK, nil)
	if err != nil {
		return nil, err
	}

	var legacyChunks []*models.DocumentChunk
	for _, chunk := range enhancedChunks {
		legacyChunks = append(legacyChunks, &models.DocumentChunk{
			ID:         chunk.ID,
			DocumentID: chunk.DocumentID,
			Text:       chunk.Text,
			Embedding:  chunk.Embedding,
		})
	}
	return legacyChunks, nil
}

// Close closes the Qdrant client connection.
func (db *VectorDB) Close() error {
	return db.client.Close()
}

// ListCollections returns all Qdrant collections with basic stats.
func (db *VectorDB) ListCollections() ([]map[string]interface{}, error) {
	collections, err := db.client.ListCollections(db.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	var result []map[string]interface{}
	for _, name := range collections {
		info, err := db.client.GetCollectionInfo(db.ctx, name)
		var pointCount uint64
		if err == nil {
			pointCount = info.GetPointsCount()
			if pointCount > 0 {
				pointCount-- // subtract the meta point
			}
		}
		result = append(result, map[string]interface{}{
			"name":        name,
			"description": "",
			"created_at":  "",
			"doc_count":   0,
			"chunk_count": pointCount,
		})
	}

	return result, nil
}

// DeleteCollection deletes a Qdrant collection entirely.
func (db *VectorDB) DeleteCollection(name string) error {
	exists, err := db.collectionExists(name)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("collection '%s' not found", name)
	}
	if err := db.client.DeleteCollection(db.ctx, name); err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}
	return nil
}

// ListDocuments returns all unique documents in a collection derived from chunk payloads.
func (db *VectorDB) ListDocuments(collectionName string) ([]map[string]interface{}, error) {
	limit := uint32(1000)
	results, err := db.client.Scroll(db.ctx, &qdrant.ScrollPoints{
		CollectionName: collectionName,
		Filter: &qdrant.Filter{
			MustNot: []*qdrant.Condition{
				qdrant.NewMatch("chunk_type", "meta"),
			},
		},
		Limit:       &limit,
		WithPayload: qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	docMap := map[string]map[string]interface{}{}
	for _, point := range results {
		payload := point.GetPayload()
		docID := payloadString(payload, "document_id")
		if docID == "" {
			continue
		}
		if _, exists := docMap[docID]; !exists {
			docMap[docID] = map[string]interface{}{
				"id":          docID,
				"source":      payloadString(payload, "source"),
				"doc_type":    payloadString(payload, "doc_type"),
				"created_at":  "",
				"chunk_count": 0,
			}
		}
		docMap[docID]["chunk_count"] = docMap[docID]["chunk_count"].(int) + 1
	}

	var documents []map[string]interface{}
	for _, doc := range docMap {
		documents = append(documents, doc)
	}
	return documents, nil
}

// DeleteDocument deletes all chunks belonging to a document ID across all collections.
func (db *VectorDB) DeleteDocument(documentID string) error {
	collections, err := db.client.ListCollections(db.ctx)
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	deleted := false
	for _, collectionName := range collections {
		_, err := db.client.Delete(db.ctx, &qdrant.DeletePoints{
			CollectionName: collectionName,
			Points: qdrant.NewPointsSelectorFilter(&qdrant.Filter{
				Must: []*qdrant.Condition{
					qdrant.NewMatch("document_id", documentID),
				},
			}),
		})
		if err == nil {
			deleted = true
			log.Printf("Deleted chunks for document '%s' from collection '%s'", documentID, collectionName)
		}
	}

	if !deleted {
		return fmt.Errorf("document with ID '%s' not found", documentID)
	}
	return nil
}

// DeleteAllDocumentsInCollection deletes all non-meta points from a collection.
func (db *VectorDB) DeleteAllDocumentsInCollection(collectionName string) error {
	limit := uint32(1)
	results, err := db.client.Scroll(db.ctx, &qdrant.ScrollPoints{
		CollectionName: collectionName,
		Filter: &qdrant.Filter{
			MustNot: []*qdrant.Condition{
				qdrant.NewMatch("chunk_type", "meta"),
			},
		},
		Limit:       &limit,
		WithPayload: qdrant.NewWithPayload(false),
	})
	if err != nil || len(results) == 0 {
		return fmt.Errorf("no documents found in collection '%s'", collectionName)
	}

	_, err = db.client.Delete(db.ctx, &qdrant.DeletePoints{
		CollectionName: collectionName,
		Points: qdrant.NewPointsSelectorFilter(&qdrant.Filter{
			MustNot: []*qdrant.Condition{
				qdrant.NewMatch("chunk_type", "meta"),
			},
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}

	log.Printf("Deleted all documents from collection '%s'", collectionName)
	return nil
}

// GetCollectionStats returns stats for a named collection.
func (db *VectorDB) GetCollectionStats(collectionName string) (map[string]interface{}, error) {
	exists, err := db.collectionExists(collectionName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("collection '%s' not found", collectionName)
	}

	info, err := db.client.GetCollectionInfo(db.ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection info: %w", err)
	}

	pointCount := info.GetPointsCount()
	if pointCount > 0 {
		pointCount--
	}

	stats := map[string]interface{}{
		"name":           collectionName,
		"description":    "",
		"created_at":     "",
		"document_count": 0,
		"chunk_count":    pointCount,
		"chunk_types":    map[string]int{},
		"document_types": map[string]int{},
	}

	limit := uint32(1000)
	results, err := db.client.Scroll(db.ctx, &qdrant.ScrollPoints{
		CollectionName: collectionName,
		Filter: &qdrant.Filter{
			MustNot: []*qdrant.Condition{
				qdrant.NewMatch("chunk_type", "meta"),
			},
		},
		Limit:       &limit,
		WithPayload: qdrant.NewWithPayload(true),
	})
	if err == nil {
		chunkTypes := map[string]int{}
		docTypes := map[string]int{}
		docIDs := map[string]bool{}

		for _, point := range results {
			payload := point.GetPayload()
			if ct := payloadString(payload, "chunk_type"); ct != "" {
				chunkTypes[ct]++
			}
			if dt := payloadString(payload, "doc_type"); dt != "" {
				docTypes[dt]++
			}
			if did := payloadString(payload, "document_id"); did != "" {
				docIDs[did] = true
			}
		}

		stats["chunk_types"] = chunkTypes
		stats["document_types"] = docTypes
		stats["document_count"] = len(docIDs)
	}

	return stats, nil
}

// ── internal helpers ──────────────────────────────────────────────────────────

func (db *VectorDB) collectionExists(name string) (bool, error) {
	collections, err := db.client.ListCollections(db.ctx)
	if err != nil {
		return false, fmt.Errorf("failed to list collections: %w", err)
	}
	for _, c := range collections {
		if c == name {
			return true, nil
		}
	}
	return false, nil
}

func (db *VectorDB) ensureCollectionDimension(collectionName string, dimension uint64) error {
	exists, err := db.collectionExists(collectionName)
	if err != nil {
		return err
	}

	if exists {
		info, err := db.client.GetCollectionInfo(db.ctx, collectionName)
		if err == nil {
			if cfg := info.GetConfig(); cfg != nil {
				if params := cfg.GetParams(); params != nil {
					if vecs := params.GetVectorsConfig(); vecs != nil {
						if vp := vecs.GetParams(); vp != nil && vp.Size == dimension {
							return nil // already correct, nothing to do
						}
					}
				}
			}
		}
		// Dimension mismatch or couldn't verify — recreate
		log.Printf("Recreating collection %s for dimension %d", collectionName, dimension)
		if err := db.client.DeleteCollection(db.ctx, collectionName); err != nil {
			return fmt.Errorf("failed to delete collection for recreation: %w", err)
		}
	}

	// Create with correct dimension
	if err := db.client.CreateCollection(db.ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     dimension,
			Distance: qdrant.Distance_Cosine,
		}),
	}); err != nil {
		return fmt.Errorf("failed to create collection %s: %w", collectionName, err)
	}

	// Always insert the meta point after creation so filters work
	zeroVec := make([]float32, dimension)
	_, err = db.client.Upsert(db.ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(0),
				Vectors: qdrant.NewVectors(zeroVec...),
				Payload: qdrant.NewValueMap(map[string]interface{}{
					"chunk_type": "meta",
					"name":       collectionName,
				}),
			},
		},
	})
	if err != nil {
		log.Printf("Warning: could not store meta point for collection %s: %v", collectionName, err)
	}
	db.createPayloadIndexes(collectionName)
	return nil
}

// chunkToPayload converts an EnhancedChunk to a flat Qdrant payload map.
// ParentChunkID is *string in the model, so we dereference it safely.
func (db *VectorDB) chunkToPayload(chunk *models.EnhancedChunk, doc *models.Document) map[string]interface{} {
	childIDsJSON, _ := json.Marshal(chunk.ChildChunkIDs)
	keywordsJSON, _ := json.Marshal(chunk.Keywords)
	metadataJSON, _ := json.Marshal(chunk.Metadata)

	parentID := ""
	if chunk.ParentChunkID != nil {
		parentID = *chunk.ParentChunkID
	}

	payload := map[string]interface{}{
		"chunk_id":        chunk.ID,
		"document_id":     chunk.DocumentID,
		"text":            chunk.Text,
		"parent_chunk_id": parentID,
		"child_chunk_ids": string(childIDsJSON),
		"section":         chunk.Section,
		"subsection":      chunk.Subsection,
		"chunk_type":      chunk.ChunkType,
		"start_pos":       chunk.StartPos,
		"end_pos":         chunk.EndPos,
		"chunk_index":     chunk.ChunkIndex,
		"keywords":        string(keywordsJSON),
		"metadata":        string(metadataJSON),
		"confidence":      chunk.Confidence,
	}

	if doc != nil {
		payload["source"] = doc.Source
		payload["doc_type"] = doc.DocType
	}

	return payload
}

// payloadToChunk converts a Qdrant payload back to an EnhancedChunk matching models.go.
func (db *VectorDB) payloadToChunk(payload map[string]*qdrant.Value) *models.EnhancedChunk {
	chunk := &models.EnhancedChunk{
		ID:         payloadString(payload, "chunk_id"),
		DocumentID: payloadString(payload, "document_id"),
		Text:       payloadString(payload, "text"),
		Section:    payloadString(payload, "section"),
		Subsection: payloadString(payload, "subsection"),
		ChunkType:  payloadString(payload, "chunk_type"),
		Confidence: payloadFloat(payload, "confidence"),
		StartPos:   payloadInt(payload, "start_pos"),
		EndPos:     payloadInt(payload, "end_pos"),
		ChunkIndex: payloadInt(payload, "chunk_index"),
	}

	// ParentChunkID is *string in the model
	if parentID := payloadString(payload, "parent_chunk_id"); parentID != "" {
		chunk.ParentChunkID = &parentID
	}

	if v := payloadString(payload, "child_chunk_ids"); v != "" && v != "null" {
		json.Unmarshal([]byte(v), &chunk.ChildChunkIDs)
	}
	if v := payloadString(payload, "keywords"); v != "" && v != "null" {
		json.Unmarshal([]byte(v), &chunk.Keywords)
	}
	if v := payloadString(payload, "metadata"); v != "" && v != "null" {
		json.Unmarshal([]byte(v), &chunk.Metadata)
	}

	return chunk
}

func payloadString(payload map[string]*qdrant.Value, key string) string {
	if v, ok := payload[key]; ok {
		return v.GetStringValue()
	}
	return ""
}

func payloadFloat(payload map[string]*qdrant.Value, key string) float64 {
	if v, ok := payload[key]; ok {
		return v.GetDoubleValue()
	}
	return 0
}

func payloadInt(payload map[string]*qdrant.Value, key string) int {
	if v, ok := payload[key]; ok {
		return int(v.GetIntegerValue())
	}
	return 0
}

func (db *VectorDB) createPayloadIndexes(collectionName string) {
	fieldIndexes := []string{"chunk_type", "section", "doc_type", "document_id"}
	for _, field := range fieldIndexes {
		_, err := db.client.CreateFieldIndex(db.ctx, &qdrant.CreateFieldIndexCollection{
			CollectionName: collectionName,
			FieldName:      field,
			FieldType:      qdrant.FieldType_FieldTypeKeyword.Enum(),
		})
		if err != nil {
			log.Printf("Warning: could not create index for field %s in %s: %v", field, collectionName, err)
		}
	}
}
