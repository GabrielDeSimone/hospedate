import utils from '../../../utils/utils'
import { withIronSessionApiRoute } from 'iron-session/next'

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard

export default withIronSessionApiRoute(PropertiesIdRoute, utils.ironOptions)

const editableFieldsToBackyardMap = {
    "title": "title",
    "description": "description",
    "maxGuests": "max_guests",
    "city": "city",
    "price": "price",
    "accommodation": "accommodation",
    "location": "location",
    "wifi": "wifi",
    "tv": "tv",
    "parking": "parking",
    "microwave": "microwave",
    "oven": "oven",
    "kettle": "kettle",
    "toaster": "toaster",
    "coffeeMachine": "coffee_machine",
    "airConditioning": "air_conditioning",
    "heating": "heating",
    "pool": "pool",
    "gym": "gym",
    "halfBathrooms": "half_bathrooms",
    "bedrooms": "bedrooms",
}

function collectEditFields(req) {
    const fieldsReceived = {}

    Object.keys(editableFieldsToBackyardMap).map((field) => {
        if (Object.keys(req.body).includes(field)) {
            fieldsReceived[editableFieldsToBackyardMap[field]] = req.body[field]
        }
    })
    return fieldsReceived
}

async function PropertiesIdRoute(req, res) {
    if (req.method === "DELETE") {
        // call delete handler
        return RemovePropertyRoute(req, res)
    } else if (req.method === "PUT") {
        return EditPropertyRoute(req, res)
    } else if (req.method !== "GET") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }
    const user = req.session.user || null;
    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }

    const [response, status] = await fetchBackyard(`/properties/${req.query.id}`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
        }
    }, res)

    if (status === 200 && response.data.status !== 'archived') {
        res.status(200).json({
            data: {
                id: response.data.id,
                title: response.data.title,
                maxGuests: response.data.max_guests,
                description: response.data.description,
                airbnb_room_id: response.data.airbnb_room_id,
                price: response.data.price,
                user_id: response.data.user_id,
                status: response.data.status,
                city: response.data.city,
                created_at: response.data.created_at,
                images: response.data.images,
                bookingOptions: response.data.booking_options,
                ...Object.fromEntries(
                    Object.keys(utils.propAmenitiesFieldsBackyardMap).map(
                        (field) => [field, response.data[utils.propAmenitiesFieldsBackyardMap[field]]]
                    )
                )
            }
        })
    } else if (status === 200) { // archived
        res.status(404).json({ error: apiErrors.NotFound })
    } else if (status === 400) {
        res.status(400).json({ error: apiErrors.BadRequest })
    } else if (status === 404) {
        res.status(404).json({ error: apiErrors.NotFound })
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}

async function RemovePropertyRoute(req, res) {
    if (req.method !== "DELETE") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }
    const user = req.session.user || null;
    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }
    if (!req.query.id) {
        res.status(400).json({ error: apiErrors.BadRequest })
        return
    }

    let [response, status] = await fetchBackyard(`/properties/${req.query.id}`, null, res)

    if (status === 404) {
        res.status(404).json({ error: apiErrors.NotFound })
        return
    } else if (status !== 200) {
        console.log(`Unexpected status ${status} when fetching a property from backyard`)
        res.status(500).json({ error: apiErrors.UnknownError })
        return
    }

    // if the caller is not the owner of the property, respond with NotFound
    if (response.data.user_id !== user.id) {
        res.status(404).json({ error: apiErrors.NotFound })
        return
    }

    // caller is the owner, now proceed with deletion
    [response, status] = await fetchBackyard(`/properties/${req.query.id}`, {
        method: "PUT",
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            status: "archived"
        })
    }, res)

    if (status === 200) {
        res.status(200).json({ data: 'ok' })
    } else if (status === 400 && response.error === "ErrPropertyHasActiveOrders") {
        res.status(400).json({ error: apiErrors.PropertyActiveOrders })
    } else {
        console.log(`Error when archiving a property from backyard. Response was: ${JSON.stringify(response)}`)
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}

async function EditPropertyRoute(req, res) {
    if (req.method !== "PUT") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }
    const user = req.session.user || null;
    if (!user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }
    if (!req.query.id) {
        res.status(400).json({ error: apiErrors.BadRequest })
        return
    }
    const receivedFields = collectEditFields(req)
    if (Object.keys(receivedFields).length === 0) {
        console.log('maybe2')
        console.log('changed fields are')
        console.log(receivedFields)
        res.status(400).json({ error: apiErrors.BadRequest })
        return
    }

    const [response, status] = await fetchBackyard(`/properties/${req.query.id}`, null, res)

    if (status === 200) {
        const ownerId = response.data.user_id
        if (ownerId === user.id) {
            const [editResponse, editStatus] = await fetchBackyard(`/properties/${req.query.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(receivedFields)
            }, res)
            if (editStatus === 200) {
                res.status(200).json({
                    id: editResponse.data.id
                })
            } else {
                console.log("Error updating property, status from backyard", editStatus, ", error code: ", editResponse.error)
                res.status(500).json({ error: apiErrors.UnknownError })
            }
        } else {
            res.status(404).json({ error: apiErrors.NotFound })
        }
    } else if (status === 404) {
        res.status(404).json({ error: apiErrors.NotFound })
    } else {
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
