package examples

import (
	"github.com/invopop/jsonschema/examples/nested"
)

// User is used as a base to provide tests for comments.
// Don't forget to checkout the nested path.
type User struct {
	// Unique sequential identifier.
	ID int `json:"id" jsonschema:"required"`
	// This comment will be ignored
	Name    string `json:"name" jsonschema:"required,minLength=1,maxLength=20,pattern=.*,description=this is a property,title=the name,example=joe,example=lucy,default=alex"`
	Friends []struct {
		// The ID of the friend
		FriendID int `json:"friend_id"`
		// A note about this friend
		FriendNote string `json:"friend_note,omitempty"`
	} `json:"friends,omitempty" jsonschema_description:"list of friends, omitted when empty"`
	Tags map[string]any `json:"tags,omitempty"`

	// An array of pets the user cares for.
	Pets nested.Pets `json:"pets"`

	// Set of animal names to pets
	NamedPets nested.NamedPets `json:"named_pets"`

	// Set of plants that the user likes
	Plants []*nested.Plant `json:"plants" jsonschema:"title=Plants"`

	// Additional data about this user
	AdditionalData struct {
		// This user's favorite color
		FavoriteColor string `json:"favorite_color"`
	} `json:"additional_data"`

	// A mapping from friend IDs to notes
	FriendToNote map[int]*struct {
		// The note for this friend
		Note string `json:"note"`
	} `json:"friend_to_note,omitempty"`

	nested.Metadata[string]
}
