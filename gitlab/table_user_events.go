// Package gitlab implements gitlab api calls for steampipe.
package gitlab

import (
	"context"
	"fmt"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	api "github.com/xanzy/go-gitlab"
)

// UserEvent contains a user's contribution.
type UserEvent struct {
	ID                      int
	ProjectID               int
	ActionName              string
	TargetID                int
	TargetIID               int
	TargetType              string
	AuthorID                int
	TargetTitle             string
	CreatedAt               *time.Time
	PushDataCommitCount     int
	PushDataAction          string
	PushDataRefType         string
	PushDataCommitFrom      string
	PushDataCommitTo        string
	PushDataRef             string
	PushDataCommitTitle     string
	NoteID                  int
	NoteType                string
	NoteBody                string
	NoteAttachment          string
	NoteTitle               string
	NoteFileName            string
	NoteAuthorID            int
	NoteExpiresAt           *time.Time
	NoteUpdatedAt           *time.Time
	NoteCreatedAt           *time.Time
	NoteNoteableID          int
	NoteNoteableType        string
	NoteCommitID            string
	NoteResolvable          bool
	NoteResolved            bool
	NoteResolvedByID        int
	NoteResolvedByUsername  string
	NoteResolvedByEmail     string
	NoteResolvedByName      string
	NoteResolvedByState     string
	NoteResolvedByAvatarURL string
	NoteResolvedByWebURL    string
	NoteResolvedAt          *time.Time
	NoteNoteableIID         int
}

func tableUserEvents() *plugin.Table {
	return &plugin.Table{
		Name:        "gitlab_user_events",
		Description: "Obtain information about a user's events.",
		List: &plugin.ListConfig{
			/*
				KeyColumns: []*plugin.KeyColumn{
					{
						Name:    "author_id",
						Require: plugin.Required,
					},
					{
						Name:      "created_at",
						Require:   plugin.Optional,
						Operators: []string{">", ">=", "=", "<", "<="},
					},
					{
						Name:      "target_type",
						Require:   plugin.Optional,
						Operators: []string{"="},
					},
					{
						Name:      "action",
						Require:   plugin.Optional,
						Operators: []string{"="},
					},
				},
			*/
			KeyColumns: plugin.SingleColumn("author_id"),
			Hydrate:    listUserEvents,
		},
		Columns: userEventColumns(),
	}
}

