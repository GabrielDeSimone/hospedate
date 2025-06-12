package sql

const Properties_new_property = `
    insert into properties (
        id,
        max_guests,
        airbnb_room_id,
        price,
        user_id,
        city
    ) values ($1, $2, $3, $4, $5, $6)
`
const Properties_fetch_by_id = `
    SELECT
        id,
        title,
        description,
        max_guests,
        airbnb_room_id,
        ex_airbnb_room_id,
        price,
        user_id,
        city,
        status,
        is_verified,
        images,
        booking_options,
        created_at,
        accommodation,
        location,
        wifi,
        tv,
        microwave,
        oven,
        kettle,
        toaster,
        coffee_machine,
        air_conditioning,
        heating,
        parking,
        pool,
        gym,
        half_bathrooms,
        bedrooms
    FROM properties
    WHERE id = $1
    LIMIT 1
`

const Properties_delete_by_id = `
    delete from properties where id=$1
`
const Properties_get_blocked_properties = `
   (SELECT dates_overlap.property_id FROM
        (SELECT 
         property_blocked_dates.property_id,greatest(property_blocked_dates.date_init, date( ? )) < least(property_blocked_dates.date_end, date( ? )) as DateRangesOverlap
         FROM
            (
              SELECT property_id,block_init as date_init, block_end as date_end
              FROM blocks
              UNION ALL
              SELECT property_id,checkin_date, checkout_date
              FROM orders
              WHERE status in ('pending', 'confirmed','ephemeral','in_progress')
            ) as property_blocked_dates
        ) dates_overlap
    group by 1 having count(*)!=sum( case when dates_overlap.DateRangesOverlap=false then 1 else 0 end ) )
`

const Block_check_colision = `
    SELECT COALESCE(SUM( case when dates_overlap.DateRangesOverlap=true then 1 else 0 end ),0) as sum_check FROM
        (SELECT 
        greatest(property_blocked_dates.date_init, date($1)) < least(property_blocked_dates.date_end, date($2)) as DateRangesOverlap
        FROM
            (SELECT property_id, block_init as date_init, block_end as date_end
            FROM blocks where property_id= $3
            UNION ALL
            SELECT property_id,checkin_date, checkout_date
            FROM orders where property_id= $3 and status in ('pending', 'confirmed','ephemeral','in_progress') ) as property_blocked_dates
        ) as dates_overlap
`

const Block_new_block = `
    insert into blocks (
        id,
        property_id,
        block_init,
        block_end
    ) values ($1, $2, date($3), date($4))
`

const Properties_fetch_blocks_by_property_id = `
    Select
        id,
        block_init,
        block_end,
        property_id,
        created_at
    from blocks
    where property_id = $1
`
const Properties_delete_block_by_id = `
    delete from blocks where property_id=$1 and id=$2
`
const Properties_search_base_query = `
    SELECT
        id,
        title,
        description,
        max_guests,
        airbnb_room_id,
        ex_airbnb_room_id,
        price,
        user_id,
        city,
        status,
        is_verified,
        images,
        booking_options,
        created_at,
        accommodation,
        location,
        wifi,
        tv,
        microwave,
        oven,
        kettle,
        toaster,
        coffee_machine,
        air_conditioning,
        heating,
        parking,
        pool,
        gym,
        half_bathrooms,
        bedrooms
    FROM properties WHERE 1=1
`

const Properties_fetch_block_by_id = `
    Select
        id,
        block_init,
        block_end,
        property_id,
        created_at
    from blocks
    where property_id = $1 and id = $2
    limit 1
`

const Properties_edit_base_query = ` UPDATE properties SET `

const Properties_set_airbnb_room_id_NULL = ` UPDATE properties SET airbnb_room_id = NULL WHERE id = $1`
