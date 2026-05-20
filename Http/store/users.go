package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Usuario   string        `bson:"usuario"`
	Email     string        `bson:"email"`
	Senha     string        `bson:"senha"`
	CreatedAt time.Time     `bson:"criado_em"`
}

type UserStore struct {
	collection *mongo.Collection
}

func NovoUsuario(db *mongo.Database) *UserStore {
	return &UserStore{
		collection: db.Collection("usuarios"),
	}
}

func (s *UserStore) CriarUsuario(usuario, email, senha string) error {
	existe, err := s.BuscarPorUsuario(usuario)
	if err == nil && existe != nil {
		return errors.New("Usuario já Existe")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := Users{
		Usuario:   usuario,
		Email:     email,
		Senha:     string(hash),
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = s.collection.InsertOne(ctx, u)

	return err
}

func (s *UserStore) BuscarPorUsuario(usuario string) (*Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var u Users
	err := s.collection.FindOne(ctx, bson.M{"usuario": usuario}).Decode(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (s *UserStore) ValidarSenha(usuario, senha string) (*Users, error) {

	u, err := s.BuscarPorUsuario(usuario)
	if err != nil {
		return nil, errors.New("Usuario não encontrado!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Senha), []byte(senha))
	if err != nil {
		return nil, errors.New("Senha Invalida")
	}
	return u, nil
}
