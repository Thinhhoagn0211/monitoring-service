-- name: GetFile :many
SELECT 
    attributes
FROM 
    file
WHERE 
    -- File name (fuzzy search with index on name column)
    (name ILIKE '%' || $1 || '%' OR $1 IS NULL)
    
    -- Exact match for file extension (index on extension column)
    AND (extension = $2 OR $2 IS NULL)
    
    -- File size range (index on size column)
    AND (size >= $3 OR $3 IS NULL)
    AND (size <= $4 OR $4 IS NULL);
    
    -- -- Creation date range (index on created_at column)
    -- AND (created_at >= COALESCE($5, '1970-01-01 00:00:00'::timestamp))
    -- AND (created_at <= COALESCE($6, CURRENT_TIMESTAMP));
    
    -- -- Modification date range (index on modified_at column)
    -- AND (modified_at >= $7 OR $7 IS NULL)
    -- AND (modified_at <= $8 OR $8 IS NULL)
    
    -- -- Last access date range (index on accessed_at column)
    -- AND (accessed_at >= $9 OR $9 IS NULL)
    -- AND (accessed_at <= $10 OR $10 IS NULL)
    
    -- -- Fuzzy search by content (assuming full-text search and index on content)
    -- AND (content ILIKE '%' || $11 || '%' OR $11 IS NULL);