package repository

import (
  "fmt"
  "strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
  return &Repository{}, nil
}

type Order struct { 
  ID    int 
  Title string 
  Image string
  Coeff float32
  Description string
}

type Basket struct {
	ID         int
	Components []Order
	Status     bool
}

func (r *Repository) GetBasket() ([]Basket, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Basket{}, err 
	}

	basket := []Basket{
		{
			ID:         1,
			Components: []Order{orders[0], orders[1]},
			Status:     false,
		},
		{
			ID:         2,
			Components: []Order{orders[3], orders[2], orders[1], orders[0]},
			Status:     true,
		},
	}

	return basket, nil
}

func (r *Repository) GetOrders() ([]Order, error) {
  orders := []Order{ 
    {
      ID:    1,
      Title: "Количество пользователей",
	  Image: "http://127.0.0.1:9000/pictures/Users.jpg",
	  Coeff: 10,
	  Description: `
		Показатель, отражающий общее число активных участников системы.
		Он определяет масштаб нагрузки и напрямую влияет на объем генерируемых данных.
		`,
    },
    {
      ID:    2,
      Title: "Частота операций",
	  Image: "http://127.0.0.1:9000/pictures/Frequency.jpg",
	  Coeff: 20,
	  Description: `
	  	Характеристика, показывающая, как часто выполняются вычислительные действия при обработке данных. 
		Она отражает интенсивность загрузки вычислительной системы.
	  `,
    },
    {
      ID:    3,
      Title: "Размер записей",
	  Image: "http://127.0.0.1:9000/pictures/Entry.jpg",
	  Coeff: 30,
	  Description: `
		Характеристика, задающая средний объём одной записи в таблице.
		Чем больше размер, тем быстрее растёт общий объём данных при их накоплении.
	  `,
    },
	{
      ID:    4,
      Title: "Количество файлов",
	  Image: "http://127.0.0.1:9000/pictures/Files.jpg",
	  Coeff: 40,
	  Description:`
		Показатель, определяющий число прикреплённых или загружаемых файлов.
		Он отражает интенсивность использования файлового хранилища в системе.
	  `,
    },
    {
      ID:    5,
      Title: "Количество версий записей",
	  Image: "http://127.0.0.1:9000/pictures/Entries.jpg",
	  Coeff: 50,
	  Description:`
		Характеристика, показывающая, сколько изменений или редакций сохраняется для каждой записи.
		Она определяет скорость роста данных при ведении истории изменений.
	  `,
    },
  }

  if len(orders) == 0 {
    return nil, fmt.Errorf("массив пустой")
  }

  return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err 
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil 
		}
	}
	return Order{}, fmt.Errorf("заказ не найден") 
}

func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}

	return result, nil
}
