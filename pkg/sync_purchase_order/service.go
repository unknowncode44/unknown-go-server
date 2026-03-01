package sync_purchase_order

import (
	"time"

	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
)

type Service interface {
	SyncOrders(req []presenter.SyncPurchaseOrderRequest) (int, error)
	GetAllOrders() ([]entities.PurchaseOrderSync, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) SyncOrders(req []presenter.SyncPurchaseOrderRequest) (int, error) {
	var orders []entities.PurchaseOrderSync
	for _, r := range req {
		// Parse date from "yyyy-mm-dd"
		t, _ := time.Parse("2006-01-02", r.FechaOC)

		orders = append(orders, entities.PurchaseOrderSync{
			CCO:             r.CCO,
			NotaPedido:      r.NotaPedido,
			FechaOC:         t,
			OCBejerman:      r.OCBejerman,
			Articulo:        r.Articulo,
			Descripcion:     r.Descripcion,
			Cantidad:        r.Cantidad,
			Moneda:          r.Moneda,
			ImporteUnitario: r.ImporteUnitario,
			Proveedor:       r.Proveedor,
		})
	}

	if err := s.repo.BatchCreate(orders); err != nil {
		return 0, err
	}
	return len(orders), nil
}

func (s *service) GetAllOrders() ([]entities.PurchaseOrderSync, error) {
	// Business logic could go here (e.g., filtering inactive records)
	return s.repo.FindAll()
}