// Hydrate Functions
func listUserEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("listUserEvents", "started")
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listUserEvents", "unable to establish a connection", err)
		return nil, fmt.Errorf("unable to establish a connection: %v", err)
	}

	userID := int(d.EqualsQuals["author_id"].GetInt64Value())
	opt := &api.ListContributionEventsOptions{ListOptions: api.ListOptions{
		Page:    1,
		PerPage: 50,
	}}

	/*
		if d.Quals["created_at"] != nil {
			for _, q := range d.Quals["created_at"].Quals {
				givenTime := q.Value.GetTimestampValue().AsTime()
				beforeTime := givenTime.Add(time.Duration(-1) * time.Second)
				afterTime := givenTime.Add(time.Second * 1)
				givenISOTime := gitlab.ISOTime(givenTime)
				beforeISOTime := gitlab.ISOTime(beforeTime)
				afterISOTime := gitlab.ISOTime(afterTime)

				switch q.Operator {
				case ">":
					opt.After = &afterISOTime
				case ">=":
					opt.After = &givenISOTime
				case "=":
					opt.After = &beforeISOTime
					opt.Before = &afterISOTime
				case "<=":
					opt.Before = &givenISOTime
				case "<":
					opt.Before = &beforeISOTime
				}
			}
		}

		if d.Quals["target_type"] != nil {
			targetType := gitlab.EventTargetTypeValue(d.EqualsQuals["target_type"].GetStringValue())
			opt.TargetType = &targetType
		}

		if d.Quals["action"] != nil {
			action := gitlab.EventTypeValue(d.EqualsQuals["action"].GetStringValue())
			opt.Action = &action
		}
	*/

	for {
		plugin.Logger(ctx).Debug("listUserEvents", "userID", userID, "page", opt.Page, "perPage", opt.PerPage)
		events, resp, err := conn.Users.ListUserContributionEvents(userID, opt)
		if err != nil {
			plugin.Logger(ctx).Error("listUserEvents", "userID", userID, "page", opt.Page, "error", err)
			return nil, fmt.Errorf("unable to obtain events for user_id %d\n%v", userID, err)
		}

		for _, event := range events {
			userEvent := &UserEvent{
				ID:          event.ID,
				ProjectID:   event.ProjectID,
				ActionName:  event.ActionName,
				TargetID:    event.TargetID,
				TargetIID:   event.TargetIID,
				TargetType:  event.TargetType,
				AuthorID:    event.AuthorID,
				TargetTitle: event.TargetTitle,
				CreatedAt:   event.CreatedAt,
			}

			//if event.PushData != nil {
			userEvent.PushDataCommitCount = event.PushData.CommitCount
			userEvent.PushDataAction = event.PushData.Action
			userEvent.PushDataRefType = event.PushData.RefType
			userEvent.PushDataCommitFrom = event.PushData.CommitFrom
			userEvent.PushDataCommitTo = event.PushData.CommitTo
			userEvent.PushDataRef = event.PushData.Ref
			userEvent.PushDataCommitTitle = event.PushData.CommitTitle
			//}
			if event.Note != nil {
				userEvent.NoteID = event.Note.ID
				userEvent.NoteType = string(event.Note.Type)
				userEvent.NoteBody = event.Note.Body
				userEvent.NoteAttachment = event.Note.Attachment
				userEvent.NoteTitle = event.Note.Title
				userEvent.NoteFileName = event.Note.FileName
				userEvent.NoteAuthorID = event.Note.Author.ID
				userEvent.NoteExpiresAt = event.Note.ExpiresAt
				userEvent.NoteUpdatedAt = event.Note.UpdatedAt
				userEvent.NoteCreatedAt = event.Note.CreatedAt
				userEvent.NoteNoteableID = event.Note.NoteableID
				userEvent.NoteNoteableType = event.Note.NoteableType
				userEvent.NoteCommitID = event.Note.CommitID
				userEvent.NoteResolvable = event.Note.Resolvable
				userEvent.NoteResolved = event.Note.Resolved
				//if event.Note.ResolvedBy != nil {
				userEvent.NoteResolvedByID = event.Note.ResolvedBy.ID
				userEvent.NoteResolvedByUsername = event.Note.ResolvedBy.Username
				userEvent.NoteResolvedByEmail = event.Note.ResolvedBy.Email
				userEvent.NoteResolvedByName = event.Note.ResolvedBy.Name
				userEvent.NoteResolvedByState = event.Note.ResolvedBy.State
				userEvent.NoteResolvedByAvatarURL = event.Note.ResolvedBy.AvatarURL
				userEvent.NoteResolvedByWebURL = event.Note.ResolvedBy.WebURL
				//}
				userEvent.NoteResolvedAt = event.Note.ResolvedAt
				userEvent.NoteNoteableIID = event.Note.NoteableIID
			}

			d.StreamListItem(ctx, userEvent)
		}

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	plugin.Logger(ctx).Debug("listUserEvents", "completed successfully")
	return nil, nil
}

