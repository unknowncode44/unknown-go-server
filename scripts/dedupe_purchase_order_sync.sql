-- Limpieza puntual de duplicados en purchase_order_syncs, previa al deploy
-- del fix de idempotencia (branch fix/po-sync-idempotency).
--
-- Por qué hace falta: BatchCreate insertaba ciegamente (sin upsert), y el
-- macro VBA reenvía TODAS las filas del rango de fechas elegido en cada
-- corrida (no lleva ningún flag de "ya enviado"). Cualquier rango solapado
-- entre corridas generó filas duplicadas para la misma línea de OC.
--
-- Clave natural real: (oc_bejerman, articulo) — el número de OC se repite
-- una vez por cada artículo que incluye esa orden, así que ninguno de los
-- dos solo alcanza. Confirmado por el dueño del proceso de compras.
--
-- IMPORTANTE — orden de despliegue: este script debe correr y comitearse
-- ANTES de desplegar el código con el uniqueIndex nuevo. AutoMigrate hace
-- log.Fatalf si falla (db/postgres.go), y crear un unique index sobre datos
-- que todavía lo violan tira el servidor entero abajo al reiniciar.
--
-- Uso: correr cada bloque a mano, revisando el resultado antes de avanzar
-- al siguiente. No pensado para correrse de un tirón sin mirar.

-- ============================================================
-- PASO 0 (solo lectura) — chequeo de riesgo antes de agrupar por la clave
-- Si esto devuelve filas, NO sigas sin revisarlas a mano: significa que hay
-- registros sin oc_bejerman cargado (p. ej. una Nota de Pedido todavía sin
-- convertir a Orden de Compra en Bejerman). Agruparlos junto a otro
-- "" + mismo articulo los trataría como duplicados y se borrarían
-- propuestas distintas que legítimamente no tienen OC todavía.
-- ============================================================
SELECT id, cco, nota_pedido, articulo, descripcion, sincronizado_en
FROM purchase_order_syncs
WHERE oc_bejerman IS NULL OR trim(oc_bejerman) = '';

-- ============================================================
-- PASO 1 (solo lectura) — ver qué se va a borrar antes de tocar nada
-- ============================================================
SELECT oc_bejerman, articulo, COUNT(*) AS filas_duplicadas,
       array_agg(id ORDER BY id) AS ids
FROM purchase_order_syncs
GROUP BY oc_bejerman, articulo
HAVING COUNT(*) > 1
ORDER BY filas_duplicadas DESC;

-- ============================================================
-- PASO 2 (solo lectura) — cuántas filas se eliminarían en total
-- ============================================================
WITH ranked AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY oc_bejerman, articulo
               ORDER BY
                   (material_id IS NOT NULL OR vendor_id IS NOT NULL) DESC,
                   sincronizado_en DESC,
                   id DESC
           ) AS rn
    FROM purchase_order_syncs
)
SELECT COUNT(*) AS filas_a_eliminar FROM ranked WHERE rn > 1;

-- ============================================================
-- PASO 3 — borrado real, dentro de una transacción
-- Criterio de qué fila se conserva por cada (oc_bejerman, articulo):
--   1) la que ya tenga material_id o vendor_id asignado a mano
--      (no perder vínculos ya cargados desde la vista de Purchase Orders),
--   2) si hay empate o ninguna tiene vínculo, la más recientemente
--      sincronizada,
--   3) si sigue empatado, el id más alto.
-- ============================================================
BEGIN;

WITH ranked AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY oc_bejerman, articulo
               ORDER BY
                   (material_id IS NOT NULL OR vendor_id IS NOT NULL) DESC,
                   sincronizado_en DESC,
                   id DESC
           ) AS rn
    FROM purchase_order_syncs
)
DELETE FROM purchase_order_syncs
WHERE id IN (SELECT id FROM ranked WHERE rn > 1);

-- Verificar antes de confirmar: no debería quedar ningún grupo con COUNT(*) > 1
SELECT oc_bejerman, articulo, COUNT(*)
FROM purchase_order_syncs
GROUP BY oc_bejerman, articulo
HAVING COUNT(*) > 1;
