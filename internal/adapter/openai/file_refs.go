package openai

import "strings"

func collectOpenAIRefFileIDs(req map[string]any) []string {
	if len(req) == 0 {
		return nil
	}
	out := make([]string, 0, 4)
	seen := map[string]struct{}{}
	for _, raw := range []any{
		req["ref_file_ids"],
		req["file_ids"],
		req["attachments"],
		req["messages"],
		req["input"],
	} {
		appendOpenAIRefFileIDs(&out, seen, raw)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func appendOpenAIRefFileIDs(out *[]string, seen map[string]struct{}, raw any) {
	switch x := raw.(type) {
	case string:
		addOpenAIRefFileID(out, seen, x)
	case []string:
		for _, item := range x {
			addOpenAIRefFileID(out, seen, item)
		}
	case []any:
		for _, item := range x {
			appendOpenAIRefFileIDs(out, seen, item)
		}
	case map[string]any:
		if fileID := strings.TrimSpace(asString(x["file_id"])); fileID != "" {
			addOpenAIRefFileID(out, seen, fileID)
		}
		if strings.Contains(strings.ToLower(strings.TrimSpace(asString(x["type"]))), "file") {
			if fileID := strings.TrimSpace(asString(x["id"])); fileID != "" {
				addOpenAIRefFileID(out, seen, fileID)
			}
		}
		if fileMap, ok := x["file"].(map[string]any); ok {
			if fileID := strings.TrimSpace(asString(fileMap["file_id"])); fileID != "" {
				addOpenAIRefFileID(out, seen, fileID)
			}
			if fileID := strings.TrimSpace(asString(fileMap["id"])); fileID != "" {
				addOpenAIRefFileID(out, seen, fileID)
			}
		}
		for _, key := range []string{"ref_file_ids", "file_ids", "attachments", "messages", "input", "content", "files", "items", "data", "source"} {
			if nested, ok := x[key]; ok {
				appendOpenAIRefFileIDs(out, seen, nested)
			}
		}
	}
}

func addOpenAIRefFileID(out *[]string, seen map[string]struct{}, fileID string) {
	fileID = strings.TrimSpace(fileID)
	if fileID == "" {
		return
	}
	if _, ok := seen[fileID]; ok {
		return
	}
	seen[fileID] = struct{}{}
	*out = append(*out, fileID)
}
