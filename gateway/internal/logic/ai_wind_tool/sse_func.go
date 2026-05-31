package ai_wind_tool

import (
	"ai-copilot-platform/ai-rpc/pb"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"unicode/utf8"
)

var citationMarkerPattern = regexp.MustCompile(`\[(\d+)\]`)

func CitationsFromRPC(citations []*pb.Citation) []StreamCitation {
	items := make([]StreamCitation, 0, len(citations))
	for _, item := range citations {
		if item == nil {
			continue
		}
		items = append(items, StreamCitation{
			DocumentId: item.DocumentId,
			ChunkId:    item.ChunkId,
			Title:      item.Title,
			Snippet:    item.Snippet,
			Score:      item.Score,
		})
	}
	return items
}

func StreamReferencesFromAnswer(answer string, citations []StreamCitation) []StreamReference {
	matches := citationMarkerPattern.FindAllStringSubmatchIndex(answer, -1)
	references := make([]StreamReference, 0, len(matches))
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}
		refIndex, err := strconv.Atoi(answer[match[2]:match[3]])
		if err != nil || refIndex <= 0 || refIndex > len(citations) {
			continue
		}
		references = append(references, StreamReference{
			RefIndex: refIndex,
			RefText:  answer[match[0]:match[1]],
			Start:    utf8.RuneCountInString(answer[:match[0]]),
			End:      utf8.RuneCountInString(answer[:match[1]]),
			Citation: citations[refIndex-1],
		})
	}
	return references
}

func StreamCitationsFromRPC(citations []*pb.Citation) []StreamCitation {
	items := make([]StreamCitation, 0, len(citations))
	for _, item := range citations {
		if item == nil {
			continue
		}
		refIndex := len(items) + 1
		items = append(items, StreamCitation{
			RefIndex:   refIndex,
			RefText:    "[" + strconv.Itoa(refIndex) + "]",
			DocumentId: item.DocumentId,
			ChunkId:    item.ChunkId,
			Title:      item.Title,
			Snippet:    item.Snippet,
			Score:      item.Score,
		})
	}
	return items
}

func WriteChatSSE(w http.ResponseWriter, event StreamEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte("data: ")); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	_, err = w.Write([]byte("\n\n"))
	return err
}
