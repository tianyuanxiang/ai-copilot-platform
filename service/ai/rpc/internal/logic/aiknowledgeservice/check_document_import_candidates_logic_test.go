package aiknowledgeservicelogic

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"ai-copilot-platform/ai-rpc/internal/config"
	"ai-copilot-platform/ai-rpc/internal/svc"
	"ai-copilot-platform/ai-rpc/pb"
)

func TestCheckDocumentImportCandidatesRejectsMoreThanOneHundredFiles(t *testing.T) {
	candidates := make([]*pb.DocumentImportCandidateReq, 101)
	logic := NewCheckDocumentImportCandidatesLogic(context.Background(), &svc.ServiceContext{})

	if _, err := logic.CheckDocumentImportCandidates(&pb.CheckDocumentImportCandidatesReq{
		UserId:     1,
		KbId:       1,
		Candidates: candidates,
	}); err == nil {
		t.Fatal("CheckDocumentImportCandidates() error = nil, want candidate limit error")
	}
}

func TestReadUploadedFileRejectsTraversalAndMissingFile(t *testing.T) {
	logic := newUploadReaderForTest(t)

	if _, err := logic.readUploadedFile("../outside.md"); err == nil {
		t.Fatal("readUploadedFile() traversal error = nil")
	}
	if _, err := logic.readUploadedFile("missing.md"); err == nil {
		t.Fatal("readUploadedFile() missing file error = nil")
	}
}

func TestReadUploadedFileAcceptsUploadRootPrefix(t *testing.T) {
	logic := newUploadReaderForTest(t)
	root := logic.svcCtx.Config.Storage.UploadPath
	nested := filepath.Join(root, "2026", "sample.md")
	if err := os.MkdirAll(filepath.Dir(nested), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(nested, []byte("# sample"), 0o644); err != nil {
		t.Fatal(err)
	}

	data, err := logic.readUploadedFile(filepath.Join(filepath.Base(root), "2026", "sample.md"))
	if err != nil {
		t.Fatalf("readUploadedFile() error = %v", err)
	}
	if string(data) != "# sample" {
		t.Fatalf("readUploadedFile() = %q, want %q", data, "# sample")
	}
}

func TestSupportedUploadedFileTypes(t *testing.T) {
	for _, fileType := range []string{"md", "markdown", "txt", "docx"} {
		if !isSupportedUploadedFileType(fileType) {
			t.Fatalf("isSupportedUploadedFileType(%q) = false", fileType)
		}
	}
	for _, fileType := range []string{"pdf", "xlsx", "xlsm", "xls", "csv", ""} {
		if isSupportedUploadedFileType(fileType) {
			t.Fatalf("isSupportedUploadedFileType(%q) = true", fileType)
		}
	}
}

func newUploadReaderForTest(t *testing.T) *IngestDocumentLogic {
	t.Helper()
	var cfg config.Config
	cfg.Storage.UploadPath = filepath.Join(t.TempDir(), "uploads")
	return NewIngestDocumentLogic(context.Background(), &svc.ServiceContext{Config: cfg})
}
