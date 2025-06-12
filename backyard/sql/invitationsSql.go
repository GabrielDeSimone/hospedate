package sql

const Invitations_fetch_by_id = `
    select
        id,
        used_by,
        kind,
        generated_by,
        created_at
    from users_invitations
    where id = $1
    limit 1;
`

const Invitations_new_invitation = `
    insert into users_invitations (
        id,
        kind,
        generated_by
    ) values ($1, $2, $3);
`

const Invitations_edit_base_query = ` UPDATE users_invitations SET `

const Invitations_search_base_query = ` 
    SELECT
        id,
        used_by,
        kind,
        generated_by,
        created_at
    from users_invitations WHERE 1=1
`
