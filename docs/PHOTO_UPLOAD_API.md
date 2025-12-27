
# Photo Upload API Documentation

## Overview

The Gymie backend provides endpoints for uploading, managing, and serving progress photos. Photos are stored locally on the server and served via static file endpoint.

## Configuration

### Environment Variables

Add these to your `.env` file:

```env
# Storage Configuration
STORAGE_BASE_PATH=./uploads              # Local file storage directory
STORAGE_PUBLIC_URL=http://localhost:8080/uploads  # Public URL for accessing files
```

### Directory Structure

Photos are organized by user ID:
```
uploads/
└── progress_photos/
    ├── 1/
    │   ├── uuid-1.jpg
    │   └── uuid-2.jpg
    └── 2/
        └── uuid-3.jpg
```

## API Endpoints

### 1. Upload Progress Photo

Upload a new progress photo with image file.

**Endpoint:** `POST /v1/progress/photos`

**Authentication:** Required (Bearer Token)

**Content-Type:** `multipart/form-data`

**Form Fields:**
- `photo` (file, required): Image file (JPEG, PNG, or WebP)
- `date` (string, required): Date in ISO 8601 format (e.g., `2024-12-26T09:00:00Z`)
- `notes` (string, optional): Optional notes about the photo
- `weight` (number, optional): Weight in kg at the time of photo

**Constraints:**
- Maximum file size: 10MB
- Allowed formats: JPEG, PNG, WebP
- File is validated for type and size

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "Progress photo uploaded successfully",
  "data": {
    "id": 1,
    "userId": 123,
    "date": "2024-12-26T09:00:00Z",
    "imageUrl": "http://localhost:8080/uploads/progress_photos/123/uuid.jpg",
    "weight": 75.5,
    "notes": "Post-workout progress",
    "createdAt": "2024-12-26T09:00:00Z",
    "updatedAt": "2024-12-26T09:00:00Z"
  }
}
```

**Error Responses:**

400 Bad Request - Invalid file:
```json
{
  "success": false,
  "error": {
    "code": "validation_error",
    "message": "Invalid file type. Only JPEG, PNG, and WebP images are allowed"
  }
}
```

400 Bad Request - File too large:
```json
{
  "success": false,
  "error": {
    "code": "validation_error",
    "message": "File size exceeds maximum limit of 10MB"
  }
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/v1/progress/photos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "photo=@/path/to/photo.jpg" \
  -F "date=2024-12-26T09:00:00Z" \
  -F "notes=Post-workout progress" \
  -F "weight=75.5"
```

**JavaScript/React Native Example:**
```javascript
const uploadPhoto = async (photoUri, date, notes, weight) => {
  const formData = new FormData();
  formData.append('photo', {
    uri: photoUri,
    type: 'image/jpeg',
    name: 'progress-photo.jpg',
  });
  formData.append('date', date);
  if (notes) formData.append('notes', notes);
  if (weight) formData.append('weight', weight.toString());

  const response = await fetch('http://localhost:8080/v1/progress/photos', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
    },
    body: formData,
  });

  return await response.json();
};
```

### 2. List Progress Photos

Get all progress photos for the authenticated user.

**Endpoint:** `GET /v1/progress/photos`

**Authentication:** Required

**Query Parameters:**
- `page` (integer, optional): Page number (default: 1)
- `page_size` (integer, optional): Items per page (default: 20)

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "userId": 123,
      "date": "2024-12-26T09:00:00Z",
      "imageUrl": "http://localhost:8080/uploads/progress_photos/123/uuid.jpg",
      "weight": 75.5,
      "notes": "Post-workout progress",
      "createdAt": "2024-12-26T09:00:00Z",
      "updatedAt": "2024-12-26T09:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

### 3. Update Progress Photo

Update photo metadata (not the image itself).

**Endpoint:** `PUT /v1/progress/photos/:id`

**Authentication:** Required

**Request Body:**
```json
{
  "date": "2024-12-26T10:00:00Z",
  "weight": 76.0,
  "notes": "Updated notes"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Progress photo updated successfully",
  "data": {
    "id": 1,
    "userId": 123,
    "date": "2024-12-26T10:00:00Z",
    "imageUrl": "http://localhost:8080/uploads/progress_photos/123/uuid.jpg",
    "weight": 76.0,
    "notes": "Updated notes",
    "createdAt": "2024-12-26T09:00:00Z",
    "updatedAt": "2024-12-26T10:00:00Z"
  }
}
```

### 4. Delete Progress Photo

Delete a progress photo (removes both database record and file).

**Endpoint:** `DELETE /v1/progress/photos/:id`

**Authentication:** Required

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Progress photo deleted successfully"
}
```

