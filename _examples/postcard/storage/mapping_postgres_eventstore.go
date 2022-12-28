package storage

import (
	"context"
	"database/sql"

	"github.com/ThreeDotsLabs/esja/eventstore"
	"github.com/ThreeDotsLabs/esja/stream"
	"github.com/ThreeDotsLabs/esja/transport"

	"postcard"
)

func NewDefaultMappingPostgresRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.NewMappingPostgresSQLConfig[postcard.Postcard](
			[]transport.Event[postcard.Postcard]{
				Created{},
				Addressed{},
				Written{},
				Sent{},
			},
		),
	)
}

func NewCustomMappingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[postcard.Postcard](),
			Mapper: transport.NewDefaultMapper(
				[]transport.Event[postcard.Postcard]{
					Created{},
					Addressed{},
					Written{},
					Sent{},
				},
			),
			Marshaler: transport.JSONMarshaler{},
		},
	)
}

func NewMappingAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[postcard.Postcard](),
			Mapper: transport.NewAESAnonymizer[postcard.Postcard](
				transport.NewDefaultMapper[postcard.Postcard](
					[]transport.Event[postcard.Postcard]{
						Created{},
						Addressed{},
						Written{},
						Sent{},
					},
				),
				ConstantSecretProvider{},
			),
			Marshaler: transport.JSONMarshaler{},
		},
	)
}

func NewMappingSQLitePostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[postcard.Postcard], error) {
	return eventstore.NewSQLStore[postcard.Postcard](
		ctx,
		db,
		eventstore.NewMappingSQLiteConfig[postcard.Postcard](
			[]transport.Event[postcard.Postcard]{
				Created{},
				Addressed{},
				Written{},
				Sent{},
			},
		),
	)
}

type Created struct {
	ID string `json:"id"`
}

func (e Created) NewFromEvent(event stream.Event[postcard.Postcard]) transport.Event[postcard.Postcard] {
	created := event.(postcard.Created)
	e.ID = created.ID
	return e
}

func (e Created) ToEvent() stream.Event[postcard.Postcard] {
	return postcard.Created{
		ID: e.ID,
	}
}

type Addressed struct {
	Sender    Address `json:"sender"`
	Addressee Address `json:"addressee"`
}

func (e Addressed) NewFromEvent(event stream.Event[postcard.Postcard]) transport.Event[postcard.Postcard] {
	addressed := event.(postcard.Addressed)
	e.Sender = Address(addressed.Sender)
	e.Addressee = Address(addressed.Addressee)
	return e
}

func (e Addressed) ToEvent() stream.Event[postcard.Postcard] {
	return postcard.Addressed{
		Sender:    postcard.Address(e.Sender),
		Addressee: postcard.Address(e.Addressee),
	}
}

type Address struct {
	Name  string `json:"name" anonymize:"true"`
	Line1 string `json:"line1"`
	Line2 string `json:"line2"`
	Line3 string `json:"line3"`
}

type Written struct {
	Content string `json:"content"`
}

func (e Written) NewFromEvent(event stream.Event[postcard.Postcard]) transport.Event[postcard.Postcard] {
	written := event.(postcard.Written)
	e.Content = written.Content
	return e
}

func (e Written) ToEvent() stream.Event[postcard.Postcard] {
	return postcard.Written{
		Content: e.Content,
	}
}

type Sent struct{}

func (e Sent) NewFromEvent(event stream.Event[postcard.Postcard]) transport.Event[postcard.Postcard] {
	_ = event.(postcard.Sent)
	return e
}

func (e Sent) ToEvent() stream.Event[postcard.Postcard] {
	return postcard.Sent{}
}
