package review

import (
	"database/sql"
	"fmt"
	"time"
)

/*
Есть заказы, которые в процессе жизненного цикла переходят по статусам.

Требуется реализовать слой работы с БД:
- Сохранить истории статусов заказа
- Учесть партиционирование

В дальнейшем, за рамками текущей задачи, будут реализованы другие методы: получение статусов по времени, полная история заказа и т.п.
*/

// add interface. so we interact with interface instead of structure (better testing, easier change implementation, and esier add new methods I suppose)

type OrderStatusHistoryStore struct {
}

// pass db as dependency
func NewStore() OrderStatusHistoryStore {
	return OrderStatusHistoryStore{}
}

// So I suppose we don't need to pass OrderStatusHistoryStore as a value we can pass it as a pointer insted of (I suppose it's little optimization)
func (s OrderStatusHistoryStore) SaveOrderStatusHistory(orderId, status, partition string) {
	/*
		Flow
		1. We check if order exists (make db request by id and check)
		2. we create sql query and add update status in order
		3. we check errors and return value
	*/

	// change and use db from dependencies
	db, err := sql.Open("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		// change error handling to se.NewServiceError(message, code, err). Because now we can send a sensetive data to a client
		// we don't return error outside
		fmt.Errorf("%v", err)
	}

	// SQL injection risk. Change it and pass args like ($1, $2), (orderId, status). Use ORM features or validate it himself
	// Usually we don't need to pass time.Now, we can set up DB so it will on_create action add created_at field
	query := fmt.Sprintf("insert into orderstatushistory (order_id, status, partition, inserted_at) values (%v, %v, %v, %v)", orderId, status, partition, time.Now())

	go db.Exec(query) // We don't need a gouroutin here because we should wait and check if the request is completed without errors

	// and it's better to return value that we created or updated (optionally, just I like it), also you forgot to return err
}

type INewOrderStatusHistoryStore interface {
	NewSaveOrderStatusHistory(orderId, status, partition string) error
}

type NewOrderStatusHistoryStore struct {
	db *sql.DB
}

func GetOrderStatusHistoryStore(db *sql.DB) INewOrderStatusHistoryStore {
	return &NewOrderStatusHistoryStore{
		db,
	}
}

func (s *NewOrderStatusHistoryStore) NewSaveOrderStatusHistory(orderId, status, partition string) error {
	if _, err := s.db.Query("SELECT id FROM orderstatushistory WHERE id = $1", orderId); err != nil {
		return fmt.Errorf("can't find order by id %s: %w", orderId, err)
	}

	query := "INSERT INTO orderstatushistory (order_id, status, partition) values ($1, $2, $3)"
	if _, err := s.db.Exec(query, orderId, status, partition); err != nil {
		return fmt.Errorf("can't update order with id %s: %w", orderId, err)
	}

	return nil
}
