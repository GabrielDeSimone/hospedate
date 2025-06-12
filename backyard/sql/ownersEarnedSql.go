package sql

const Owners_earned_new_order = `
    INSERT INTO owners_earned (
            id,
            order_id,
            user_id,
            earned_amount_cents,
            earned_currency)
    VALUES ($1, $2, $3, $4, $5);
`

const Owners_earned_fetch_by_id = `
    Select
        id,
        order_id,
        user_id,
        earned_amount_cents,
        earned_currency,
        created_at
    from owners_earned
    where id = $1
    limit 1
`
const Owners_earned_search_query = ` 
    SELECT
        id,
        order_id,
        user_id,
        earned_amount_cents,
        earned_currency,
        created_at
    from owners_earned where user_id = $1
`
