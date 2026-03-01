package presenter

import "time"

// SyncPurchaseOrderRequest receive from Excel
type SyncPurchaseOrderRequest struct {
	CCO             string  `json:"cco"`
	NotaPedido      string  `json:"nota_pedido"`
	FechaOC         string  `json:"fecha_oc"` // Formato "yyyy-mm-dd"
	OCBejerman      string  `json:"oc_bejerman"`
	Articulo        string  `json:"articulo"`
	Descripcion     string  `json:"descripcion"`
	Cantidad        float64 `json:"cantidad"`
	ImporteUnitario float64 `json:"importe_unitario"`
	Moneda          string  `json:"moneda"`
	Proveedor       string  `json:"proveedor"` // Agregado para consistencia
}

// SyncResponse split standar response
type SyncResponse struct {
	Ok      bool        `json:"ok"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PurchaseOrderSyncResponse defines the output shape for a single record
type PurchaseOrderSyncResponse struct {
	ID              uint      `json:"id"`
	CCO             string    `json:"cco"`
	NotaPedido      string    `json:"nota_pedido"`
	FechaOC         time.Time `json:"fecha_oc"`
	OCBejerman      string    `json:"oc_bejerman"`
	Articulo        string    `json:"articulo"`
	Descripcion     string    `json:"descripcion"`
	Cantidad        float64   `json:"cantidad"`
	Proveedor       string    `json:"proveedor"`
	Moneda          string    `json:"moneda"`
	ImporteUnitario float64   `json:"importe_unitario"`
	SincronizadoEn  time.Time `json:"sincronizado_en"`
}
