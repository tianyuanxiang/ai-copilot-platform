package ai_document

import (
	aiknowledgeclient "ai-copilot-platform/ai-rpc/client/aiknowledgeservice"
	"ai-copilot-platform/gateway/internal/types"
)

func documentItemFromRPC(doc *aiknowledgeclient.DocumentItem) types.AiDocumentItem {
	if doc == nil {
		return types.AiDocumentItem{}
	}
	return types.AiDocumentItem{
		DocumentId:   doc.DocumentId,
		KbId:         doc.KbId,
		FileName:     doc.FileName,
		FileType:     doc.FileType,
		Status:       doc.Status,
		ErrorMessage: doc.ErrorMessage,
		UploadedBy:   doc.UploadedBy,
		CreatedAt:    doc.CreatedAt,
		UpdatedAt:    doc.UpdatedAt,
	}
}
