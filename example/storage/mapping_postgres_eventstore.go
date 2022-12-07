package storage

import (
	"context"
	"database/sql"

	"github.com/ThreeDotsLabs/esja/example/aggregate/postcard"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/eventstore"
	"github.com/ThreeDotsLabs/esja/pkg/transport"
)

func NewDefaultMappingPostgresRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard.Postcard], error) {
	return eventstore.NewSQLStore[*postcard.Postcard](
		ctx,
		db,
		eventstore.NewMappingPostgresSQLConfig[*postcard.Postcard](
			[]transport.EventMapper[*postcard.Postcard]{
				CreatedMapper{},
				AddressedMapper{},
				WrittenMapper{},
				SentMapper{},
			},
		),
	)
}

func NewCustomMappingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard.Postcard], error) {
	return eventstore.NewSQLStore[*postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[*postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardMapping"),
			Serializer: transport.NewMappingSerializer(
				transport.JSONMarshaler{},
				[]transport.EventMapper[*postcard.Postcard]{
					CreatedMapper{},
					AddressedMapper{},
					WrittenMapper{},
					SentMapper{},
				},
			),
		},
	)
}

func NewMappingAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard.Postcard], error) {
	return eventstore.NewSQLStore[*postcard.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[*postcard.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardMappingAnonymizing"),
			Serializer: transport.NewAESAnonymizingSerializer[*postcard.Postcard](
				transport.NewMappingSerializer[*postcard.Postcard](
					transport.JSONMarshaler{},
					[]transport.EventMapper[*postcard.Postcard]{
						CreatedMapper{},
						AddressedMapper{},
						WrittenMapper{},
						SentMapper{},
					},
				),
				ConstantSecretProvider{},
			),
		},
	)
}

func NewMappingSQLiteRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard.Postcard], error) {
	return eventstore.NewSQLStore[*postcard.Postcard](
		ctx,
		db,
		eventstore.NewMappingSQLiteConfig[*postcard.Postcard](
			[]transport.EventMapper[*postcard.Postcard]{
				CreatedMapper{},
				AddressedMapper{},
				WrittenMapper{},
				SentMapper{},
			},
		),
	)
}

type Created struct {
	ID string `json:"id"`
}

type CreatedMapper struct{}

func (CreatedMapper) SupportedEvent() aggregate.Event[*postcard.Postcard] {
	return postcard.Created{}
}

func (CreatedMapper) StorageEvent() any {
	return &Created{}
}

func (CreatedMapper) ToStorage(event aggregate.Event[*postcard.Postcard]) any {
	e := event.(*postcard.Created)
	return &Created{
		ID: e.ID,
	}
}

func (CreatedMapper) FromStorage(event any) aggregate.Event[*postcard.Postcard] {
	e := event.(*Created)
	return &postcard.Created{
		ID: e.ID,
	}
}

type Addressed struct {
	Sender    Address `json:"sender"`
	Addressee Address `json:"addressee"`
}

type AddressedMapper struct{}

func (AddressedMapper) SupportedEvent() aggregate.Event[*postcard.Postcard] {
	return postcard.Addressed{}
}

func (AddressedMapper) StorageEvent() any {
	return &Addressed{}
}

func (AddressedMapper) ToStorage(event aggregate.Event[*postcard.Postcard]) any {
	e := event.(*postcard.Addressed)
	return &Addressed{
		Sender:    Address(e.Sender),
		Addressee: Address(e.Addressee),
	}
}

func (AddressedMapper) FromStorage(event any) aggregate.Event[*postcard.Postcard] {
	e := event.(*Addressed)
	return &postcard.Addressed{
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

type WrittenMapper struct{}

func (WrittenMapper) SupportedEvent() aggregate.Event[*postcard.Postcard] {
	return postcard.Written{}
}

func (WrittenMapper) StorageEvent() any {
	return &Written{}
}

func (WrittenMapper) ToStorage(e aggregate.Event[*postcard.Postcard]) any {
	ev := e.(*postcard.Written)
	return &Written{
		Content: ev.Content,
	}
}

func (WrittenMapper) FromStorage(event any) aggregate.Event[*postcard.Postcard] {
	e := event.(*Written)
	return &postcard.Written{
		Content: e.Content,
	}
}

type Sent struct{}

type SentMapper struct{}

func (SentMapper) SupportedEvent() aggregate.Event[*postcard.Postcard] {
	return postcard.Sent{}
}

func (SentMapper) StorageEvent() any {
	return &Sent{}
}

func (SentMapper) ToStorage(event aggregate.Event[*postcard.Postcard]) any {
	return &Sent{}
}

func (SentMapper) FromStorage(event any) aggregate.Event[*postcard.Postcard] {
	return &postcard.Sent{}
}
