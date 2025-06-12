package sql

const Users_credit_new_instance = `
    INSERT INTO users_credit (
            id,
            user_id,
            invitation_id,
            earned_amount,
            earned_currency)
    VALUES ($1, $2, $3, $4, $5);
`

const Users_credit_fetch_by_id = `
    Select
        id,
        user_id,
        invitation_id,
        earned_amount,
        earned_currency,
        created_at
    from users_credit
    where id = $1
    limit 1
`
const Users_credit_fetch_by_invitationid = `
    Select
        id,
        user_id,
        invitation_id,
        earned_amount,
        earned_currency,
        created_at
    from users_credit
    where invitation_id = $1
    limit 1
`

const User_credit_search_query = ` 
    SELECT
        id,
        user_id,
        invitation_id,
        earned_amount,
        earned_currency,
        created_at
    from users_credit where user_id = $1
`
