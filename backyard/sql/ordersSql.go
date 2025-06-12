package sql

const Order_new_order = `
    insert into orders (
        id,
        user_id,
        property_id,
        checkin_date,
        checkout_date,
        number_guests,
        price,
        total_billed_cents,
        order_type
    ) values ($1, $2, $3, date($4), date($5), $6, $7, $8, $9)
`

const Orders_fetch_by_id = `
    Select
        id,
        user_id,
        status,
        property_id,
        checkin_date,
        checkout_date,
        number_guests,
        price,
        price_currency,
        total_billed_cents,
        canceled_by,
        order_type,
        created_at,
        wallet_address
    from orders
    where id = $1
    limit 1
`
const Orders_delete_by_id = `
    delete from orders where id=$1
`

const Orders_search_base_query = ` 
    SELECT
        id,
        user_id,
        status,
        property_id,
        checkin_date,
        checkout_date,
        number_guests,
        price,
        price_currency,
        total_billed_cents,
        canceled_by,
        order_type,
        created_at,
        wallet_address
    FROM orders WHERE 1=1
`

const Orders_edit_base_query = ` UPDATE orders SET `

const Orders_fetch_ephemeral = `
    Select
        id,
        user_id,
        status,
        property_id,
        checkin_date,
        checkout_date,
        number_guests,
        price,
        price_currency,
        total_billed_cents,
        canceled_by,
        order_type,
        created_at,
        wallet_address
    from orders
    where status = 'ephemeral'
`
