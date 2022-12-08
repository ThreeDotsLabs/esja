package storage

import (
	"context"
	"database/sql"
	postcard2 "github.com/ThreeDotsLabs/esja/example/postcard"

	"github.com/ThreeDotsLabs/esja/pkg/eventstore"
	"github.com/ThreeDotsLabs/esja/pkg/transport"
)

func NewDefaultMappingPostgresRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard2.Postcard], error) {
	return eventstore.NewSQLStore[*postcard2.Postcard](
		ctx,
		db,
		eventstore.NewMappingPostgresSQLConfig[*postcard2.Postcard](
			[]transport.EventMapper[*postcard2.Postcard]{
				CreatedMapper{},
				AddressedMapper{},
				WrittenMapper{},
				SentMapper{},
			},
		),
	)
}

func NewCustomMappingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard2.Postcard], error) {
	return eventstore.NewSQLStore[*postcard2.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[*postcard2.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[*postcard2.Postcard]("PostcardMapping"),
			Serializer: transport.NewMappingSerializer(
				transport.JSONMarshaler{},
				[]transport.EventMapper[*postcard2.Postcard]{
					CreatedMapper{},
					AddressedMapper{},
					WrittenMapper{},
					SentMapper{},
				},
			),
		},
	)
}

func NewMappingAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (eventstore.EventStore[*postcard2.Postcard], error) {
	return eventstore.NewSQLStore[*postcard2.Postcard](
		ctx,
		db,
		eventstore.SQLConfig[*postcard2.Postcard]{
			SchemaAdapter: eventstore.NewPostgresSchemaAdapter[*postcard2.Postcard]("PostcardMappingAnonymizing"),
			Serializer: transport.NewAESAnonymizingSerializer[*postcard2.Postcard](
				transport.NewMappingSerializer[*postcard2.Postcard](
					transport.JSONMarshaler{},
					[]transport.EventMapper[*postcard2.Postcard]{
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

type Created struct {
	ID string `json:"id"`
}

type CreatedMapper struct{}

func (CreatedMapper) SupportedEvent() stream.Event[*postcard2.Postcard] {
	return postcard2.Created{}
}

func (CreatedMapper) StorageEvent() any {
	return &Created{}
}

func (CreatedMapper) ToStorage(event stream.Event[*postcard2.Postcard]) any {
	e := event.(*postcard2.Created)
	return &Created{
		ID: e.ID,
	}
}

func (CreatedMapper) FromStorage(event any) stream.Event[*postcard2.Postcard] {
	e := event.(*Created)
	return &postcard2.Created{
		ID: e.ID,
	}
}

type Addressed struct {
	Sender    Address `json:"sender"`
	Addressee Address `json:"addressee"`
}

type AddressedMapper struct{}

func (AddressedMapper) SupportedEvent() stream.Event[*postcard2.Postcard] {
	return postcard2.Addressed{}
}

func (AddressedMapper) StorageEvent() any {
	return &Addressed{}
}

func (AddressedMapper) ToStorage(event stream.Event[*postcard2.Postcard]) any {
	e := event.(*postcard2.Addressed)
	return &Addressed{
		Sender:    Address(e.Sender),
		Addressee: Address(e.Addressee),
	}
}

func (AddressedMapper) FromStorage(event any) stream.Event[*postcard2.Postcard] {
	e := event.(*Addressed)
	return &postcard2.Addressed{
		Sender:    postcard2.Address(e.Sender),
		Addressee: postcard2.Address(e.Addressee),
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

func (WrittenMapper) SupportedEvent() stream.Event[*postcard2.Postcard] {
	return postcard2.Written{}
}

func (WrittenMapper) StorageEvent() any {
	return &Written{}
}

func (WrittenMapper) ToStorage(e stream.Event[*postcard2.Postcard]) any {
	ev := e.(*postcard2.Written)
	return &Written{
		Content: ev.Content,
	}
}

func (WrittenMapper) FromStorage(event any) stream.Event[*postcard2.Postcard] {
	e := event.(*Written)
	return &postcard2.Written{
		Content: e.Content,
	}
}

type Sent struct{}

type SentMapper struct{}

func (SentMapper) SupportedEvent() stream.Event[*postcard2.Postcard] {
	return postcard2.Sent{}
}

func (SentMapper) StorageEvent() any {
	return &Sent{}
}

func (SentMapper) ToStorage(event stream.Event[*postcard2.Postcard]) any {
	return &Sent{}
}

func (SentMapper) FromStorage(event any) stream.Event[*postcard2.Postcard] {
	return &postcard2.Sent{}
}
