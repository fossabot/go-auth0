package management

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/auth0/go-auth0"
)

func TestJobManager_VerifyEmail(t *testing.T) {
	user := givenAUser(t)
	job := &Job{UserID: user.ID}

	err := m.Job.VerifyEmail(job)
	assert.NoError(t, err)

	actualJob, err := m.Job.Read(job.GetID())
	assert.NoError(t, err)
	assert.Equal(t, job.GetID(), actualJob.GetID())
}

func TestJobManager_ExportUsers(t *testing.T) {
	givenAUser(t)
	conn, err := m.Connection.ReadByName("Username-Password-Authentication")
	assert.NoError(t, err)

	job := &Job{
		ConnectionID: conn.ID,
		Format:       auth0.String("json"),
		Limit:        auth0.Int(5),
		Fields: []map[string]interface{}{
			{"name": "name"},
			{"name": "email"},
			{"name": "identities[0].connection"},
		},
	}

	err = m.Job.ExportUsers(job)
	assert.NoError(t, err)
}

func TestJobManager_ImportUsers(t *testing.T) {
	conn, err := m.Connection.ReadByName("Username-Password-Authentication")
	assert.NoError(t, err)

	job := &Job{
		ConnectionID:        conn.ID,
		Upsert:              auth0.Bool(true),
		SendCompletionEmail: auth0.Bool(false),
		Users: []map[string]interface{}{
			{"email": "auzironian@example.com", "email_verified": true},
		},
	}
	err = m.Job.ImportUsers(job)
	assert.NoError(t, err)

	t.Cleanup(func() {
		users, err := m.User.ListByEmail("auzironian@example.com")
		assert.NoError(t, err)
		assert.Len(t, users, 1)

		cleanupUser(t, users[0].GetID())
	})
}