### 5. Access Photo File

Photos are served as static files and can be accessed directly via their URL.

**Endpoint:** `GET /uploads/progress_photos/:userId/:filename`

**Authentication:** Not required (public access)

**Example:**
```
http://localhost:8080/uploads/progress_photos/123/a1b2c3d4-e5f6-7890-abcd-ef1234567890.jpg
```

## Implementation Details

### File Storage

- Photos are stored in the local filesystem
- Each user has a dedicated directory
- Filenames are UUIDs to prevent conflicts
- Original file extension is preserved

### File Validation

- Content-Type is validated
- File size is limited to 10MB
- Only image formats are allowed (JPEG, PNG, WebP)

### Security

- User can only access/modify their own photos
- Authentication required for all write operations
- User ID is embedded in file path for isolation

### Error Handling

- If database insert fails, uploaded file is deleted
- If file upload fails, no database record is created
- Graceful handling of missing files during deletion

## Integration with Frontend

Update your frontend `photoSyncService.ts`:

```typescript
const BACKEND_URL = 'http://localhost:8080/v1';

async uploadPhoto(photo: ProgressPhoto): Promise<{ success: boolean; cloudUrl?: string }> {
  try {
    const formData = new FormData();
    formData.append('photo', {
      uri: photo.uri,
      type: 'image/jpeg',
      name: `${photo.id}.jpg`,
    } as any);
    formData.append('date', photo.date);
    formData.append('notes', photo.notes || '');

    const response = await fetch(`${BACKEND_URL}/progress/photos`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authToken}`,
      },
      body: formData,
    });

    if (response.ok) {
      const data = await response.json();
      return { 
        success: true, 
        cloudUrl: data.data.imageUrl 
      };
    }
    
    return { success: false };
  } catch (error) {
    console.error('Upload error:', error);
    return { success: false };
  }
}
```

## Production Considerations

### For Production Deployment:

1. **Use Cloud Storage:**
   - Integrate AWS S3, Google Cloud Storage, or Cloudflare R2
   - Implement presigned URLs for secure uploads
   - Use CDN for faster image delivery

2. **Image Processing:**
   - Generate thumbnails for faster loading
   - Compress images to reduce storage
   - Convert to WebP format for better compression

3. **Security:**
   - Implement rate limiting on upload endpoint
   - Scan uploaded files for malware
   - Use HTTPS for all requests
   - Implement proper CORS policies

4. **Scalability:**
   - Use object storage instead of local filesystem
   - Implement caching for frequently accessed images
   - Consider using a queue for image processing

5. **Backup:**
   - Regular backups of uploaded files
   - Implement soft delete (keep files for recovery)
   - Consider multi-region replication

## Testing

### Test Upload:
```bash
# Create a test image
convert -size 800x600 xc:blue test.jpg

# Upload it
curl -X POST http://localhost:8080/v1/progress/photos \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "photo=@test.jpg" \
  -F "date=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  -F "notes=Test photo"
```

### Test List:
```bash
curl http://localhost:8080/v1/progress/photos \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Test Delete:
```bash
curl -X DELETE http://localhost:8080/v1/progress/photos/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Troubleshooting

### Upload fails with "validation_error"
- Check file type (must be JPEG, PNG, or WebP)
- Check file size (must be < 10MB)
- Ensure date is in ISO 8601 format

### Upload fails with "upload_failed"
- Check `STORAGE_BASE_PATH` exists and is writable
- Verify disk space is available
- Check server logs for detailed error

### Photos not accessible
- Verify `STORAGE_PUBLIC_URL` is correctly configured
- Check file permissions on upload directory
- Ensure static file serving is enabled in routes

### Database error on upload
- Verify database connection
- Check ProgressPhoto table exists
- Review foreign key constraints
