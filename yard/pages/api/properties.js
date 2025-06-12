import utils from '../../utils/utils'
import { withIronSessionApiRoute } from "iron-session/next";

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard

export default withIronSessionApiRoute(newPropertyRoute, utils.ironOptions)

async function newPropertyRoute(req, res) {
    if (req.method !== "POST") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }
    if (!req.session.user || !req.session.user.isHost) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }
    const body = req.body
    if (!body.maxGuests || !body.airbnb_room_id || !body.price || !body.city) {
        res.status(400).json({ error: 'maxGuests, airbnb_room_id, price and city are required' })
        return
    }
    if (! utils.cities.includes(body.city)) {
        res.status(400).json({ error: 'City entered is invalid' })
        return
    }

    const data = {
        "max_guests": body.maxGuests,
        "airbnb_room_id": body.airbnb_room_id,
        "price": body.price,
        "city": body.city,
        "user_id": req.session.user.id
    }

    const options = {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data)
    }

    const [response, status] = await fetchBackyard('/properties', options, res)

    if (status === 201) {
        res.status(201).json({ data: { id: response.id } })
    } else if (status === 400 && response.error === "ErrPropertyAlreadyTaken") {
        res.status(400).json({ error: apiErrors.PropertyAlreadyTaken })
    } else if (status === 400) {
        res.status(400).json({ error: apiErrors.BadRequest })
    } else {
        console.log("Error while creating a property, got this from backend: " + JSON.stringify(response))
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}