// Column Function
func userEventColumns() []*plugin.Column {
	return []*plugin.Column{
		{
			Name:        "id",
			Type:        proto.ColumnType_INT,
			Description: "id",
		},
		{
			Name:        "project_id",
			Type:        proto.ColumnType_INT,
			Description: "project_id",
		},
		{
			Name:        "action_name",
			Type:        proto.ColumnType_STRING,
			Description: "action_name",
		},
		{
			Name:        "target_id",
			Type:        proto.ColumnType_INT,
			Description: "target_id",
		},
		{
			Name:        "target_iid",
			Type:        proto.ColumnType_INT,
			Description: "target_iid",
		},
		{
			Name:        "target_type",
			Type:        proto.ColumnType_STRING,
			Description: "target_type",
		},
		{
			Name:        "author_id",
			Type:        proto.ColumnType_INT,
			Description: "author_id",
		},
		{
			Name:        "target_title",
			Type:        proto.ColumnType_STRING,
			Description: "target_title",
		},
		{
			Name:        "created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "created_at",
		},
		{
			Name:        "push_data_commit_count",
			Type:        proto.ColumnType_INT,
			Description: "push_data_commit_count",
		},
		{
			Name:        "push_data_action",
			Type:        proto.ColumnType_STRING,
			Description: "push_data_action",
		},
		{
			Name:        "push_data_ref_type",
			Type:        proto.ColumnType_STRING,
			Description: "push_data_ref_type",
		},
		{
			Name:        "push_data_commit_from",
			Type:        proto.ColumnType_STRING,
			Description: "push_data_commit_from",
		},
		{
			Name:        "push_data_commit_to",
			Type:        proto.ColumnType_STRING,
			Description: "push_data_commit_to",
		},
		{
			Name:        "push_data_ref",
			Type:        proto.ColumnType_STRING,
			Description: "push_data_ref",
		},
		{
			Name:        "push_data_commit_title",
			Type:        proto.ColumnType_STRING,
			Description: "push_data_commit_title",
		},
		{
			Name:        "note_id",
			Type:        proto.ColumnType_INT,
			Description: "note_id",
		},
		{
			Name:        "note_type",
			Type:        proto.ColumnType_STRING,
			Description: "note_type",
		},
		{
			Name:        "note_body",
			Type:        proto.ColumnType_STRING,
			Description: "note_body",
		},
		{
			Name:        "note_attachment",
			Type:        proto.ColumnType_STRING,
			Description: "note_attachment",
		},
		{
			Name:        "note_title",
			Type:        proto.ColumnType_STRING,
			Description: "note_title",
		},
		{
			Name:        "note_file_name",
			Type:        proto.ColumnType_STRING,
			Description: "note_file_name",
		},
		{
			Name:        "note_author_id",
			Type:        proto.ColumnType_INT,
			Description: "note_author_id",
		},
		{
			Name:        "note_expires_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "note_expires_at",
		},
		{
			Name:        "note_updated_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "note_updated_at",
		},
		{
			Name:        "note_created_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "note_created_at",
		},
		{
			Name:        "note_noteable_id",
			Type:        proto.ColumnType_INT,
			Description: "note_noteable_id",
		},
		{
			Name:        "note_noteable_type",
			Type:        proto.ColumnType_STRING,
			Description: "note_noteable_type",
		},
		{
			Name:        "note_commit_id",
			Type:        proto.ColumnType_STRING,
			Description: "note_commit_id",
		},
		{
			Name:        "note_resolvable",
			Type:        proto.ColumnType_BOOL,
			Description: "note_resolvable",
		},
		{
			Name:        "note_resolved",
			Type:        proto.ColumnType_BOOL,
			Description: "note_resolved",
		},
		{
			Name:        "note_resolved_by_id",
			Type:        proto.ColumnType_INT,
			Description: "note_resolved_by_id",
		},
		{
			Name:        "note_resolved_by_username",
			Type:        proto.ColumnType_STRING,
			Description: "note_resolved_by_username",
		},
		{
			Name:        "note_resolved_by_email",
			Type:        proto.ColumnType_STRING,
			Description: "note_resolved_by_email",
		},
		{
			Name:        "note_resolved_by_name",
			Type:        proto.ColumnType_STRING,
			Description: "note_resolved_by_name",
		},
		{
			Name:        "note_resolved_by_state",
			Type:        proto.ColumnType_STRING,
			Description: "note_resolved_by_state",
		},
		{
			Name:        "note_resolved_by_avatar_url",
			Type:        proto.ColumnType_STRING,
			Description: "note_resolved_by_avatar_url",
		},
		{
			Name:        "note_resolved_by_web_url",
			Type:        proto.ColumnType_STRING,
			Description: "note_resolved_by_web_url",
		},
		{
			Name:        "note_resolved_at",
			Type:        proto.ColumnType_TIMESTAMP,
			Description: "note_resolved_at",
		},
		{
			Name:        "note_noteable_iid",
			Type:        proto.ColumnType_INT,
			Description: "note_noteable_iid",
		},
	}
}
