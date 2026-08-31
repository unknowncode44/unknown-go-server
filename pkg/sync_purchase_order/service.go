package sync_purchase_order

import (
	"time"

	"github.com/google/uuid"
	"github.com/unknowncode44/unknown-go-server/api/presenter"
	"github.com/unknowncode44/unknown-go-server/pkg/entities"
	"github.com/unknowncode44/unknown-go-server/pkg/material"
	"github.com/unknowncode44/unknown-go-server/pkg/vendor"
)

type Service interface {
	SyncOrders(req []presenter.SyncPurchaseOrderRequest) (int, error)
	GetAllOrders() ([]entities.PurchaseOrderSync, error)
	AssignLinks(id uint, materialId uuid.UUID, taxId string) (*entities.PurchaseOrderSync, error)
}

type service struct {
	repo  Repository
	mRepo material.Repository
	vRepo vendor.Repository
}

func NewService(r Repository, materialRepo material.Repository, vendorRepo vendor.Repository) Service {
	return &service{
		repo:  r,
		mRepo: materialRepo,
		vRepo: vendorRepo,
	}
}

// poKey es la clave natural de una línea de OC: el mismo par de columnas
// sobre el que corren el uniqueIndex de la entidad y el upsert del repo.
type poKey struct {
	ocBejerman string
	articulo   string
}

func (s *service) SyncOrders(req []presenter.SyncPurchaseOrderRequest) (int, error) {
	// El export de Excel puede traer la misma línea (oc_bejerman + articulo)
	// más de una vez dentro del MISMO envío — p. ej. filas todavía sin número
	// de OC asignado, que quedan todas con oc_bejerman vacío. Postgres rechaza
	// un INSERT ... ON CONFLICT DO UPDATE que afecte la misma fila dos veces en
	// una sola sentencia ("command cannot affect row a second time",
	// SQLSTATE 21000), así que hay que deduplicar acá, antes del batch.
	// Gana la última ocurrencia: es la versión que el Excel deja vigente.
	var orders []entities.PurchaseOrderSync
	seen := make(map[poKey]int, len(req))

	for _, r := range req {
		// Parse date from "yyyy-mm-dd"
		t, _ := time.Parse("2006-01-02", r.FechaOC)

		order := entities.PurchaseOrderSync{
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
		}

		k := poKey{ocBejerman: r.OCBejerman, articulo: r.Articulo}
		if i, dup := seen[k]; dup {
			orders[i] = order
			continue
		}
		seen[k] = len(orders)
		orders = append(orders, order)
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

func (s *service) AssignLinks(id uint, materialId uuid.UUID, taxId string) (*entities.PurchaseOrderSync, error) {
	m, err := s.mRepo.FindByID(materialId)
	if err != nil {
		return nil, err
	}
	v, err := s.vRepo.FindByTaxID(taxId)
	if err != nil {
		return nil, err
	}
	po, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdateLinks(po, &m.ID, &v.ID)
}
