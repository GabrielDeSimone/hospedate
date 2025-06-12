package sql

const Users_fetch_user = `
    Select
        id,
        name,
        email,
        created_at,
        phone_number,
        is_host
    from users
    where email = $1 and crypt($2, password_enc) = password_enc
    limit 1
`

const Users_new_user = `
    insert into users (
        id,
        name,
        email,
        password_enc,
        phone_number
    ) values ($1, $2, $3, crypt($4, gen_salt('bf', 8)), $5)
`

const Users_fetch_by_id = `
    Select
        id,
        name,
        email,
        created_at,
        phone_number,
        is_host
    from users
    where id = $1
    limit 1
`

const Users_edit_base_query = ` UPDATE users SET `

const Users_new_withdrawal = `
    insert into withdrawals (
        id,
        user_id,
        reclaimed_amount_cents,
        reclaimed_currency,
        wallet_address
    ) values ($1, $2, $3, $4, $5)
`

const Users_get_withdrawal = `
    Select
        id,
        user_id,
        method,
        status,
        reclaimed_amount_cents,
        reclaimed_currency,
        created_at,
        processed_at,
        wallet_address
    from withdrawals
    where id = $1
    limit 1
`

const Users_get_balance_by_id = `
    SELECT user_earned.user_id, (user_earned.earned  - COALESCE(user_withdrawn.reclaimed,0)) as amount_cents
    FROM 
    (Select
        user_id,
        sum(earned_amount_cents) as earned
    from owners_earned
    where user_id = $1 and earned_currency = 'USDT'
    group by user_id) as user_earned
    LEFT JOIN 
    (Select
        user_id,
        sum(reclaimed_amount_cents) as reclaimed
    from withdrawals
    where user_id = $1 and reclaimed_currency = 'USDT' and status != 'aborted' 
    group by user_id) as user_withdrawn on user_earned.user_id = user_withdrawn.user_id
`
const Withdrawals_edit_base_query = ` UPDATE withdrawals SET `

const Withdrawals_search_query = ` 
    SELECT
        id,
        user_id,
        method,
        status,
        reclaimed_amount_cents,
        reclaimed_currency,
        created_at,
        processed_at,
        wallet_address
    FROM withdrawals where user_id = $1
`
