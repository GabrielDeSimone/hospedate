import utils from '../../utils/utils'
import { withIronSessionApiRoute } from "iron-session/next";

const apiErrors = utils.apiErrors
const fetchBackyard = utils.fetchBackyard

export default withIronSessionApiRoute(newReservationRoute, utils.ironOptions)

async function newReservationRoute(req, res) {
    if (req.method !== "POST") {
        res.status(405).json({ error: apiErrors.UnsupportedMethod })
        return
    }
    if (!req.session.user) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    }
    // params check
    const body = req.body
    if (!body.checkinDate || !body.checkoutDate || !body.guests || !body.propertyId || !body.reservationType) {
        res.status(400).json({ error: 'checkinDate, checkoutDate, guests, propertyId and reservationType are required' })
        return
    }

    // let [response, status] = await fetchBackyard(`/properties/${body.propertyId}`, null, res)
    let [response, status] = await fetchBackyard(`/properties/${body.propertyId}`, null, res)

    if (status !== 200) {
        res.status(400).json({ error: apiErrors.BadRequest })
        return
    }
    // owner cannot make a reservation of a property that is their own!
    if (response.data && response.data.user_id && response.data.user_id === req.session.user.id) {
        res.status(403).json({ error: apiErrors.Forbidden })
        return
    } else if (response.data && !response.data.user_id) {
        console.log("Could not find a user_id when fetching a property")
        res.status(500).json({ error: apiErrors.UnknownError })
        return
    } else if (response.error) {
        console.log(`Encountered an error when fetching a property from backyard: ${JSON.stringify(response.error)}`)
        res.status(500).json({ error: apiErrors.UnknownError })
        return
    }

    [response, status] = await fetchBackyard('/orders', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            date_start: body.checkinDate,
            date_end: body.checkoutDate,
            number_guests: body.guests,
            user_id: req.session.user.id,
            property_id: body.propertyId,
            order_type: body.reservationType,
        })
    }, res)

    if (status === 201) {
        res.status(201).json({ data: {
            id: response.data.id,
            checkinDate: response.data.date_start,
            checkoutDate: response.data.date_end,
            guests: response.data.number_guests,
            propertyId: response.data.property_id,
            reservationType: response.data.order_type,
        }})
    } else if (status === 400) {
        res.status(400).json({ error: apiErrors.BadRequest })
    } else {
        console.log("Error while creating a reservation, got this from backend: " + JSON.stringify(response))
        res.status(500).json({ error: apiErrors.UnknownError })
    }
}