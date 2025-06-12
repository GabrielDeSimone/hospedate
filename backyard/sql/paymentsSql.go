package sql

const Payments_new_payment = `
    insert into payments (
        id,
        order_id,
        traveler_amount_cents,
        traveler_currency
    ) values ($1, $2, $3, $4)
`
const Payments_fetch_by_id = `
    Select
        id,
        order_id,
        method,
        status,
        traveler_amount_cents,
        traveler_currency,
        received_amount_cents,
        received_currency,
        reverted_amount_cents,
        reverted_currency,
        created_at
    from payments
    where id = $1
    limit 1
`

const Payments_fetch_by_order_id = `
    Select
        id,
        order_id,
        method,
        status,
        traveler_amount_cents,
        traveler_currency,
        received_amount_cents,
        received_currency,
        reverted_amount_cents,
        reverted_currency,
        created_at
    from payments
    where order_id = $1
    limit 1
`

const Payments_delete_by_id = `
    delete from payments where id=$1
`

const Payments_delete_by_order_id = `
    delete from payments where order_id=$1
`

const Payments_search_base_query = ` 
    SELECT
        id,
        order_id,
        method,
        status,
        traveler_amount_cents,
        traveler_currency,
        received_amount_cents,
        received_currency,
        reverted_amount_cents,
        reverted_currency,
        created_at
    FROM payments WHERE 1=1
`

const Payments_edit_base_query = ` UPDATE payments SET `

const Payments_fetch_PK_by_order_id = `
    Select
        private_key,
    from payments
    where order_id = $1
    limit 1
`
