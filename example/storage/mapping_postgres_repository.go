package storage

import (
	"context"
	"database/sql"

	"github.com/ThreeDotsLabs/esja/example/aggregate/postcard"
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	sql2 "github.com/ThreeDotsLabs/esja/pkg/repository/sql"
)

func NewMappingPostcardRepository(ctx context.Context, db *sql.DB) (sql2.Repository[*postcard.Postcard], error) {
	return sql2.NewRepository[*postcard.Postcard](
		ctx,
		db,
		sql2.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardMapping"),
		sql2.NewMappingSerializer(
			sql2.JSONMarshaler{},
			[]sql2.EventMapper[*postcard.Postcard]{
				CreatedMapper{},
				AddressedMapper{},
				WrittenMapper{},
				SentMapper{},
			},
		),
	)
}

func NewMappingAnonymizingPostcardRepository(ctx context.Context, db *sql.DB) (sql2.Repository[*postcard.Postcard], error) {
	return sql2.NewRepository[*postcard.Postcard](
		ctx,
		db,
		sql2.NewPostgresSchemaAdapter[*postcard.Postcard]("PostcardMappingAnonymizing"),
		sql2.NewMappingSerializer(
			sql2.NewAnonymizingMarshaler(
				sql2.JSONMarshaler{},
				sql2.NewAESAnonymizer(ConstantSecretProvider{}),
			),
			[]sql2.EventMapper[*postcard.Postcard]{
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
