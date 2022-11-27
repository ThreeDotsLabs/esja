package storage

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ThreeDotsLabs/esja/example/aggregate/postcard"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/repository"
)

// NewPostcardRepository returns a new postgres repository for the Postcard aggregate.
// Events are marshaled using the provided marshalers.
// Marshalers should use dedicated models to marshal/unmarshal events and map the domain events to them.
// This is to decouple the implementation of the aggregate from the storage layer.
func NewPostcardRepository(ctx context.Context, db *sql.DB) (repository.PostgresRepository[*postcard.Postcard], error) {
	return repository.NewPostgresRepository[*postcard.Postcard](
		ctx,
		db,
		"Postcard",
		[]repository.EventMarshaler[*postcard.Postcard]{
			CreatedMarshaler{},
			AddressedMarsahler{},
			WritterMarshaler{},
			SentMarshaler{},
		},
	)
}

type Created struct {
	ID string `json:"id"`
}

type CreatedMarshaler struct{}

func (CreatedMarshaler) SupportedEvent() aggregate.Event[*postcard.Postcard] {
	return postcard.Created{}
}

func (CreatedMarshaler) Marshal(e aggregate.Event[*postcard.Postcard]) ([]byte, error) {
	ev := e.(postcard.Created)
	c := Created{
		ID: ev.ID,
	}

	return json.Marshal(c)
}

func (CreatedMarshaler) Unmarshal(payload []byte) (aggregate.Event[*postcard.Postcard], error) {
	var c Created
	err := json.Unmarshal(payload, &c)
	if err != nil {
		return nil, err
	}

	return postcard.Created{
		ID: c.ID,
	}, nil
}

type Addressed struct {
	Sender    Address `json:"sender"`
	Addressee Address `json:"addressee"`
}

type AddressedMarsahler struct{}

func (AddressedMarsahler) SupportedEvent() aggregate.Event[*postcard.Postcard] {
	return postcard.Addressed{}
}

func (AddressedMarsahler) Marshal(e aggregate.Event[*postcard.Postcard]) ([]byte, error) {
	ev := e.(postcard.Addressed)
	a := Addressed{
		Sender:    Address(ev.Sender),
		Addressee: Address(ev.Addressee),
	}
	return json.Marshal(a)
}

func (AddressedMarsahler) Unmarshal(payload []byte) (aggregate.Event[*postcard.Postcard], error) {
	var a Addressed
	err := json.Unmarshal(payload, &a)
	if err != nil {
		return nil, err
	}

	return postcard.Addressed{
		Sender:    postcard.Address(a.Sender),
		Addressee: postcard.Address(a.Addressee),
	}, nil
}

type Address struct {
	Name  string `json:"name"`
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	Line3 string `json:"line3"`
}

type Written struct {
	Content string `json:"content"`
}

type WritterMarshaler struct{}

func (WritterMarshaler) SupportedEvent() aggregate.Event[*postcard.Postcard] {
	return postcard.Written{}
}

func (WritterMarshaler) Marshal(e aggregate.Event[*postcard.Postcard]) ([]byte, error) {
	ev := e.(postcard.Written)
	w := Written{
		Content: ev.Content,
	}
	return json.Marshal(w)
}

func (WritterMarshaler) Unmarshal(payload []byte) (aggregate.Event[*postcard.Postcard], error) {
	var w Written
	err := json.Unmarshal(payload, &w)
	if err != nil {
		return nil, err
	}

	return postcard.Written{
		Content: w.Content,
	}, nil
}

type Sent struct{}

type SentMarshaler struct{}

func (SentMarshaler) SupportedEvent() aggregate.Event[*postcard.Postcard] {
	return postcard.Sent{}
}

func (SentMarshaler) Marshal(e aggregate.Event[*postcard.Postcard]) ([]byte, error) {
	return []byte("{}"), nil
}

func (SentMarshaler) Unmarshal(payload []byte) (aggregate.Event[*postcard.Postcard], error) {
	return postcard.Sent{}, nil
}
