package elements

import "time"

type SubdivxMovie struct {
	Title          string    `json:"title,omitempty"`
	Url            string    `json:"url,omitempty"`
	Description    string    `json:"description,omitempty"`
	DownloadsCount int       `json:"downloadsCount,omitempty"`
	Cds            int       `json:"cds,omitempty"`
	CommentsCount  int       `json:"commentsCount,omitempty"`
	Format         string    `json:"format,omitempty"`
	UploadedBy     string    `json:"uploadedBy,omitempty"`
	UploadedDate   time.Time `json:"uploadedDate,omitempty"`
}
