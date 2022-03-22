package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

type MergeRequestNote struct {
	GitlabId        int    `json:"id"`
	MergeRequestId  int    `json:"noteable_id"`
	MergeRequestIid int    `json:"noteable_iid"`
	NoteableType    string `json:"noteable_type"`
	Body            string
	GitlabCreatedAt core.Iso8601Time `json:"created_at"`
	Confidential    bool
	Resolvable      bool `json:"resolvable"`
	System          bool `json:"system"`
	Author          struct {
		Username string `json:"username"`
	}
}

type ApiMergeRequestNoteResponse []MergeRequestNote

func ExtractApiMergeRequestsNotes(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, _ := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_NOTES_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			mrNote := &MergeRequestNote{}
			err := json.Unmarshal(row.Data, mrNote)
			if err != nil {
				return nil, err
			}

			gitlabMergeRequestNote, err := convertMergeRequestNote(mrNote)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)

			results = append(results, gitlabMergeRequestNote)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertMergeRequestNote(mrNote *MergeRequestNote) (*models.GitlabMergeRequestNote, error) {
	gitlabMergeRequestNote := &models.GitlabMergeRequestNote{
		GitlabId:        mrNote.GitlabId,
		MergeRequestId:  mrNote.MergeRequestId,
		MergeRequestIid: mrNote.MergeRequestIid,
		NoteableType:    mrNote.NoteableType,
		AuthorUsername:  mrNote.Author.Username,
		Body:            mrNote.Body,
		GitlabCreatedAt: mrNote.GitlabCreatedAt.ToTime(),
		Confidential:    mrNote.Confidential,
		Resolvable:      mrNote.Resolvable,
		System:          mrNote.System,
	}
	return gitlabMergeRequestNote, nil
}
